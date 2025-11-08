namespace Nova.Data
{
    using System;

    /// <summary>
    /// 索引特性：标记在字段或方法上，声明该成员作为索引
    /// </summary>
    [AttributeUsage(AttributeTargets.Field | AttributeTargets.Method, AllowMultiple = false)]
    public class IndexAttribute : Attribute
    {
        /// <summary>索引名称（未指定则使用字段/方法名）</summary>
        public string Name { get; set; } = "";

        /// <summary>是否唯一索引</summary>
        public bool IsUnique { get; set; } = false;

        public IndexAttribute() { }

        public IndexAttribute(string name, bool isUnique = false)
        {
            Name = name;
            IsUnique = isUnique;
        }
    }
}