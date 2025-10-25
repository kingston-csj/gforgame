using System;
using System.Reflection;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 字段编码元数据
    /// 包含字段信息和对应的编码解码器
    /// </summary>
    public class FieldCodecMeta
    {
        public FieldInfo Field { get; }
        public Codec Codec { get; }

        public FieldCodecMeta(FieldInfo field, Codec codec)
        {
            Field = field ?? throw new ArgumentNullException(nameof(field));
            Codec = codec ?? throw new ArgumentNullException(nameof(codec));
        }

        public static FieldCodecMeta Create(FieldInfo field, Codec codec)
        {
            return new FieldCodecMeta(field, codec);
        }
    }
    
}