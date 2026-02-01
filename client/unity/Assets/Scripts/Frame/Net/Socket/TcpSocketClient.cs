using System;
using Nova.Commons.Util;
using Nova.Logger;
using Nova.Net.UnitySocket;

namespace Nova.Net.Socket
{
    public class TcpSocketClient : SocketClient, IDisposable
    {
        /// <summary>
        /// 原生socket
        /// </summary>
        private UnitySocket.UnitySocket _socket;

        private bool _isConnected = false;

        /// <summary>
        /// 连接地址，例如：192.168.1.100:8080
        /// </summary>
        private string address;

        public TcpSocketClient(string address, SocketRuntimeEnvironment runtimeEnvironment)
        {
            this.address = address;
            _isConnected = false;
            this._runtimeEnvironment = runtimeEnvironment;
            _dispatcher = new DefaultSocketIoDispatcher(runtimeEnvironment);
        }

        /// <summary>
        /// 异步连接到服务器
        /// 会先关闭已有连接，然后建立新连接
        /// </summary>
        public override void ConnectAsync(Action connectAsyncSuccessHandler)
        {
            _connectSuccessCallback = connectAsyncSuccessHandler;
            CloseAsync();
            _socket = new UnitySocket.UnitySocket(this.address);
            _socket.OnOpen += OnOpen;
            _socket.OnMessage += OnMessage;
            _socket.OnClose += OnClose;
            _socket.OnError += OnError;
            _socket.ConnectAsync();
            LoggerUtil.Info("连接服务器 @ " + this.address);
        }

        /// <summary>
        ///     发送二进制消息
        /// </summary>
        public override void Send(SocketDataFrame frame)
        {
            if (_socket != null && _isConnected)
            {
                frame.rawData = _runtimeEnvironment.MessageCodec.Encode(frame.message);
                _socket.Send(frame);
            }
        }

        protected void OnOpen(object sender, OpenEventArgs e)
        {
            _isConnected = true;
            _OnOpen();
        }

        protected virtual void _OnOpen()
        {
            _connectSuccessCallback?.Invoke();
        }

        /// <summary>
        /// WebSocket连接关闭回调
        /// </summary>
        protected void OnClose(object sender, CloseEventArgs e)
        {
            _isConnected = false;
            _OnClose();
        }

        protected virtual void _OnClose()
        {
        }

        /// <summary>
        ///     接收消息处理
        ///     解析消息并调用对应的回调函数
        /// </summary>
        protected void OnMessage(object sender, MessageEventArgs e)
        {
            if (_isConnected)
            {
                byte[] body = e.DataFrame.rawData;
                if (!_runtimeEnvironment.MessageFactory.Contains(e.DataFrame.cmd))
                {
                    LoggerUtil.Error($"未注册的消息类型: {e.DataFrame.cmd}");
                    return;
                }

                Type type = _runtimeEnvironment.MessageFactory.GetMessageType(e.DataFrame.cmd);
                e.DataFrame.message = (Message)_runtimeEnvironment.MessageCodec.Decode(type, body);
                if (e.DataFrame.message == null)
                {
                    LoggerUtil.Error($"数据解析错误: {type}");
                }

                if (_socketLog)
                {
                    LoggerUtil.Info("接收消息: " + e.DataFrame.msgJson);
                }

                if (!_runtimeEnvironment.MessageFactory.Contains(e.DataFrame.cmd))
                {
                    LoggerUtil.Error($"未注册消息类型:{e.DataFrame.cmd}");
                    return;
                }

                _dispatcher.OnMessage(e.DataFrame);
            }
        }

        /// <summary>
        /// WebSocket错误回调
        /// </summary>
        protected void OnError(object sender, ErrorEventArgs e)
        {
            _isConnected = false;
            _OnError(e.message);
        }

        protected virtual void _OnError(Exception e)
        {
        }


        /// <summary>
        /// 关闭WebSocket连接
        /// </summary>
        public void CloseAsync()
        {
            if (_socket != null)
            {
                _socket.OnOpen -= OnOpen;
                _socket.OnMessage -= OnMessage;
                _socket.OnClose -= OnClose;
                _socket.OnError -= OnError;
                _socket.CloseAsync();
                _socket = null;
            }
        }


        public void Dispose()
        {
            CloseAsync();
            _connectSuccessCallback = null;

            if (_socket != null)
            {
                _socket.Disconnect();
            }
        }
    }
}