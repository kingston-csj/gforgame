using System;
using System.Collections.Generic;

namespace Nova.Net.Socket
{
    public class MessageFactory
    {
        private Dictionary<int, Type> _id2Clazz = new();

        private Dictionary<Type, int> _clazz2Id = new();

        public void Register(int cmd, Type type)
        {
            _id2Clazz.TryAdd(cmd, type);
            _clazz2Id.TryAdd(type, cmd);
        }

        /// <summary>
        /// 获取消息的cmd
        /// </summary>
        /// <param name="type"></param>
        /// <returns></returns>
        public int GetMessageCmd(Type type)
        {
            return _clazz2Id[type];
        }

        /// <summary>
        /// 获取消息的类型
        /// </summary>
        /// <param name="cmd"></param>
        /// <returns></returns>
        public Type GetMessageType(int cmd)
        {
            return _id2Clazz[cmd];
        }

        public bool Contains(int cmd)
        {
            return _id2Clazz.ContainsKey(cmd);
        }
    }
    
}