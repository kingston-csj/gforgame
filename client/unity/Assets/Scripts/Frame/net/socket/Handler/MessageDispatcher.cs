using System;
using System.Collections.Generic;
using Nova.Logger;

namespace Nova.Net.Socket
{
    public class MessageDispatcher
    {
        static Dictionary<int, MessageCallback> _handlers = new();

        static Dictionary<int, Type> _handlerTypes = new();

        /// <summary>
        /// 注册消息处理器
        /// </summary>
        /// <param name="cmd"></param>
        /// <param name="handler"></param>
        public static void Register(int cmd, MessageCallback handler)
        {
            _handlers.Add(cmd, handler);
            _handlerTypes.Add(cmd, handler.type);
        }

        public static void Dispatch(int cmd, Message msg)
        {
            if (_handlers.ContainsKey(cmd))
            {
                MessageCallback handler;
                _handlers.TryGetValue(cmd, out handler);
                if (handler == null)
                {
                    LoggerUtil.Error("未注册数据监听 cmd:" + cmd);
                }
                handler?.callback(msg);
            }
        }
    }
    
}