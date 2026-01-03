using System;

namespace Nova.Codec
{
    /// <summary>
    /// 由于二进制是一种字节流，若服务器与客户端的消息定义不一致，则无法进行序列化
    /// 对于一些客户端自定义的字段，服务器没有相应的定义，为了防止序列化错误，
    /// 可以使用该特性标记需要忽略的字段
    /// </summary>
    [AttributeUsage(AttributeTargets.Field)]
    public class FieldIgnore : Attribute
    {
    }
}