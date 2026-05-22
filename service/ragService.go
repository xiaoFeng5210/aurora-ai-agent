package service

import (
	"aurora-agent/service/embedding"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"go.uber.org/zap"

	rabbitmq_module "aurora-agent/database/rabbitmq"
	"aurora-agent/database/redis"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	exg        = "rag"
	routingKey = "vectorize_key"
)

// FileParser 文件解析器接口
// 不同文件类型实现此接口，将文件二进制数据解析为纯文本
// 扩展新格式时只需：1. 实现此接口  2. 注册到 fileParsers 即可
type FileParserTrait interface {
	Parse2String(data []byte) (string, error)
}

type RagVectorizeMsg struct {
	Uid      int    `json:"uid"`
	Filename string `json:"filename"`
	// FileContent []byte `json:"file_content"`
	FileCloudPath string `json:"file_cloud_path"`
}

// fileParsers 文件扩展名 → 解析器的注册表
// 新增格式在此注册，handler 层的白名单也需同步更新
var fileParsers = map[string]FileParserTrait{
	".txt": &TextParser{},
	".md":  &TextParser{},
	".csv": &TextParser{},
	// TODO: 未来扩展
	// ".pdf":  &PdfParser{},
	// ".docx": &WordParser{},
}

// ---------------------- 解析器实现 ----------------------

// TextParser 纯文本解析器，适用于 .txt / .md / .csv 等文本格式
type TextParser struct{}

func (p *TextParser) Parse2String(data []byte) (string, error) {
	text := strings.TrimSpace(string(data))
	if text == "" {
		return "", fmt.Errorf("文件内容为空")
	}
	return text, nil
}

// TODO: PdfParser PDF 解析器（待实现）
// type PdfParser struct{}
// func (p *PdfParser) Parse(data []byte) (string, error) { ... }

// TODO: WordParser Word(.docx) 解析器（待实现）
// type WordParser struct{}
// func (p *WordParser) Parse(data []byte) (string, error) { ... }

// ---------------------- Service ----------------------

// CreateRag 将上传的文件写入当前用户的个人 RAG 知识库
// 流程：根据扩展名选择解析器 → 解析为文本 → 分块 → 向量化 → 按 user_id 写入 Qdrant
func CreateRag(uid int, file io.Reader, filename string) (*qdrant.UpdateResult, error) {
	// 1. 根据文件扩展名匹配解析器
	ext := strings.ToLower(filepath.Ext(filename))
	parser, ok := fileParsers[ext]
	if !ok {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 2. 读取文件二进制数据
	var data []byte
  
	buffer := make([]byte, 1024*4)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("读取文件失败: %w", err)
		}

		data = append(data, buffer[:n]...)
	}

	// 3. 通过解析器将文件数据转为纯文本
	text, err := parser.Parse2String(data)
	if err != nil {
		return nil, fmt.Errorf("文件解析失败: %w", err)
	}

	// 4. 文本分块 → 向量化
	doc := &embedding.MdDocument{
		Content: text,
	}
	doc.Chunk()

	_, err = doc.Embedding()
	if err != nil {
		return nil, fmt.Errorf("向量化失败: %w", err)
	}

	// 5. 写入 Qdrant，user_id 用于隔离不同用户的知识库
	logger.Info("RAG 创建: 正在写入知识库",
		zap.Int("uid", uid),
		zap.String("filename", filename),
	)
	result, err := doc.UpsertQdrantVector(uid, filename)
	if err != nil {
		logger.Error("RAG 创建: 写入 Qdrant 失败", zap.Error(err))
		return nil, fmt.Errorf("写入知识库失败: %w", err)
	}

	logger.Info("RAG 创建成功",
		zap.Int("uid", uid),
		zap.String("filename", filename),
		zap.Any("result", result),
	)
	return result, nil
}

func SendRagVectorizeTask(uid int, file io.Reader, filename string, fileCloudPath string) error {
	msg := &RagVectorizeMsg{
		Uid:           uid,
		Filename:      filename,
		FileCloudPath: fileCloudPath,
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}
	conn, err := rabbitmq_module.Connect()
	if err != nil {
		return fmt.Errorf("connect to rabbitmq failed: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel failed: %w", err)
	}
	defer ch.Close()

	err = rabbitmq_module.DeclareExchange(ch, exg)
	if err != nil {
		return err
	}

	// 发送消息
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		exg,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonMsg,
		},
	)
	if err != nil {
		return fmt.Errorf("publish message failed: %w", err)
	}

	// 需要存进redis里
	client := redis.Client()
	redisCtx := context.Background()
	pipe := client.TxPipeline()
	key := fmt.Sprintf("rag_vectorize.%d", uid)
	redisData := map[string]interface{}{
		"filename": filename,
		"status":   "pending",
	}
	pipe.HSet(redisCtx, key, redisData)
	pipe.Expire(redisCtx, key, 3*time.Hour)
	if _, err = pipe.Exec(redisCtx); err != nil {
		return err
	}
	return nil
}

func ConsumeRagVectorizeTask() {
	conn, err := rabbitmq_module.Connect()
	if err != nil {
		panic("connect to rabbitmq failed: " + err.Error())
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic("open channel failed: " + err.Error())
	}
	defer ch.Close()
	// 声明exchange,这个很关键，最好都声明下保证exchange存在
	err = rabbitmq_module.DeclareExchange(ch, exg)
	if err != nil {
		panic("declare exchange failed: " + err.Error())
	}

	// 声明队列
	q, err := ch.QueueDeclare(
		"upload_and_vectorize_queue",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		panic("declare queue failed: " + err.Error())
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exg,
		false,
		nil,
	)
	if err != nil {
		panic("bind queue failed: " + err.Error())
	}

	deliverCh, err := ch.Consume(
		q.Name,
		"",
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic("consume failed: " + err.Error())
	}

	for delivery := range deliverCh {
		logger.Info("VectorizeTaskConsumer Received message: ", zap.String("message", string(delivery.Body)))
		var message RagVectorizeMsg
		err := json.Unmarshal(delivery.Body, &message)
		if err != nil {
			logger.Error("unmarshal message failed: " + err.Error())
			delivery.Nack(false, true)
			continue
		}
		fileData, err := DonwloadFileFromBaiduNetworkdisk(message.FileCloudPath)
		if err != nil {
			logger.Error("download file from baidu networkdisk failed: " + err.Error())
			delivery.Nack(false, true)
			continue
		}
		_, err = CreateRag(message.Uid, bytes.NewReader(fileData), message.Filename)
		if err != nil {
			logger.Error("create rag failed: " + err.Error())
			delivery.Nack(false, true)
			continue
		}

		// 更新redis状态
		client := redis.Client()
		key := fmt.Sprintf("rag_vectorize.%d", message.Uid)
		err = client.Del(context.Background(), key).Err()
		if err != nil {
			logger.Error("delete redis key failed: " + err.Error())
		}
		delivery.Ack(false)
	}
}
