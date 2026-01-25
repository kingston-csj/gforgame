using System;
using System.Reflection;
using Nova.Commons.Util;
using Nova.Logger;
using UnityEngine;

namespace Nova.Net.Socket
{
    /// <summary>
    /// socket/websocket客户端基类
    /// </summary>
    public abstract class SocketClient
    {
        /// <summary>
        ///     是否输出日志
        /// </summary>
        public static bool _socketLog = true;

        /// <summary>
        ///     连接成功回调
        /// </summary>
        protected Action _connectSuccessCallback;

        /// <summary>
        ///     连接状态
        /// </summary>
        protected bool _isConnected;

        /// <summary>
        ///     请求消息序列号，自增长，客户端根据这个序列号来标记属于回调
        /// </summary>
        protected static int idCounter = 1;

        protected SocketIoDispatcher _dispatcher;

        protected SocketRuntimeEnvironment _runtimeEnvironment;

        /// <summary>
        ///     发送消息的数据结构
        /// </summary>
        protected SocketDataFrame _sendData = new();

        /// <summary>
        ///     异步连接到服务器
        ///     会先关闭已有连接，然后建立新连接
        /// </summary>
        public virtual void ConnectAsync(Action connectAsyncSuccessHandler)
        {
        }

        /// <summary>
        ///     发送请求，异步回调
        /// </summary>
        /// <typeparam name="RES">响应数据类型</typeparam>
        /// <param name="request">请求体数据</param>
        /// <param name="callback">响应回调函数</param>
        public void Send<RES>(Message request, Action<RES> callback)
            where RES : Response
        {
            //  获取协议类的cmd;
            int reqCmd = request.GetType().GetCustomAttribute<MessageMeta>().Cmd;
            _sendData.index = idCounter++;
            _sendData.cmd = reqCmd;
            _sendData.message = request;

            if (_socketLog)
            {
                string className = request.GetType().Name;
                LoggerUtil.Info("发送消息: " + className + " >> " + JsonUtility.ToJson(_sendData.message));
            }

            // 设置响应处理器
            Action<Message> reqCallBack = data =>
            {
                RES msg = (RES)data;
                callback(msg);
            };
            var call = new MessageCallback(typeof(RES), reqCallBack);
            CallbackMgr.Register(_sendData.index, call);

            Send(_sendData);
        }

        public virtual void Send(SocketDataFrame dataFrame)
        {
        }
    }
}