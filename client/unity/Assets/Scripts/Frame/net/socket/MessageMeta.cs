using System;

namespace Nova.Net.Socket
{
    /// <summary>
    /// 消息元数据
    /// 用于标记Message 类
    /// </summary>
    [AttributeUsage(AttributeTargets.Class)]
    public class MessageMeta : Attribute
    {
        /// <summary>
        /// 消息cmd 每个消息都有一个唯一的cmd
        /// </summary>
        public int Cmd;
    }
}