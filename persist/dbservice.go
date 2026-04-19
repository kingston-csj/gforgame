package persist

type DBService interface {

	// 数据实体持久化到数据库
	// 该接口会保证无论表记录是否已存在于数据库
	// 若存在则执行更新操作，否则执行插入动作
	SaveToDb(entity Entity)

	// 删除数据（游戏业务一般只作更新，不作删除，这个接口使用场景很少）
	DeleteEntityFromDb(entity Entity)

	// 关闭服务, 确保将所有缓存数据写入到数据库
	Shutdown()
}