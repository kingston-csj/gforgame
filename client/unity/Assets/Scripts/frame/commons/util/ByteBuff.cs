namespace Nova.Commons.Util
{
    using System;
    using System.Text;

    /// <summary>
    /// 字节缓冲区工具类，用于高效读写各种基础数据类型（支持自动扩容、字节序转换、读写位置标记/重置）
    /// 核心特性：
    /// 1. 默认采用大端序（网络字节序），自动适配系统字节序（小端系统自动转换）；
    /// 2. 支持布尔、字节、短整数、整数、长整数、浮点数、双精度浮点数、字符串等类型的读写；
    /// 3. 自动扩容机制（不足时翻倍或直接扩至所需大小），无需手动管理缓冲区容量；
    /// 4. 支持读写位置标记与重置，适用于临时读写后回退场景；
    /// 5. 提供缓冲区压缩、清空、重置等辅助操作，优化内存使用。
    /// </summary>
    public class ByteBuff
    {
        /// <summary>
        /// 当前系统的字节序（true：小端序；false：大端序）
        /// </summary>
        private static readonly bool IsLittleEndian = BitConverter.IsLittleEndian;

        /// <summary>
        /// 是否需要进行字节序转换（系统为小端序且强制大端序时生效）
        /// </summary>
        private readonly bool needSwap;

        /// <summary>
        /// 内部存储字节数据的缓冲区
        /// </summary>
        private byte[] buffer;

        /// <summary>
        /// 下一次写入数据的起始位置（从0开始，随写入递增）
        /// </summary>
        private int writeIndex;

        /// <summary>
        /// 下一次读取数据的起始位置（从0开始，随读取递增）
        /// </summary>
        private int readIndex;

        /// <summary>
        /// 标记的读取位置（用于 MarkReadIndex/ResetReadIndex 临时保存读取位置）
        /// </summary>
        private int markedReadIndex;

        /// <summary>
        /// 标记的写入位置（用于 MarkWriteIndex/ResetWriteIndex 临时保存写入位置）
        /// </summary>
        private int markedWriteIndex;

        /// <summary>
        /// 缓冲区当前的有效容量（即已使用的最大长度，随写入动态更新）
        /// </summary>
        private int capacity;

        /// <summary>
        /// 构造函数，初始化字节缓冲区
        /// </summary>
        /// <param name="size">初始缓冲区大小（单位：字节），默认1024字节</param>
        /// <param name="forceBigEndian">是否强制使用大端序（网络字节序），默认true</param>
        /// <remarks>
        /// 1. 初始大小仅为默认值，缓冲区会根据写入数据量自动扩容，无需担心溢出；
        /// 2. 大端序是跨平台数据传输的标准字节序（如网络通信、文件存储），建议保持默认；
        /// 3. 小端序仅适用于单一平台的本地数据交互，不推荐跨环境使用。
        /// </remarks>
        public ByteBuff(int size = 1024, bool forceBigEndian = true)
        {
            buffer = new byte[size];
            writeIndex = 0;
            readIndex = 0;
            markedReadIndex = 0;
            markedWriteIndex = 0;
            capacity = size;
            // 系统是小端序且需要强制大端序时，才需要进行字节序转换
            needSwap = IsLittleEndian && forceBigEndian;
        }

        /// <summary>
        /// 根据配置自动转换字节序（仅当 needSwap 为 true 时执行反转）
        /// </summary>
        /// <param name="bytes">需要转换的字节数组</param>
        /// <returns>转换后的字节数组（无需转换时返回原数组，避免额外内存分配）</returns>
        public byte[] SwapBytes(byte[] bytes)
        {
            if (!needSwap) return bytes;
            Array.Reverse(bytes);
            return bytes;
        }

        #region 标记/重置位置方法

        /// <summary>
        /// 标记当前的读取位置，用于后续通过 ResetReadIndex 恢复
        /// </summary>
        /// <remarks>适用于临时读取数据（如解析协议头）后，需要回退到原位置继续读取的场景</remarks>
        public void MarkReadIndex()
        {
            markedReadIndex = readIndex;
        }

        /// <summary>
        /// 重置读取位置到最近一次 MarkReadIndex 标记的位置
        /// </summary>
        /// <remarks>若未调用过 MarkReadIndex，默认恢复到初始位置（0）</remarks>
        public void ResetReadIndex()
        {
            readIndex = markedReadIndex;
        }

        /// <summary>
        /// 标记当前的写入位置，用于后续通过 ResetWriteIndex 恢复
        /// </summary>
        /// <remarks>适用于临时写入数据（如预占协议长度字段）后，需要回退到原位置重新写入的场景</remarks>
        public void MarkWriteIndex()
        {
            markedWriteIndex = writeIndex;
        }

        /// <summary>
        /// 重置写入位置到最近一次 MarkWriteIndex 标记的位置
        /// </summary>
        /// <remarks>若未调用过 MarkWriteIndex，默认恢复到初始位置（0）</remarks>
        public void ResetWriteIndex()
        {
            writeIndex = markedWriteIndex;
        }

        #endregion

        #region 写入方法（自动更新写入位置，支持基础数据类型）

        /// <summary>
        /// 写入一个32位单精度浮点数（Float）到缓冲区
        /// </summary>
        /// <param name="value">要写入的浮点数（范围：±1.5×10^-45 至 ±3.4×10^38）</param>
        /// <remarks>
        /// 1. 占用4字节存储空间；
        /// 2. 自动处理字节序转换（小端系统强制大端时反转字节）；
        /// 3. 缓冲区不足时自动扩容。
        /// </remarks>
        public void WriteFloat(float value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个布尔值（Bool）到缓冲区
        /// </summary>
        /// <param name="value">要写入的布尔值（true/false）</param>
        public void WriteBool(bool value)
        {
            WriteByte(value ? (byte)1 : (byte)0);
        }

        /// <summary>
        /// 写入一个16位有符号短整数（Short）到缓冲区
        /// </summary>
        /// <param name="value">要写入的短整数（范围：-32768 至 32767）</param>
        public void WriteShort(short value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个32位有符号整数（Int）到缓冲区
        /// </summary>
        /// <param name="value">要写入的整数（范围：-2147483648 至 2147483647）</param>
        public void WriteInt(int value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个64位有符号长整数（Long）到缓冲区
        /// </summary>
        /// <param name="value">要写入的长整数（范围：-9223372036854775808 至 9223372036854775807）</param>
        public void WriteLong(long value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个字符串到缓冲区（先写长度，再写内容）
        /// </summary>
        /// <param name="value">要写入的字符串（支持空字符串）</param>
        /// <param name="encoding">字符串编码格式（默认UTF-8）</param>
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

        /// <summary>
        /// 写入一个64位双精度浮点数（Double）到缓冲区
        /// </summary>
        /// <param name="value">要写入的双精度浮点数（范围：±5.0×10^-324 至 ±1.7×10^308）</param>
        public void WriteDouble(double value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            WriteBytes(SwapBytes(bytes), bytes.Length);
        }

        /// <summary>
        /// 写入一个8位无符号字节（Byte）到缓冲区
        /// </summary>
        /// <param name="value">要写入的字节（范围：0 至 255）</param>
        public void WriteByte(byte value)
        {
            EnsureCapacity(writeIndex + 1);
            buffer[writeIndex] = value;
            writeIndex++;
            capacity = Math.Max(capacity, writeIndex);
        }

        /// <summary>
        /// 写入指定长度的字节数组到缓冲区
        /// </summary>
        /// <param name="bytes">要写入的字节数组（不可为null）</param>
        /// <param name="len">要写入的字节数（需≤bytes.Length，避免数组越界）</param>
        /// <remarks>
        /// 1. 采用 Buffer.BlockCopy 进行内存复制，效率高于 Array.Copy；
        /// 2. 自动检查缓冲区容量，不足时自动扩容；
        /// 3. 写入后自动更新 writeIndex 和 capacity；
        /// 4. 适用于批量写入二进制数据（如文件片段、协议体）。
        /// </remarks>
        /// <exception cref="ArgumentNullException">bytes 为 null 时抛出</exception>
        /// <exception cref="ArgumentOutOfRangeException">len 大于 bytes.Length 时抛出</exception>
        public void WriteBytes(byte[] bytes, int len)
        {
            if (bytes == null)
                throw new ArgumentNullException(nameof(bytes), "写入的字节数组不可为null");
            if (len < 0 || len > bytes.Length)
                throw new ArgumentOutOfRangeException(nameof(len), "写入长度超出字节数组范围");

            EnsureCapacity(writeIndex + len);
            Buffer.BlockCopy(bytes, 0, buffer, writeIndex, len);
            writeIndex += len;
            capacity = Math.Max(capacity, writeIndex);
        }

        #endregion

        #region 读取方法（自动更新读取位置，支持基础数据类型）

        /// <summary>
        /// 从缓冲区读取指定长度的字节数组
        /// </summary>
        /// <param name="count">要读取的字节数（需≥0）</param>
        /// <returns>包含读取数据的新字节数组（长度=count）</returns>
        /// <remarks>
        /// 1. 采用 Buffer.BlockCopy 进行内存复制，效率高于 Array.Copy；
        /// 2. 读取后自动递增 readIndex；
        /// 3. 适用于批量读取二进制数据（如协议体、文件片段）。
        /// </remarks>
        /// <exception cref="ArgumentOutOfRangeException">count 小于0时抛出</exception>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足 count 时抛出</exception>
        public byte[] ReadBytes(int count)
        {
            if (count < 0)
                throw new ArgumentOutOfRangeException(nameof(count), "读取长度不可为负数");
            if (ReadableBytes < count)
                throw new IndexOutOfRangeException(
                    $"读取字节数不足：需要 {count} 字节，可用 {ReadableBytes} 字节");

            byte[] result = new byte[count];
            Buffer.BlockCopy(buffer, readIndex, result, 0, count);
            readIndex += count;
            return result;
        }

        /// <summary>
        /// 从缓冲区读取一个布尔值（Bool）
        /// </summary>
        /// <returns>读取的布尔值（0x01 → true，0x00 → false）</returns>
        /// <remarks>
        /// 1. 读取1字节数据；
        /// 2. 读取后自动递增 readIndex；
        /// 3. 非0x00的字节都会被解析为 true（兼容异常数据场景）。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足1字节时抛出</exception>
        public bool ReadBool()
        {
            if (ReadableBytes < 1)
                throw new IndexOutOfRangeException(
                    $"读取布尔值失败：需要1字节，可用 {ReadableBytes} 字节");

            return ReadByte() != 0;
        }

        /// <summary>
        /// 从缓冲区读取一个16位有符号短整数（Short）
        /// </summary>
        /// <returns>读取的16位有符号短整数（范围：-32768 至 32767）</returns>
        /// <remarks>
        /// 1. 读取2字节数据，自动处理字节序转换；
        /// 2. 读取后自动递增 readIndex。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足2字节时抛出</exception>
        public short ReadShort()
        {
            if (ReadableBytes < 2)
                throw new IndexOutOfRangeException(
                    $"读取短整数失败：需要2字节，可用 {ReadableBytes} 字节");

            byte[] bytes = ReadBytes(2);
            return BitConverter.ToInt16(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从缓冲区读取一个32位单精度浮点数（Float）
        /// </summary>
        /// <returns>读取的32位单精度浮点数（范围：±1.5×10^-45 至 ±3.4×10^38）</returns>
        /// <remarks>
        /// 1. 读取4字节数据，自动处理字节序转换；
        /// 2. 读取后自动递增 readIndex；
        /// 3. 精度较低，适用于无需高精度的场景（如普通游戏参数）。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足4字节时抛出</exception>
        public float ReadFloat()
        {
            if (ReadableBytes < 4)
                throw new IndexOutOfRangeException(
                    $"读取浮点数失败：需要4字节，可用 {ReadableBytes} 字节");

            byte[] bytes = ReadBytes(4);
            return BitConverter.ToSingle(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从缓冲区读取一个64位双精度浮点数（Double）
        /// </summary>
        /// <returns>读取的64位双精度浮点数（范围：±5.0×10^-324 至 ±1.7×10^308）</returns>
        /// <remarks>
        /// 1. 读取8字节数据，自动处理字节序转换；
        /// 2. 读取后自动递增 readIndex；
        /// 3. 精度高于Float，适用于需要高精度的场景（如坐标、金额）。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足8字节时抛出</exception>
        public double ReadDouble()
        {
            if (ReadableBytes < 8)
                throw new IndexOutOfRangeException(
                    $"读取双精度浮点数失败：需要8字节，可用 {ReadableBytes} 字节");

            byte[] bytes = ReadBytes(8);
            return BitConverter.ToDouble(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从缓冲区读取一个32位有符号整数（Int）
        /// </summary>
        /// <returns>读取的32位有符号整数（范围：-2147483648 至 2147483647）</returns>
        /// <remarks>
        /// 1. 读取4字节数据，自动处理字节序转换；
        /// 2. 读取后自动递增 readIndex；
        /// 3. 适用于存储普通整数参数（如数量、ID、长度）。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足4字节时抛出</exception>
        public int ReadInt()
        {
            if (ReadableBytes < 4)
                throw new IndexOutOfRangeException(
                    $"读取整数失败：需要4字节，可用 {ReadableBytes} 字节");

            byte[] bytes = ReadBytes(4);
            return BitConverter.ToInt32(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从缓冲区读取一个64位有符号长整数（Long）
        /// </summary>
        /// <returns>读取的64位有符号长整数（范围：-9223372036854775808 至 9223372036854775807）</returns>
        /// <remarks>
        /// 1. 读取8字节数据，自动处理字节序转换；
        /// 2. 读取后自动递增 readIndex；
        /// 3. 适用于存储大范围整数（如时间戳、大金额、长ID）。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足8字节时抛出</exception>
        public long ReadLong()
        {
            if (ReadableBytes < 8)
                throw new IndexOutOfRangeException(
                    $"读取长整数失败：需要8字节，可用 {ReadableBytes} 字节");

            byte[] bytes = ReadBytes(8);
            return BitConverter.ToInt64(SwapBytes(bytes), 0);
        }

        /// <summary>
        /// 从缓冲区读取一个字符串（先读长度，再读内容）
        /// </summary>
        /// <param name="encoding">字符串编码格式（默认UTF-8）</param>
        /// <returns>读取的字符串（空字符串或有效内容）</returns>
        /// <remarks>
        /// 1. 读取流程：先读4字节长度 → 再读对应长度的字节内容 → 解码为字符串；
        /// 2. 若长度为0，直接返回空字符串；
        /// 3. 编码需与写入时一致（如写入用UTF-8，读取也需用UTF-8），否则会出现乱码。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">
        /// 1. 可读字节数不足4字节（长度字段）；
        /// 2. 可读字节数不足长度字段指定的字节数（内容字段）时抛出。
        /// </exception>
        public string ReadString(Encoding encoding = null)
        {
            encoding = encoding ?? Encoding.UTF8;

            if (ReadableBytes < 4)
                throw new IndexOutOfRangeException(
                    $"读取字符串长度失败：需要4字节，可用 {ReadableBytes} 字节");

            int length = ReadInt();
            if (length == 0) return string.Empty;

            if (ReadableBytes < length)
                throw new IndexOutOfRangeException(
                    $"读取字符串内容失败：需要 {length} 字节，可用 {ReadableBytes} 字节");

            byte[] bytes = ReadBytes(length);
            return encoding.GetString(bytes);
        }

        /// <summary>
        /// 从缓冲区读取一个8位无符号字节（Byte）
        /// </summary>
        /// <returns>读取的8位无符号字节（范围：0 至 255）</returns>
        /// <remarks>
        /// 1. 读取1字节数据，无字节序转换；
        /// 2. 读取后自动递增 readIndex；
        /// 适用于存储单个字节参数（如状态位、类型标识）。
        /// </remarks>
        /// <exception cref="IndexOutOfRangeException">可读字节数不足1字节时抛出</exception>
        public byte ReadByte()
        {
            if (ReadableBytes < 1)
                throw new IndexOutOfRangeException(
                    $"读取字节失败：需要1字节，可用 {ReadableBytes} 字节");

            byte result = buffer[readIndex];
            readIndex++;
            return result;
        }

        #endregion

        /// <summary>
        /// 确保缓冲区容量不小于指定的最小需求（不足时自动扩容）
        /// </summary>
        /// <param name="min">所需的最小容量（单位：字节）</param>
        /// <remarks>
        protected void EnsureCapacity(int min)
        {
            if (buffer.Length < min)
            {
                int newCapacity = Math.Max(buffer.Length * 2, min);
                byte[] newBuffer = new byte[newCapacity];
                Buffer.BlockCopy(buffer, 0, newBuffer, 0, capacity);
                buffer = newBuffer;
            }
        }

        /// <summary>
        /// 将缓冲区中已写入的有效数据转换为新的字节数组（从0到writeIndex）
        /// </summary>
        /// <returns>包含所有有效数据的字节数组（长度=writeIndex）</returns>
        /// <remarks>
        public byte[] ToArray()
        {
            int validLength = writeIndex;
            byte[] result = new byte[validLength];
            Buffer.BlockCopy(buffer, 0, result, 0, validLength);
            return result;
        }

        /// <summary>
        /// 重置所有读写位置到初始状态（0），不清空缓冲区数据
        /// </summary>
        /// <remarks>
        /// 1. 重置后可重新从缓冲区起始位置读写数据；
        /// 2. 缓冲区数据仍保留，若需清空数据请使用 Clear() 方法；
        /// 3. 适用于循环复用缓冲区（如频繁发送固定格式的协议）。
        /// </remarks>
        public void Reset()
        {
            readIndex = 0;
            writeIndex = 0;
            markedReadIndex = 0;
            markedWriteIndex = 0;
        }

        /// <summary>
        /// 重置读取位置到初始状态（0），并重置标记的读取位置
        /// </summary>
        /// <remarks>写入位置和缓冲区数据不受影响，适用于重新读取已写入的数据</remarks>
        public void ResetRead()
        {
            readIndex = 0;
            markedReadIndex = 0;
        }

        /// <summary>
        /// 重置写入位置到初始状态（0），并重置标记的写入位置
        /// </summary>
        /// <remarks>读取位置和缓冲区数据不受影响，适用于覆盖写入数据</remarks>
        public void ResetWrite()
        {
            writeIndex = 0;
            markedWriteIndex = 0;
        }

        /// <summary>
        /// 获取当前缓冲区中可读取的字节数（writeIndex - readIndex）
        /// </summary>
        public int ReadableBytes
        {
            get { return writeIndex - readIndex; }
        }

        /// <summary>
        /// 清空缓冲区所有数据，并重置所有读写位置到初始状态
        /// </summary>
        /// <remarks>
        /// 1. 清空逻辑：
        ///    a. 用 Array.Clear 清空缓冲区字节数据（设为0x00）；
        ///    b. 重置 readIndex、writeIndex、markedReadIndex、markedWriteIndex 为0；
        ///    c. 重置 capacity 为0；
        /// 2. 适用于缓冲区复用前的彻底清理（如切换不同类型的协议读写时）。
        /// </remarks>
        public void Clear()
        {
            Array.Clear(buffer, 0, buffer.Length);
            Reset();
            capacity = 0;
        }
    }
}