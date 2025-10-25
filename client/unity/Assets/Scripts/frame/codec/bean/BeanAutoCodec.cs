using System;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    
    /// <summary>
    /// 基于结构体的消息体编码解码（自动适配基础类型、数组、集合、自定义类型）
    /// </summary>
    public class BeanAutoCodec : MessageCodec
    {
        /// <summary>
        /// 默认写入缓冲区大小（1M）
        /// </summary>
        private const int DefaultWriteBuffSize = 1024 * 1024;

        /// <summary>
        /// 线程本地缓冲区（避免线程安全问题）
        /// </summary>
        private readonly ByteBuff _sendBuff;

        private ByteBuff _receiveBuff;

        /// <summary>
        /// 构造函数（指定缓冲区大小）
        /// </summary>
        /// <param name="writeBuffSize">编码单个消息的最大缓冲区大小，超出会抛出 BufferOverflowException</param>
        public BeanAutoCodec(int readBuffSize, int writeBuffSize)
        {
            if (readBuffSize <= 0 || writeBuffSize <= 0)
            {
                throw new ArgumentOutOfRangeException(nameof(writeBuffSize), "缓冲区大小必须大于0");
            }

            _receiveBuff = new ByteBuff(readBuffSize);
            _sendBuff = new ByteBuff(writeBuffSize);
        }


        /// <summary>
        /// 构造函数（使用默认缓冲区大小1M）
        /// </summary>
        public BeanAutoCodec() : this(DefaultWriteBuffSize, DefaultWriteBuffSize)
        {
        }

        /// <summary>
        /// 解码：byte[] → 目标Bean类型
        /// </summary>
        /// <param name="msgType">目标Bean类型</param>
        /// <param name="body">待解码的字节数组（完整包体）</param>
        public object Decode(Type msgType, byte[] body)
        {
            if (msgType == null)
                throw new ArgumentNullException(nameof(msgType));
            if (body == null || body.Length == 0)
                return Activator.CreateInstance(msgType);

            try
            {
                // 字节数组包装为ByteBuff（读模式）
                _receiveBuff.Reset();
                _receiveBuff.WriteBytes(body, body.Length);
                _receiveBuff.ResetRead(); // 重置读取位置到开头

                // 获取对应类型的编解码器并解码
                Codec codec = Codec.GetCodec(msgType);
                return codec.Decode(_receiveBuff, msgType);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"解码失败，消息类型：{msgType.FullName}", ex);
                return null;
            }
        }

        /// <summary>
        /// 编码：Bean → byte[]
        /// </summary>
        /// <param name="message">待编码的Bean实例</param>
        public byte[] Encode(object message)
        {
            if (message == null)
                throw new ArgumentNullException(nameof(message));

            // 清空发送缓冲区
            _sendBuff.Clear();
            try
            {
                // 获取Bean类型的编解码器并编码
                Codec codec = Codec.GetCodec(message.GetType());
                codec.Encode(_sendBuff, message);
                // 转换为字节数组返回（仅包含有效数据）
                return _sendBuff.ToArray();
            }
            catch (Exception ex)
            {
                Console.WriteLine($"编码失败，消息类型：{message.GetType().FullName}", ex);
                return Array.Empty<byte>();
            }
        }
    }
}