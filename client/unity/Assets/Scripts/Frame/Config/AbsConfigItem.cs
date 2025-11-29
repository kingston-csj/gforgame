using System;
using System.Collections.Generic;

namespace Nova.Data
{
    /// <summary>
    ///     配置项基类，用于存储基础配置数据
    ///     作为所有配置项的父类，提供通用属性和访问方法
    /// </summary>
    [Serializable]
    public class AbsConfigItem
    {
        /// <summary>
        ///     唯一标识ID
        ///     用于在配置系统中唯一标识该配置项
        /// </summary>
        public int id;

        /// <summary>
        ///     配置项名称
        /// </summary>
        public string name;

        /// <summary>
        ///     描述信息
        /// </summary>
        public string desc;

        protected Dictionary<string, object> _properties;

        public Dictionary<string, object> properties
        {
            get
            {
                if (_properties == null)
                {
                    InitProperties();
                }

                return _properties;
            }
        }

        protected virtual void InitProperties()
        {
            _properties = new Dictionary<string, object>();
            _properties.Add("id", id);
            _properties.Add("name", name);
            _properties.Add("desc", desc);
        }
    }
}