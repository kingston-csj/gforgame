namespace Nova.Data
{
    using System;
    using System.Reflection;

    /// <summary>
    /// 字段索引元数据：基于字段的索引实现
    /// </summary>
    public class FieldIndexMeta<I> : IIndexMeta<I> where I : AbsConfigData
    {
        private readonly FieldInfo _field;
        public string Name { get; }
        public bool IsUnique { get; }

        public FieldIndexMeta(FieldInfo field)
        {
            _field = field;
            // 获取特性配置
            var indexAttr = field.GetCustomAttribute<IndexAttribute>();
            if (indexAttr == null)
                throw new ArgumentException($"字段 {field.Name} 未标记 [Index] 特性");

            // 确定索引名称（特性指定优先，否则用字段名）
            Name = string.IsNullOrEmpty(indexAttr.Name) ? field.Name : indexAttr.Name;
            IsUnique = indexAttr.IsUnique;

            // 允许访问私有字段
            if (!_field.IsPublic)
                throw new ArgumentException($"字段 {field.Name} 不是公共字段");
        }

        /// <summary>从配置项中获取字段值作为索引值</summary>
        public object GetValue(I item)
        {
            try
            {
                return _field.GetValue(item);
            }
            catch (Exception e)
            {
                throw new InvalidOperationException($"获取字段 {_field.Name} 的值失败", e);
            }
        }
    }
}