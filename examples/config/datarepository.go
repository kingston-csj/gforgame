package config

import "reflect"

// 暂未适配
type DataRepository interface {
	/**
	 * 查询配置容器（对应原 QueryContainer<T>）
	 * @return any 实际返回 Container 子类，调用方通过类型断言转换
	 */
	QueryContainer(tableClass reflect.Type, containerClass reflect.Type) any

	/**
	 * 根据主键查询（对应原 QueryById<E>）
	 * @return any 实际返回 E 类型，调用方通过类型断言转换
	 */
	QueryById(clazz reflect.Type, id any) any

	/**
	 * 查询所有数据（对应原 QueryAll<E>）
	 * @return []any 实际返回 []E，调用方通过类型断言转换
	 */
	QueryAll(clazz reflect.Type) []any

	/**
	 * 根据索引查询列表（对应原 QueryByIndex<E>）
	 */
	QueryByIndex(clazz reflect.Type, name string, index any) []any

	/**
	 * 根据唯一索引查询（对应原 QueryByUniqueIndex<E>）
	 */
	QueryByUniqueIndex(clazz reflect.Type, name string, index any) any

	/**
	 * 重载数据（保持不变）
	 */
	Reload(table string)
}