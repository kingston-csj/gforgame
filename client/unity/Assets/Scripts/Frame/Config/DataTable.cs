using System;

namespace Nova.Net.Socket
{
    /// <summary>
    /// 标记一个类代表配置表
    /// </summary>
    [AttributeUsage(AttributeTargets.Class)]
    public class DataTable : Attribute
    {
        /// <summary>
        /// 对应的文件名称
        /// </summary>
        public string name;
    }
}