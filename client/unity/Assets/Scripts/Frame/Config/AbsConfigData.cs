using System;
using System.Collections.Generic;

namespace Nova.Data
{
    /// <summary>
    ///     配置项基类，提供通用属性和访问方法
    /// </summary>
    [Serializable]
    public abstract class AbsConfigData
    {
        /// <summary>
        ///     唯一标识
        /// </summary>
        public int id;

        /// <summary>
        ///    名称
        /// </summary>
        public string name;

        /// <summary>
        ///     描述信息
        /// </summary>
        public string desc;
    }
}