using System;
using System.Collections.Generic;
using Frame.Codec.Bean;
using Frame.Commons.Utils;
using Nova.Commons.Util;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 数组属性序列化
    /// 数组的元素可以是父类或抽象类
    /// 数组长度不能超过Short.MAX_VALUE，即最多65535
    /// </summary>
    public class ArrayCodec2 : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            // 读取数组长度
            int length = buff.ReadShort();
            var elementType = type.GetElementType();
            if (length <= 0)
            {
                return Array.CreateInstance(elementType, 0);
            }

            // 元素类型状态： 0代表基本类型或者元素类型一致，1代表元素类型不一致
            byte status = buff.ReadByte();
            // 创建数组并逐个解码元素
            Array array = Array.CreateInstance(elementType, length);
            Codec elementCodec = GetCodec(elementType);

            for (int i = 0; i < length; i++)
            {
                if (status == 1)
                {
                    int messageId = buff.ReadInt();
                    elementType = BeanCodecContext.MessageFactory.GetMessageType(messageId);
                    elementCodec = GetCodec(elementType);
                }

                object element = elementCodec.Decode(buff, elementType);
                array.SetValue(element, i);
            }

            return array;
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteShort(0);
                return;
            }

            Array array = value as Array ?? throw new ArgumentException("value 必须是数组类型");
            int length = array.Length;
            buff.WriteShort((short)length);

            if (length == 0)
            {
                return;
            }


            var elementType = array.GetType().GetElementType();
            byte status = 0;
            // 基本类型，写入状态码：0
            if (TypeUtil.IsPrimitiveOrString(elementType))
            {
                ISet<Type> set = new HashSet<Type>();
                foreach (object element in array)
                {
                    set.Add(element.GetType());
                }

                if (set.Count > 1 || elementType.IsInterface || elementType.IsAbstract)
                {
                    status = 1;
                }
            }

            buff.WriteByte(status);
            // 逐个编码元素
            Codec elementCodec = GetCodec(elementType);
            for (int i = 0; i < length; i++)
            {
                object element = array.GetValue(i);
                if (status == 1)
                {
                    buff.WriteInt(BeanCodecContext.MessageFactory.GetMessageCmd(element.GetType()));
                }

                elementCodec.Encode(buff, element);
            }
        }
    }
}