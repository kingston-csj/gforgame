using System;
using System.Reflection;
using Game.Net;
using Game.Net.Message;
using Nova.Codec;
using Nova.Commons.Util;
using Nova.Logger;
using Nova.Net.Socket;
using UnityEngine;

namespace Game.Core
{
    public class WebManager
    {
        private static WebManager _instance;

        private bool connecting = false;

        public static WebManager Instance
        {
            get
            {
                if (_instance == null)
                {
                    _instance = new WebManager();
                }

                return _instance;
            }
        }

        /// <summary>
        /// 连接到服务器
        /// </summary>
        /// <param name="url">游戏服务器（网关）地址</param>
        /// <param name="callback">连接成功回调</param>
        public void ConnectToServer(string url, Action callback)
        {
            if (connecting)
            {
                return;
            }
            SocketRuntimeEnvironment runtimeEnvironment =
                new SocketRuntimeEnvironment(typeof(MessageRouter), new JsonCodec(), new MessageFactory());
            // 自动注册所有的消息类型
            foreach (Type item in ClassScanner.ListClassesWithAttribution<MessageMeta>())
            {
                // 获得class对应MessageMeta特性的cmd
                MessageMeta messageMeta = item.GetCustomAttribute(typeof(MessageMeta)) as MessageMeta;
                int cmd = messageMeta.Cmd;
                runtimeEnvironment.MessageFactory.Register(cmd, item);
            }
            connecting = true;
            SocketClient socketClient = new TcpSocketClient(url, runtimeEnvironment);
            AppContext.socketClient = socketClient;

            socketClient.ConnectAsync(() =>
            {
                // 连接成功
                LoggerUtil.Info("连接成功");
                callback();
            });
        }
    }
}