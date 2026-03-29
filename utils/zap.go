package utils

import (
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	Logger = InitZap("log/zap")
}

func InitZap(logFile string) *zap.Logger {
	rotateOut, err := rotatelogs.New(
		logFile+".%Y%m%d%H.log",                   //指定日志文件的路径和名称，路径不存在时会创建
		rotatelogs.WithLinkName(logFile+".log"),   //为最新的一份日志创建软链接
		rotatelogs.WithRotationTime(24*time.Hour), //每隔24小时生成一份新的日志文件
		rotatelogs.WithMaxAge(2*24*time.Hour),     //只留最近7天的日志，或使用WithRotationCount只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") //指定时间格式
	encoderConfig.TimeKey = "time"                                                    //默认是ts
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder                           //指定level的显示样式
	core := zapcore.NewCore(
		// zapcore.NewJSONEncoder(encoderConfig),  //日志为json格式
		zapcore.NewConsoleEncoder(encoderConfig), //日志为console格式（Field还是json格式）
		// zapcore.AddSync(file),                    //指定输出到文件
		// zapcore.AddSync(lumberJackLogger), //指定输出到文件
		zapcore.AddSync(rotateOut), //指定输出到文件
		zapcore.InfoLevel,          //设置最低级别
	)

	logger := zap.New(
		core,
		zap.AddCaller(),                       //上报文件名和行号
		zap.AddStacktrace(zapcore.ErrorLevel), //error级别及其以上的日志打印调用堆栈
		zap.Hooks(func(e zapcore.Entry) error { //可以添加多个钩子，在输出日志【之前】执行钩子
			if e.Level >= zapcore.ErrorLevel {
				fmt.Println(e.Message)
			}
			return nil
		}),
	)

	logger = logger.With(
		zap.Namespace("aurora-agent"), //后续的Field都记录在此命名空间中
		//通过zap.String、zap.Int等显式指定类型；fmt.Printf之类的方法大量使用interface{}和反射，性能损失不少
	)
	return logger
}
