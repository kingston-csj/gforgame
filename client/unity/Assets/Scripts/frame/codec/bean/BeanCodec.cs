using System;
using System.Collections.Generic;
using Nova.Commons.Util;
using Nova.Logger;

namespace Nova.Codec.bean
{
    
    /// <summary>
    /// 自定义Bean编解码器（自动解析字段并递归编解码）
    /// </summary>
    public class BeanCodec : Codec
    {
        //所有用于编解码的字段（按类定义顺序，递归 包含所有子对象的字段）
        private readonly List<FieldCodecMeta> _fieldMetas;

        private BeanCodec(List<FieldCodecMeta> fieldMetas)
        {
            _fieldMetas = fieldMetas ?? throw new ArgumentNullException(nameof(fieldMetas));
        }

        public static BeanCodec Create(List<FieldCodecMeta> fieldMetas)
        {
            return new BeanCodec(fieldMetas);
        }

        public override object Decode(ByteBuff buff, Type type)
        {
            // 创建Bean实例（要求Bean有默认构造函数，支持可空值类型）
            object bean;
            try
            {
                bean = Activator.CreateInstance(type, nonPublic: true); // 支持私有无参构造函数
            }
            catch (Exception ex)
            {
                throw new InvalidOperationException($"无法创建类型 {type.FullName} 的实例，请确保存在无参构造函数", ex);
            }

            // 递归解码每个字段并赋值（严格按字段顺序，与Java保持一致）
            foreach (FieldCodecMeta meta in _fieldMetas)
            {
                try
                {
                    object fieldValue = meta.Codec.Decode(buff, meta.Field.FieldType);
                    meta.Field.SetValue(bean, fieldValue);
                }
                catch (Exception ex)
                {
                    LoggerUtil.Error($"解码字段 {type.FullName}.{meta.Field.Name} 失败");
                    throw new InvalidOperationException(
                        $"解码字段 {meta.Field.DeclaringType.FullName}.{meta.Field.Name} 失败", ex);
                }
            }

            return bean;
        }

        public override void Encode(ByteBuff buff, object value)
        {
            if (value == null)
            {
                throw new ArgumentNullException(nameof(value), "Bean实例不能为null");
            }

            // 递归编码每个字段（严格按字段顺序，与解码顺序一致，否则会数据错乱）
            foreach (FieldCodecMeta meta in _fieldMetas)
            {
                try
                {
                    // 获取字段值（支持私有字段）
                    object fieldValue = meta.Field.GetValue(value);
                    // 处理可空类型：null时传入默认值
                    if (fieldValue == null && meta.Field.FieldType.IsNullableType())
                    {
                        fieldValue = Activator.CreateInstance(Nullable.GetUnderlyingType(meta.Field.FieldType));
                    }

                    // 调用字段对应的编解码器编码
                    meta.Codec.Encode(buff, fieldValue);
                }
                catch (Exception ex)
                {
                    throw new InvalidOperationException(
                        $"编码字段 {meta.Field.DeclaringType.FullName}.{meta.Field.Name} 失败", ex);
                }
            }
        }
    }

    /// <summary>
    /// 类型扩展方法（辅助判断可空类型）
    /// </summary>
    internal static class TypeExtensions
    {
        /// <summary>
        /// 判断是否为可空类型（Nullable<T> 或引用类型）
        /// </summary>
        public static bool IsNullableType(this Type type)
        {
            // 核心条件：是值类型 + 是泛型 + 泛型定义是 Nullable<>
            return type.IsValueType
                   && type.IsGenericType
                   && type.GetGenericTypeDefinition() == typeof(Nullable<>);
        }
    }
}