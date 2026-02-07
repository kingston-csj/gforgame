using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using Nova.Commons.Util;
using Nova.Logger;
using UnityEngine;

namespace Nova.Data
{
    /// <summary>
    ///     基础配置类容器，用于管理对应表配置数据
    ///     如果配置类无需二级缓存数据，直接使用该类作为配置容器即可
    /// </summary>
    public class ConfigContainer<E> : IDisposable where E : AbsConfigData
    {
        // 组件集合
        protected Dictionary<int, E> _primaryMap = new();

        // 索引映射
        protected Dictionary<string, List<E>> _indexMapper = new();
        protected List<IIndexMeta<E>> _indexMetas = new();

        /// <summary>
        ///     配置记录（相当于excel的一行记录）
        ///  key为id主键，value为记录数据
        /// </summary>
        protected Dictionary<int, E> _items;

        /// <summary>
        ///     构造函数，从资源绑定中加载配置文件
        /// </summary>
        /// <param name="textAsset">配置文件</param>
        public ConfigContainer(TextAsset textAsset)
        {
            // 自动扫描索引元数据（字段和方法）
            ScanIndexMetas();
            _initFromJsonItems(textAsset);
        }

        /// <summary>自动扫描记录类型中标记 [Index] 的字段和方法，构建索引元数据</summary>
        private void ScanIndexMetas()
        {
            Type itemType = typeof(E);

            // 1. 扫描字段索引
            foreach (var field in itemType.GetFields(BindingFlags.Public | BindingFlags.NonPublic |
                                                     BindingFlags.Instance))
            {
                if (field.GetCustomAttribute<IndexAttribute>() != null)
                {
                    _indexMetas.Add(new FieldIndexMeta<E>(field));
                }
            }

            // 2. 扫描方法索引（无参数的方法）
            foreach (var method in itemType.GetMethods(BindingFlags.Public | BindingFlags.NonPublic |
                                                       BindingFlags.Instance))
            {
                if (method.GetCustomAttribute<IndexAttribute>() != null && method.GetParameters().Length == 0)
                {
                    _indexMetas.Add(new MethodIndexMeta<E>(method));
                }
            }

            // 校验索引名称唯一性
            var duplicateNames = _indexMetas.GroupBy(m => m.Name)
                .Where(g => g.Count() > 1)
                .Select(g => g.Key);
            if (duplicateNames.Any())
            {
                throw new InvalidOperationException($"存在重复的索引名称：{string.Join(",", duplicateNames)}");
            }
        }

        /// <summary>
        ///     从TextAsset初始化配置数据(json文本)
        /// </summary>
        /// <param name="textAsset">包含json配置的文本资源</param>
        private void _initFromJsonItems(TextAsset textAsset)
        {
            var items = JsonUtil.FromJson<E[]>(textAsset.text);
            // list转字典
            _items = new Dictionary<int, E>();
            foreach (var item in items)
            {
                if (item.id == 0)
                {
                    // id为0，代表空白行，直接跳过
                    continue;
                }

                _items.Add(item.id, item);
            }

            HashSet<string> uniqueIndexKeys = new();
            foreach (var item in items)
            {
                if (item.id == 0) continue;

                // 主键校验
                if (_primaryMap.ContainsKey(item.id))
                    throw new InvalidOperationException($"主键重复：{item.id}");
                _primaryMap[item.id] = item;

                // 构建所有索引
                foreach (var indexMeta in _indexMetas)
                {
                    object indexValue = indexMeta.GetValue(item); // 调用字段/方法获取索引值
                    string indexKey = GenerateIndexKey(indexMeta.Name, indexValue);

                    // 唯一索引校验
                    if (indexMeta.IsUnique)
                    {
                        if (uniqueIndexKeys.Contains(indexKey))
                            throw new InvalidOperationException($"唯一索引 {indexMeta.Name} 重复：{indexKey}");
                        uniqueIndexKeys.Add(indexKey);
                    }

                    // 添加到索引映射
                    if (!_indexMapper.ContainsKey(indexKey))
                        _indexMapper[indexKey] = new List<E>();
                    _indexMapper[indexKey].Add(item);
                }
            }

            AfterLoad();
        }

        /// <summary>生成索引键（格式："索引名称:值"）</summary>
        private string GenerateIndexKey(string indexName, object value)
        {
            return $"{indexName}:{value}";
        }

        /// <summary>
        ///     初始化后的钩子，子类可重写此方法进行额外处理
        ///     例如：新建二级缓存，进行数据完整性验证等等
        /// </summary>
        protected virtual void AfterLoad()
        {
            // 默认为空
        }

        /// <summary>
        ///     根据ID获取记录
        /// </summary>
        /// <param name="id">记录ID</param>
        /// <returns>匹配的记录，未找到时返回null</returns>
        public E GetItem(int id)
        {
            if (_items.TryGetValue(id, out E item))
            {
                return item;
            }

            return null;
        }

        /// <summary>
        ///     获取配置列表
        /// </summary>
        /// <returns></returns>
        public E[] GetItems()
        {
            return _items.Values.ToArray();
        }

        /// <summary>
        /// 按索引名称查询列表
        /// </summary>
        public List<E> GetItemsByIndex(string indexName, object indexValue)
        {
            var indexMeta = _indexMetas.FirstOrDefault(m => m.Name == indexName);
            if (indexMeta == null)
            {
                LoggerUtil.Error($"索引 {indexName} 未定义");
                return new List<E>();
            }

            string indexKey = GenerateIndexKey(indexName, indexValue);
            return _indexMapper.TryGetValue(indexKey, out var items) ? new List<E>(items) : new List<E>();
        }

        /// <summary>
        /// 按唯一索引查询单个结果
        /// </summary>
        public E GetUniqueItemByIndex(string indexName, object indexValue)
        {
            var items = GetItemsByIndex(indexName, indexValue);
            return items.Count == 0 ? null : items[0];
        }

        /// <summary>
        /// 清理资源，防止内存泄漏
        /// </summary>
        public void Dispose()
        {
            _items = null;
            _indexMapper = null;
            _indexMetas = null;
            _primaryMap = null;
        }
    }
}