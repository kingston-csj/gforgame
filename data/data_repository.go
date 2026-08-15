package data

import "reflect"

// 暂未适配
type DataRepository interface {
	// 根据模型类型+id查询单行配置
	QueryById(modelType reflect.Type, id any) any
	// 查询模型全量数据
	QueryAll(modelType reflect.Type) []any
	// 根据数据模型类型获取对应容器
	QueryContainer(modelType reflect.Type, containerType reflect.Type) any
	// 直接根据容器类型获取容器实例
	GetSpecificContainer(containerType reflect.Type) any
	// 热更重载表
	Reload(table string)
}