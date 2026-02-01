using System;
using System.Collections.Generic;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using Nova.Commons.Util;
using Nova.Net.Socket;
using UnityEngine;

namespace Nova.Net.UnitySocket

{
    public class UnitySocket : IDisposable
    {
        #region 常量定义

        private const int DEFAULT_BUFFER_SIZE = 102400; // 100k，统一缓冲区大小
        private const int MESSAGE_HEADER_SIZE = 12; // 消息头大小（长度+索引+cmd）
        private const int RECEIVE_BUFFER_CHUNK = 8192; // 接收块大小（8k，平衡性能和GC）
        private const int SOCKET_CONNECT_TIMEOUT = 5000; // 连接超时时间（5秒）

        #endregion

        #region 私有字段

        private readonly Queue<EventArgs> eventQueue = new();
        private readonly object eventQueueLock = new();

        private string _Address;
        private string _host;
        private NetworkStream _networkStream;
        private int _port;

        private static readonly int MESSAGE_SIZE = 12;

        private readonly CancellationTokenSource _cts = new(); // 取消令牌，用于优雅关闭

        private ByteBuff _receiveBuff;
        private ByteBuff _sendBuff;
        private SocketState _state;

        private TcpClient _tcpClient;

        #endregion


        #region 公共事件

        /// <summary>
        ///     连接成功时触发的事件
        /// </summary>
        public event EventHandler<OpenEventArgs> OnOpen;

        /// <summary>
        ///     连接关闭时触发的事件
        /// </summary>
        public event EventHandler<CloseEventArgs> OnClose;

        /// <summary>
        ///     发生错误时触发的事件
        /// </summary>
        public event EventHandler<ErrorEventArgs> OnError;

        /// <summary>
        ///     收到消息时触发的事件
        /// </summary>
        public event EventHandler<MessageEventArgs> OnMessage;

        #endregion

        /// <summary>
        /// </summary>
        /// <param name="address">地址，比如 192.168.0.1:8000</param>
        public UnitySocket(string address)
        {
            // 初始化内存 100k
            _receiveBuff = new ByteBuff(DEFAULT_BUFFER_SIZE);
            // 初始化内存 100k
            _sendBuff = new ByteBuff(DEFAULT_BUFFER_SIZE);

            _Address = address;

            string[] arr = address.Split(":");
            if (arr.Length == 2)
            {
                _host = arr[0];
                _port = int.Parse(arr[1]);
            }
            else
            {
                Debug.LogError("无效的连接地址");
            }

            _state = SocketState.None;
        }

        /// <summary>
        /// </summary>
        /// <param name="host">192.168.0.1</param>
        /// <param name="port">8000</param>
        public UnitySocket(string host, int port)
        {
            // 初始化内存 1M
            _receiveBuff = new ByteBuff(DEFAULT_BUFFER_SIZE);
            _sendBuff = new ByteBuff(DEFAULT_BUFFER_SIZE);

            _Address = host + ":" + port;
            _host = host;
            _port = port;
            _state = SocketState.None;
        }

        private bool isOpening => _state == SocketState.Opened;


        public string Address => _Address;

        /// <summary>
        ///     异步连接到服务器
        /// </summary>
        public void ConnectAsync()
        {
            if (_state != SocketState.None)
            {
                Debug.LogWarning("当前Socket已处于连接/连接中状态，无需重复连接");
                return;
            }

            SocketManager.Instance.Add(this);
            Task.Run(ConnectTask);
        }

        /// <summary>
        ///     异步关闭连接
        /// </summary>
        public async Task CloseAsync()
        {
            if (_state == SocketState.Closed) return;

            _state = SocketState.Closed;
            _cts.Cancel(); // 取消所有异步任务

            // 优雅释放资源
            if (_networkStream != null)
            {
                try
                {
                    await _networkStream.FlushAsync();
                }
                catch
                {
                    /* 忽略关闭时的刷新异常 */
                }

                _networkStream.Dispose();
                _networkStream = null;
            }

            if (_tcpClient != null)
            {
                _tcpClient.Close();
                _tcpClient.Dispose();
                _tcpClient = null;
            }

            HandleClose();
            SocketManager.Instance.Remove(this);
        }

        private async Task ConnectTask()
        {
            try
            {
                _state = SocketState.Connecting;
                _tcpClient = new TcpClient();
                // 异步连接，带超时（避免无限等待）
                var connectTask = _tcpClient.ConnectAsync(_host, _port);
                if (await Task.WhenAny(connectTask, Task.Delay(SOCKET_CONNECT_TIMEOUT)) == connectTask)
                {
                    // 连接成功
                    _networkStream = _tcpClient.GetStream();
                    _state = SocketState.Opened;
                    HandleOpen();

                    // 启动异步接收（无轮询、无Sleep）
                    await ReceiveTask();
                }
                else
                {
                    // 连接超时
                    throw new TimeoutException($"连接服务器超时（{SOCKET_CONNECT_TIMEOUT}ms）：{this.Address}");
                }
            }
            catch (Exception e)
            {
                Debug.LogError($"连接失败: {e.Message}");
                HandleClose();
                Disconnect();
            }
        }

        public void Send(SocketDataFrame frame)
        {
            if (!isOpening || _networkStream == null)
            {
                Debug.LogError("未连接到服务器，无法发送数据");
                return;
            }

            try
            {
                _sendBuff.Reset();
                int frameSize = frame.rawData.Length + 12;
                _sendBuff.WriteInt(frameSize); // 写入总长度
                _sendBuff.WriteInt(frame.index); // 写入索引
                _sendBuff.WriteInt(frame.cmd); // 写入命令ID
                _sendBuff.WriteBytes(frame.rawData, frame.rawData.Length); // 写入消息体
                // 转换为最终字节数组
                byte[] rawData = _sendBuff.ToArray();
                _networkStream.Write(rawData, 0, rawData.Length);
                _networkStream.Flush(); // 立即发送
            }
            catch (Exception ex)
            {
                Debug.LogError($"发送数据失败：{ex.Message}");
            }
        }

        private async Task ReceiveTask()
        {
            byte[] buffer = new byte[RECEIVE_BUFFER_CHUNK]; // 8k缓冲区，减少GC
            while (isOpening && !_cts.Token.IsCancellationRequested)
            {
                try
                {
                    // 异步读取，无轮询，有数据才触发，CPU占用极低
                    int bytesRead = await _networkStream.ReadAsync(buffer, 0, buffer.Length, _cts.Token);
                    if (bytesRead == 0)
                    {
                        // 读取到0字节表示服务器关闭连接
                        Debug.Log("服务器主动关闭连接");
                        break;
                    }

                    // 写入接收缓冲区（线程安全）
                    lock (_receiveBuff)
                    {
                        _receiveBuff.WriteBytes(buffer, bytesRead);
                    }

                    // 解析所有完整包
                    _ReadAllPacks();
                }
                catch (OperationCanceledException)
                {
                    // 正常关闭，忽略
                    break;
                }
                catch (Exception e)
                {
                    Debug.LogError($"接收消息失败: {e.Message}");
                    HandleError(e);
                    break;
                }
            }

            // 接收异常/断开，关闭连接
            await CloseAsync();
        }

        /// <summary>
        /// 循环读取所有包
        /// </summary>
        private void _ReadAllPacks()
        {
            lock (_receiveBuff) // 保证解析时缓冲区不被修改
            {
                for (;;)
                {
                    if (_receiveBuff.ReadableBytes < 12)
                    {
                        // 如果不够一个字节，则跳出
                        break;
                    }

                    // 保存现场
                    _receiveBuff.MarkReadIndex();
                    // 整个封包的大小
                    int msgLength = _receiveBuff.ReadInt();
                    int index = _receiveBuff.ReadInt();
                    int cmd = _receiveBuff.ReadInt();
                    int bodySize = msgLength - MESSAGE_SIZE;
                    if (_receiveBuff.ReadableBytes >= bodySize)
                    {
                        // 可以读取整个包
                        SocketDataFrame socketPack = new();
                        socketPack.size = msgLength;
                        socketPack.index = index;
                        socketPack.cmd = cmd;
                        socketPack.rawData = _receiveBuff.ReadBytes(bodySize);
                        HandleMessage(socketPack);
                    }
                    else
                    {
                        // 半包，恢复读取位置，等待下一次数据
                        _receiveBuff.ResetRead();
                        break;
                    }
                }
            }
        }

        // 断开连接
        public void Disconnect()
        {
            _state = SocketState.Closed;

            if (_networkStream != null)
            {
                _networkStream.Close();
                _networkStream = null;
            }

            if (_tcpClient != null)
            {
                _tcpClient.Close();
                _tcpClient = null;
            }
        }

        private void HandleOpen()
        {
            HandleEventSync(new OpenEventArgs());
        }

        private void HandleMessage(SocketDataFrame data)
        {
            HandleEventSync(new MessageEventArgs(data));
        }

        private void HandleClose()
        {
            HandleEventSync(new CloseEventArgs());
        }

        private void HandleError(Exception exception)
        {
            HandleEventSync(new ErrorEventArgs(exception));
        }

        private void HandleEventSync(EventArgs eventArgs)
        {
            lock (eventQueueLock)
            {
                eventQueue.Enqueue(eventArgs);
            }
        }

        internal void Update()
        {
            EventArgs e;
            while (eventQueue.Count > 0)
            {
                lock (eventQueueLock)
                {
                    e = eventQueue.Dequeue();
                }

                if (e is CloseEventArgs)
                {
                    OnClose?.Invoke(this, e as CloseEventArgs);
                    SocketManager.Instance.Remove(this);
                }
                else if (e is OpenEventArgs)
                {
                    OnOpen?.Invoke(this, e as OpenEventArgs);
                }
                else if (e is MessageEventArgs)
                {
                    OnMessage?.Invoke(this, e as MessageEventArgs);
                }
                else if (e is ErrorEventArgs)
                {
                    OnError?.Invoke(this, e as ErrorEventArgs);
                }
            }
        }

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (disposing)
            {
                // 释放托管资源
                _cts.Cancel();
                _cts.Dispose();
                _ = CloseAsync(); // 异步关闭连接

                if (_receiveBuff != null)
                {
                    // 如果ByteBuff实现了IDisposable，这里释放
                    // _receiveBuff.Dispose();
                }

                if (_sendBuff != null)
                {
                    // _sendBuff.Dispose();
                }
            }
        }

        ~UnitySocket()
        {
            Dispose(false);
        }
    }
}