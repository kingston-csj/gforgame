package persist

// 持久化容器， 提供以下类型
// 基于队列的持久化容器
// 基于延迟的持久化容器
type PersistContainer interface {

	// 接收实体
	Receive(entity Entity)

	// 关闭容器, 确保将所有缓存数据写入到数据库
	ShutdownGraceful()

	// 当前等待入库的队列大小
	Size() int
}