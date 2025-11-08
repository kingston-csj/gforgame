namespace Nova.Data
{
    using System;
    using System.Reflection;

    /// <summary>
    /// 方法索引元数据：基于方法返回值的索引实现
    /// </summary>
    public class MethodIndexMeta<I> : IIndexMeta<I> where I : AbsConfigItem
    {
        private readonly MethodInfo _method;
        public string Name { get; }
        public bool IsUnique { get; }

        public MethodIndexMeta(MethodInfo method)
        {
            _method = method;
            // 校验方法（无参数）
            if (method.GetParameters().Length > 0)
                throw new ArgumentException($"索引方法 {method.Name} 不能有参数");

            // 获取特性配置
            var indexAttr = method.GetCustomAttribute<IndexAttribute>();
            if (indexAttr == null)
                throw new ArgumentException($"方法 {method.Name} 未标记 [Index] 特性");

            // 确定索引名称（特性指定优先，否则用方法名）
            Name = string.IsNullOrEmpty(indexAttr.Name) ? method.Name : indexAttr.Name;
            IsUnique = indexAttr.IsUnique;

            // 允许访问私有方法
            if (!_method.IsPublic)
                throw new ArgumentException($"索引方法 {method.Name} 必须为公共方法");
        }

        /// <summary>调用方法获取返回值作为索引值</summary>
        public object GetValue(I item)
        {
            try
            {
                return _method.Invoke(item, parameters: null);
            }
            catch (Exception e)
            {
                throw new InvalidOperationException($"调用方法 {_method.Name} 获取索引值失败", e);
            }
        }
    }
}