using System;
using System.Text;

namespace Nova.Commons.Util
{
    /// <summary>
    /// 字节数组处理类，用于处理二进制数据的读写操作
    /// 支持基本数据类型（int, float, long, string）的序列化和反序列化
    /// 典型使用场景:
    /// 1. 网络通信：序列化数据包
    /// 2. 二进制文件读写
    /// 3. 内存数据缓存
    /// </summary>
    public class ByteBuff
    {
        /// <summary>
        /// 系统字节序
        /// </summary>
        private static readonly bool IsLittleEndian = BitConverter.IsLittleEndian;

        /// <summary>
        /// 是否需要转换字节序（默认使用大端序）
        /// </summary>
        private readonly bool needSwap;

        /// <summary>
        /// 内部字节缓冲区
        /// </summary>
        private byte[] buffer;

        /// <summary>
        /// 写入位置
        /// </summary>
        private int writeIndex;

        /// <summary>
        /// 读取位置
        /// </summary>
        private int readIndex;

        /// <summary>
        /// 标记的读取位置
        /// </summary>
        private int markedReadIndex;

        /// <summary>
        /// 标记的写入位置
        /// </summary>
        private int markedWriteIndex;

        /// <summary>
        /// 构造函数
        /// </summary>
        /// <param name="size">初始缓冲区大小</param>
        /// <param name="forceBigEndian">是否强制使用大端序</param>
        public ByteBuff(int size = 1024, bool forceBigEndian = true)
        {
            buffer = new byte[size];
            writeIndex = 0;
            readIndex = 0;
            markedReadIndex = 0;
            markedWriteIndex = 0;
            needSwap = IsLittleEndian && forceBigEndian;
        }

        /// <summary>
        /// 转换字节序
        /// </summary>
        public byte[] SwapBytes(byte[] bytes)
        {
            if (!needSwap || bytes == null || bytes.Length <= 1) return bytes;
            Array.Reverse(bytes);
            return bytes;
        }

        #region 标记/重置位置方法（统一命名，匹配Socket代码调用）

        public void MarkReadIndex()
        {
            markedReadIndex = readIndex;
        }

        /// <summary>
        /// 重置读取位置到最近标记的位置
        /// </summary>
        public void ResetReadIndex()
        {
            readIndex = markedReadIndex;
            markedReadIndex = 0;
        }

        /// <summary>
        /// 标记当前写入位置
        /// </summary>
        public void MarkWriteIndex()
        {
            markedWriteIndex = writeIndex;
        }

        /// <summary>
        /// 重置写入位置到最近标记的位置
        /// </summary>
        public void ResetWriteIndex()
        {
            writeIndex = markedWriteIndex;
        }

        #endregion

        #region 写入方法（修复扩容和容量计算）

        /// <summary>
        /// 写入一个浮点数到字节数组中
        /// </summary>
        /// <param name="value">要写入的浮点数值</param>
        public void WriteFloat(float value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个布尔值到字节数组中
        /// </summary>
        public void WriteBool(bool value)
        {
            WriteByte(value ? (byte)1 : (byte)0);
        }

        public void WriteShort(short value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个整数到字节数组中
        /// </summary>
        /// <param name="value">要写入的整数值</param>
        public void WriteInt(int value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个长整数到字节数组中
        /// </summary>
        /// <param name="value">要写入的长整数值</param>
        public void WriteLong(long value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个字符串到字节数组中
        /// </summary>
        public void WriteString(string value, Encoding encoding = null)
        {
            encoding = encoding ?? Encoding.UTF8;
            if (string.IsNullOrEmpty(value))
            {
                WriteInt(0);
                return;
            }

            byte[] bytes = encoding.GetBytes(value);
            WriteInt(bytes.Length);
            WriteBytes(bytes, bytes.Length);
        }

        public void WriteDouble(double value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个字节到缓冲区
        /// </summary>
        /// <param name="value">要写入的字节值</param>
        public void WriteByte(byte value)
        {
            EnsureCapacity(writeIndex + 1);
            buffer[writeIndex] = value;
            writeIndex++;
        }

        /// <summary>
        /// 将字节数组写入缓冲区
        /// </summary>
        public void WriteBytes(byte[] bytes, int len)
        {
            if (bytes == null || len <= 0) return; // 防御性判断
            EnsureCapacity(writeIndex + len);
            Buffer.BlockCopy(bytes, 0, buffer, writeIndex, len);
            writeIndex += len;
        }

        #endregion

        #region 读取方法（保持原有逻辑，增加防御性判断）

        /// <summary>
        /// 从缓冲区读取指定数量的字节
        /// </summary>
        /// <param name="count">要读取的字节数</param>
        /// <returns>包含读取数据的新字节数组</returns>
        /// <remarks>
        /// 1. 会检查是否越界
        /// 2. 使用Buffer.BlockCopy进行高效的内存复制
        /// 3. 自动更新读取位置
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">当读取位置超出有效数据范围时抛出</exception>
        public byte[] ReadBytes(int count)
        {
            if (count <= 0) return Array.Empty<byte>();
            if (ReadableBytes < count)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read. Required: {count}, Available: {ReadableBytes}");

            byte[] result = new byte[count];
            Buffer.BlockCopy(buffer, readIndex, result, 0, count);
            readIndex += count;
            return result;
        }

        /// <summary>
        /// 从字节数组中读取一个布尔值
        /// </summary>
        public bool ReadBool()
        {
            if (ReadableBytes < 1)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read bool. Required: 1, Available: {ReadableBytes}");

            return ReadByte() != 0;
        }

        /// <summary>
        /// 从字节数组中读取一个Short值
        /// </summary>
        public short ReadShort()
        {
            if (ReadableBytes < 2)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read short. Required: 2, Available: {ReadableBytes}");

            byte[] bytes = ReadBytes(2);
            return BitConverter.ToInt16(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从字节数组中读取一个Float值
        /// </summary>
        public float ReadFloat()
        {
            if (ReadableBytes < 4)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read float. Required: 4, Available: {ReadableBytes}");

            byte[] bytes = ReadBytes(4);
            return BitConverter.ToSingle(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从字节数组中读取一个Double值
        /// </summary>
        public double ReadDouble()
        {
            if (ReadableBytes < 8)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read double. Required: 8, Available: {ReadableBytes}");

            byte[] bytes = ReadBytes(8);
            return BitConverter.ToDouble(SwapBytes(bytes), 0);
        }

        public int ReadInt()
        {
            if (ReadableBytes < 4)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read int. Required: 4, Available: {ReadableBytes}");

            byte[] bytes = ReadBytes(4);
            return BitConverter.ToInt32(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从字节数组中读取一个长整数
        /// </summary>
        /// <returns>读取的长整数值</returns>
        public long ReadLong()
        {
            if (ReadableBytes < 8)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read long. Required: 8, Available: {ReadableBytes}");

            byte[] bytes = ReadBytes(8);
            return BitConverter.ToInt64(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从字节数组中读取一个字符串
        /// 先读取字符串长度，再读取字符串内容
        /// </summary>
        /// <param name="encoding">字符串编码格式，默认UTF8</param>
        /// <returns>读取的字符串</returns>
        public string ReadString(Encoding encoding = null)
        {
            encoding = encoding ?? Encoding.UTF8;

            if (ReadableBytes < 4)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read string length. Required: 4, Available: {ReadableBytes}");

            int length = ReadInt();
            if (length == 0) return string.Empty;

            if (ReadableBytes < length)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read string content. Required: {length}, Available: {ReadableBytes}");

            byte[] bytes = ReadBytes(length);
            return encoding.GetString(bytes);
        }

        /// <summary>
        /// 从缓冲区读取一个字节
        /// </summary>
        /// <returns>读取的字节值</returns>
        /// <exception cref="IndexOutOfRangeException">当读取位置超出有效数据范围时抛出</exception>
        public byte ReadByte()
        {
            if (ReadableBytes < 1)
                throw new IndexOutOfRangeException(
                    $"Not enough bytes to read byte. Required: 1, Available: {ReadableBytes}");

            byte result = buffer[readIndex];
            readIndex++;
            return result;
        }

        #endregion

        #region 核心修复：扩容方法

        /// <summary>
        /// 确保缓冲区容量足够
        /// </summary>
        /// <param name="min">所需的最小容量</param>
        protected void EnsureCapacity(int min)
        {
            if (buffer.Length >= min) return;

            // 扩容策略：翻倍或直接到所需大小，取最大值
            int newCapacity = Math.Max(buffer.Length * 2, min);
            byte[] newBuffer = new byte[newCapacity];

            // 复制已写入的有效数据（writeIndex）
            if (writeIndex > 0)
            {
                Buffer.BlockCopy(buffer, 0, newBuffer, 0, writeIndex);
            }

            // 替换缓冲区，保留有效数据
            buffer = newBuffer;
        }

        #endregion

        #region 辅助方法（修复命名和逻辑）

        /// <summary>
        /// 获取可读取的字节数
        /// </summary>
        public int ReadableBytes => writeIndex - readIndex;

        /// <summary>
        /// 将有效数据转换为新数组
        /// </summary>
        public byte[] ToArray()
        {
            byte[] result = new byte[writeIndex];
            Buffer.BlockCopy(buffer, 0, result, 0, writeIndex);
            return result;
        }

        /// <summary>
        /// 重置读写位置
        /// </summary>
        public void Reset()
        {
            readIndex = 0;
            writeIndex = 0;
            markedReadIndex = 0;
            markedWriteIndex = 0;
        }

        /// <summary>
        /// 压缩缓冲区，移除已读取的无效数据
        /// </summary>
        public void Compress()
        {
            if (readIndex <= 0) return;

            int readable = ReadableBytes;
            if (readable > 0)
            {
                // 将未读取的数据移到缓冲区头部
                Buffer.BlockCopy(buffer, readIndex, buffer, 0, readable);
            }

            // 更新读写索引
            writeIndex = readable;
            readIndex = 0;
            markedReadIndex = 0;
            // 下面这行，可能有问题，避免在使用合并后，再恢复markedWriteIndex
            markedWriteIndex = Math.Min(markedWriteIndex, writeIndex);
        }

        /// <summary>
        /// 清空缓冲区（保留物理容量，仅重置数据）
        /// </summary>
        public void Clear()
        {
            // 无需清空数组，重置索引即可（性能更高）
            Reset();
        }

        #endregion
    }
}