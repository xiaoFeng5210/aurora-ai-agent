# 数据结构和算法
作为想做全栈的开发者（而不是搞竞赛），算法不用贪多。按优先级，我建议你先掌握这 4 个，足够覆盖 90% 的实际工作和面试：

  优先学（按顺序）

  1. 哈希表 + 数组/字符串
  最高频、最实用。去重、计数、查找、缓存几乎天天用。Go 的 map、Python 的 dict、前端的 Object/Map 都是它。性价比第一。

  2. 双指针 / 滑动窗口
  处理数组、字符串、区间问题的利器。思路简单、应用广（限流、滑动统计、字符串匹配）。

  3. 二分查找 + 排序
  不是让你手写快排，而是理解「有序 → log n」的思想。数据库索引、版本回滚定位、Rust/Go 标准库的 sort 都基于此。

  4. 树与图的遍历（BFS / DFS + 递归）
  全栈强相关：前端的 DOM 树、组件树、路由树；后端的依赖解析、文件系统、权限继承。递归思维也是这里练出来的。

  暂时缓一缓

  - 动态规划：思维门槛高，实际工作用得少，等前面扎实了再碰。
  - 复杂图算法（最短路、最小生成树、网络流）：用到时再查，不必现在啃。
  - 高级数据结构（线段树、并查集、Trie）：并查集和 Trie 可以了解概念，其余按需。

  一句话路线

  ▎ 哈希表 → 双指针 → 二分/排序 → 树图遍历，每个配 10~20 道题练手，就够支撑你做全栈了。

  语言学习上也有个建议的先后顺序，如果你想我帮你排一下 Go / Python / 前端 / Flutter / Rust
  的学习路径，告诉我你目前最熟哪个、想先做哪类项目（Web 后端 / App / 全栈产品），我可以给你一条更具体的路线。


  ## 难点
  1. 流式响应里拼装 Function Call（最硬的一块）

  ai/llm/glm.go:208 的 mergeToolCalls 是整个项目最微妙的地方。GLM 用 SSE
  流式返回，工具调用的参数是被切成碎片一段段吐出来的，你得自己按 Index 把它们重新拼回完整 JSON：

  case strings.HasPrefix(toolCall.Function.Arguments, current.Function.Arguments):
      current.Function.Arguments = toolCall.Function.Arguments
  default:
      current.Function.Arguments += toolCall.Function.Arguments

  这里既要处理"增量追加"又要处理"全量覆盖"两种流式协议，边界判断很容易出 bug，拼错一个字符整个 json.Unmarshal 就崩。这是 Agent
  能不能正常工作的命门。

  2. Agent 多轮循环的状态管理

  ai/agent/core.go:68 的 RunAgent 里同时维护 conversation（本轮上下文）和 a.History（持久历史）两份消息列表，但工具结果只追加进
  conversation，最终回答才写回 History。两份状态的同步规则不一致，多轮对话 +
  工具调用混在一起时很容易出现上下文丢失或重复。MaxLoop=6 的死循环保护是必须的，但目前到达上限直接报错，没有"降级回答"。

  3. 异步向量化的可靠性（有隐患）

  service/ragService.go:249 的消费者链路是：RabbitMQ → 下载百度网盘文件 → 解析 → 向量化 → 写 Qdrant → 更新
  Redis。问题在于失败时一律：

  delivery.Nack(false, true)  // requeue=true

  如果某个文件本身有问题（比如永久解析失败），它会被无限重新入队、无限重试——典型的"毒消息"问题。这里缺少死信队列和重试次数上限，
  生产环境会出事。

  4. RAG 检索质量

  service/embedding/mdDocument.go:28 的分块是按 500 个字符硬切，没有重叠（overlap）、不考虑句子/段落边界。一句话被从中间切断后向
  量语义就乱了，直接影响召回质量。而且检索后没有重排序（rerank），让大模型自己判断"要不要查知识库"（rag.go 里那段超长
  prompt）也很脆弱。

  5. 其他工程难点

  - SSE 端到端链路：GLM 的 SSE → 后端解析 → 后端再 SSE 给前端，中间任何一环断连都要兜底（CHANGELOG 里专门记了"SSE
  连接异常终止日志"）。
  - 百度网盘当文件存储：OAuth token 存 Redis、要处理刷新和限流，是个不常规但成本低的选择，复杂度都在 token 生命周期管理上。
  - 多服务编排：docker-compose 跑 6
  个服务（PG/Redis/RabbitMQ/Qdrant/前端×2/后端），依赖顺序、健康检查、以及为绕开服务器构建慢而改成"本地预编译再塞进镜像"（见近期
   git log），都是部署上踩过坑的体现。
