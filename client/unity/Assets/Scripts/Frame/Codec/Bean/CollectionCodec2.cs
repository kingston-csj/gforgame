using System;
using System.Collections;
using System.Collections.Generic;
using System.Reflection;
using Frame.Codec.Bean;
using Frame.Commons.Utils;
using Nova.Commons.Util;
using Nova.Net.Socket;

namespace Nova.Codec.bean
{
    /// <summary>
    /// 集合编解码器（支持List/HashSet）
    /// 数组长度不能超过Short.MAX_VALUE，即最多65535
    /// 数组的元素可以是父类或抽象类
    /// </summary>
    public class CollectionCodec2 : Codec
    {
        public override object Decode(ByteBuff buff, Type type)
        {
            // 读取集合长度
            int length = buff.ReadShort();
            if (length <= 0)
            {
                return Activator.CreateInstance(type);
            }

            // 元素类型状态： 0代表基本类型或者元素类型一致，1代表元素类型不一致
            byte status = buff.ReadByte();

            // 获取集合元素类型 
            Type elementType = type.GetGenericArguments()[0];
            Codec elementCodec = GetCodec(elementType);

            // 创建集合并逐个添加元素
            ICollection collection = (ICollection)Activator.CreateInstance(type);
            MethodInfo addMethod = type.GetMethod("Add", new[] { elementType });

            for (int i = 0; i < length; i++)
            {
                if (status == 1)
                {
                    int messageId = buff.ReadInt();
                    elementType = BeanCodecContext.MessageFactory.GetMessageType(messageId);
                    elementCodec = GetCodec(elementType);
                }

                object element = elementCodec.Decode(buff, elementType);
                addMethod.Invoke(collection, new[] { element });
            }

            return collection;
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                buff.WriteShort(0);
                return;
            }

            ICollection collection = value as ICollection
                                     ?? throw new ArgumentException(
                                         $"value 必须是集合类型（实现 ICollection 接口），当前类型：{value.GetType().FullName}");

            buff.WriteShort((short)collection.Count);
            if (collection.Count == 0)
            {
                return;
            }

            // 获取集合元素类型（ 
            Type elementType = GetCollectionElementType(value.GetType());
            byte status = 0;
            // 基本类型，写入状态码：0
            if (TypeUtil.IsPrimitiveOrString(elementType))
            {
                ISet<Type> set = new HashSet<Type>();
                foreach (object element in collection)
                {
                    set.Add(element.GetType());
                }

                if (set.Count > 1 || elementType.IsInterface || elementType.IsAbstract)
                {
                    status = 1;
                }
            }

            buff.WriteByte(status);

            Codec elementCodec = GetCodec(elementType);

            // 遍历集合元素编码（非泛型ICollection需强制转换为IEnumerable）
            foreach (object element in (IEnumerable)collection)
            {
                if (status == 1)
                {
                    buff.WriteInt(BeanCodecContext.MessageFactory.GetMessageCmd(element.GetType()));
                }

                elementCodec.Encode(buff, element);
            }
        }

        /// <summary>
        /// 解析集合的元素类型（泛型集合）
        /// </summary>
        private Type GetCollectionElementType(Type collectionType)
        {
            if (!collectionType.IsGenericType)
            {
                throw new NotSupportedException("仅支持泛型集合类型");
            }

            Type[] genericArgs = collectionType.GetGenericArguments();
            if (genericArgs.Length != 1)
            {
                throw new NotSupportedException("仅支持单泛型参数的集合类型");
            }

            return genericArgs[0];
        }
    }
}