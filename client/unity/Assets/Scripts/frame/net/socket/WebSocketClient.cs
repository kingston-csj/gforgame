using System;
using Nova.Commons.Util;
using Nova.Logger;
using UnityEngine;
using UnityWebSocket;
using CloseEventArgs = UnityWebSocket.CloseEventArgs;
using ErrorEventArgs = UnityWebSocket.ErrorEventArgs;
using MessageEventArgs = UnityWebSocket.MessageEventArgs;
using OpenEventArgs = UnityWebSocket.OpenEventArgs;

namespace Nova.Net.Socket
{
    /// <summary>
    ///     WebSocket通信客户端
    ///     负责处理WebSocket连接、消息收发、回调管理等核心功能
    /// </summary>
    public class WebSocketClient : SocketClient, IDisposable
    {
        /// <summary>
        ///     WebSocket实例
        /// </summary>
        private WebSocket _socket;

        /// <summary>
        /// 连接地址，例如：ws://192.168.1.100:8080/ws
        /// </summary>
        private string address;


        /// <summary>
        ///     构造函数：初始化基本数据结构
        /// </summary>
        public WebSocketClient(string address, SocketRuntimeEnvironment runtimeEnvironment)
        {
            this.address = address;
            _isConnected = false;
            _sendData = new SocketDataFrame();
            this._runtimeEnvironment = runtimeEnvironment;
            _dispatcher = new DefaultSocketIoDispatcher(runtimeEnvironment);
        }

        /// <summary>
        ///     异步连接到服务器
        ///     会先关闭已有连接，然后建立新连接
        /// </summary>
        /// <param name="connectAsyncSuccessHandler"> 连接成功回调函数 </param>
        public override void ConnectAsync(Action connectAsyncSuccessHandler)
        {
            _connectSuccessCallback = connectAsyncSuccessHandler;
            CloseAsync();
            _socket = new WebSocket(address);
            _socket.OnOpen += OnOpen;
            _socket.OnMessage += OnMessage;
            _socket.OnClose += OnClose;
            _socket.OnError += OnError;
            _socket.ConnectAsync();
            LoggerUtil.Info("连接中... ip: " + address);
        }

        /// <summary>
        ///     WebSocket连接成功回调
        /// </summary>
        protected void OnOpen(object sender, OpenEventArgs e)
        {
            _isConnected = true;
            _OnOpen();
        }

        protected virtual void _OnOpen()
        {
            _connectSuccessCallback?.Invoke();
        }

        /// 接收缓冲区
        private ByteBuff _receiveBuff = new ByteBuff();

        /// <summary>
        ///     接收消息处理
        ///     解析消息并调用对应的回调函数
        /// </summary>
        protected void OnMessage(object sender, MessageEventArgs e)
        {
            if (_isConnected)
            {
                SocketDataFrame dataFrame = null;
                if (!_runtimeEnvironment.UsedBinaryFrame)
                {
                    dataFrame = JsonUtility.FromJson<SocketDataFrame>(e.Data);
                }
                else
                {
                    dataFrame = new();
                    _receiveBuff.Reset();
                    _receiveBuff.WriteBytes(e.RawData, e.RawData.Length);
                    // 整包长度:包头(12个字节)+包体(不定长)
                    dataFrame.size = _receiveBuff.ReadInt();
                    dataFrame.index = _receiveBuff.ReadInt();
                    dataFrame.cmd = _receiveBuff.ReadInt();
                    int bodySize = dataFrame.size - 12;
                    byte[] body = _receiveBuff.ReadBytes(bodySize);
                    if (!_runtimeEnvironment.MessageFactory.Contains(dataFrame.cmd))
                    {
                        LoggerUtil.Error($"未注册的消息类型: {dataFrame.cmd}");
                        return;
                    }

                    Type type = _runtimeEnvironment.MessageFactory.GetMessageType(dataFrame.cmd);
                    dataFrame.data = (Message)_runtimeEnvironment.MessageCodec.Decode(type, body);
                    if (dataFrame.data == null)
                    {
                        LoggerUtil.Error($"数据解析错误: {type}");
                    }
                }

                if (_socketLog)
                {
                    LoggerUtil.Info("接收消息: " + dataFrame.msg);
                }

                if (!_runtimeEnvironment.MessageFactory.Contains(dataFrame.cmd))
                {
                    LoggerUtil.Error($"未注册消息类型:{dataFrame.cmd}");
                    return;
                }

                _dispatcher.OnMessage(dataFrame);
            }
        }

        /// <summary>
        ///     WebSocket连接关闭回调
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
        ///     WebSocket错误回调
        /// </summary>
        protected void OnError(object sender, ErrorEventArgs e)
        {
            _isConnected = false;
            _OnError();
        }

        protected void _OnError()
        {
        }

        // 发送缓冲区
        private ByteBuff _sendBuff = new();

        /// <summary>
        ///     发送二进制消息
        /// </summary>
        public override void Send(SocketDataFrame frame)
        {
            if (_socket != null && _isConnected)
            {
                // 发送二进制消息
                if (_runtimeEnvironment.UsedBinaryFrame)
                {
                    byte[] body = _runtimeEnvironment.MessageCodec.Encode(frame.data);
                    _sendBuff.Reset();
                    int frameSize = body.Length + 12;
                    _sendBuff.WriteInt(frameSize); // 写入总长度
                    _sendBuff.WriteInt(frame.index); // 写入索引
                    _sendBuff.WriteInt(frame.cmd); // 写入命令ID
                    _sendBuff.WriteBytes(body, body.Length); // 写入消息体
                    _socket.SendAsync(_sendBuff.ToArray());
                }
                else
                {
                    // 发送文本消息
                    frame.msg = JsonUtility.ToJson(frame.data);
                    string data = JsonUtility.ToJson(frame);
                    Send(data);
                }
            }
            else
            {
                LoggerUtil.Error("发送消息失败,socket为空或未连接");
            }
        }

        public void Send(string message)
        {
            if (_socket != null && _isConnected)
            {
                _socket.SendAsync(message);
            }
            else
            {
                LoggerUtil.Error("发送消息失败,socket为空或未连接");
            }
        }


        /// <summary>
        ///     关闭WebSocket连接
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

        /// <summary>
        ///     销毁对象时清理资源
        /// </summary>
        public void Dispose()
        {
            CloseAsync();
            _connectSuccessCallback = null;
        }
    }
    
}