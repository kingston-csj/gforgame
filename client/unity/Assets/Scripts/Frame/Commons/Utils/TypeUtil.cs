using System;

namespace Frame.Commons.Utils
{
    /// <summary>
    /// 类型工具类，提供基本类型/字符串判断、类型兼容性校验
    /// </summary>
    public static class TypeUtil
    {
        /// <summary>
        /// 判断类型是否为基本类型（值类型）、基础类型的可空类型或字符串
        /// </summary>
        /// <param name="type">要判断的类型</param>
        /// <returns>是基本类型/包装类/字符串返回true，否则返回false</returns>
        public static bool IsPrimitiveOrString(Type type)
        {
            if (type == null)
            {
                return false;
            }

            // 判断是否为C#基本类型
            if (type.IsPrimitive)
            {
                return true;
            }

            // 判断是否为基础值类型的可空类型
            if (IsWrapperType(type))
            {
                return true;
            }

            // 判断是否为字符串类型
            if (type == typeof(string))
            {
                return true;
            }

            return false;
        }

        /// <summary>
        /// 判断是否为C#基础值类型的可空类型
        /// </summary>
        /// <param name="type">要判断的类型</param>
        /// <returns>是基础值类型的可空类型返回true，否则返回false</returns>
        private static bool IsWrapperType(Type type)
        {
            // C#中无专门"包装类"
            if (Nullable.GetUnderlyingType(type) is Type underlyingType)
            {
                return underlyingType == typeof(int)
                       || underlyingType == typeof(long)
                       || underlyingType == typeof(short)
                       || underlyingType == typeof(byte)
                       || underlyingType == typeof(float)
                       || underlyingType == typeof(double)
                       || underlyingType == typeof(bool)
                       || underlyingType == typeof(char);
            }

            // 兼容直接传入基础值类型（如typeof(int)
            return type == typeof(int)
                   || type == typeof(long)
                   || type == typeof(short)
                   || type == typeof(byte)
                   || type == typeof(float)
                   || type == typeof(double)
                   || type == typeof(bool)
                   || type == typeof(char);
        }

        /// <summary>
        /// 判断值是否与类型兼容
        /// 示例：
        /// 1. 类型为 int，值为 1 时，返回 true
        /// 2. 类型为 int，值为 1L 时，返回 true
        /// 3. 类型为 int，值为 "1" 时，返回 true
        /// 4. 类型为 int，值为 "1.0" 时，返回 false
        /// </summary>
        /// <param name="value">待校验的值</param>
        /// <param name="type">目标类型</param>
        /// <returns>值与类型兼容返回true，否则返回false</returns>
        public static bool IsCompatibleType(object value, Type type)
        {
            // 空值兼容所有类型；值是目标类型的实例直接兼容
            if (value == null || type.IsInstanceOfType(value))
            {
                return true;
            }

            // 处理值类型
            Type valueType = value.GetType();
            if (type == typeof(int) && (valueType == typeof(int) || valueType == typeof(long)))
            {
                return true;
            }

            if (type == typeof(long) && valueType == typeof(long))
            {
                return true;
            }

            if (type == typeof(double) && valueType == typeof(double))
            {
                return true;
            }

            if (type == typeof(float) && (valueType == typeof(float) || valueType == typeof(double)))
            {
                return true;
            }

            if (type == typeof(short) &&
                (valueType == typeof(short) || valueType == typeof(int) || valueType == typeof(long)))
            {
                return true;
            }

            if (type == typeof(byte) && (valueType == typeof(byte) || valueType == typeof(short) ||
                                         valueType == typeof(int) || valueType == typeof(long)))
            {
                return true;
            }

            if (type == typeof(char) && valueType == typeof(char))
            {
                return true;
            }

            if (type == typeof(bool) && valueType == typeof(bool))
            {
                return true;
            }

            return false;
        }
    }
}