namespace Game.Util
{
    using System;

    /// <summary>
    /// 数值类型转换工具类
    /// 功能：安全地将任意Object类型转换为指定数值类型，支持默认值和异常兜底
    /// </summary>
    public sealed class NumberUtil
    {
        // 私有构造函数，禁止实例化（工具类）
        private NumberUtil()
        {
        }

        #region 布尔类型转换

        public static bool BooleanValue(object obj)
        {
            return BooleanValue(obj, false);
        }

        public static bool BooleanValue(object obj, bool defaultValue)
        {
            if (obj == null)
            {
                return defaultValue;
            }

            // 直接匹配bool类型
            if (obj is bool)
            {
                return (bool)obj;
            }

            try
            {
                // C#的bool.Parse区分，因此转小写兼容
                return bool.Parse(obj.ToString().ToLower());
            }
            catch (Exception)
            {
                return defaultValue;
            }
        }

        #endregion

        #region 字节类型转换

        public static byte ByteValue(object obj)
        {
            return ByteValue(obj, (byte)0);
        }

        public static byte ByteValue(object obj, byte defaultValue)
        {
            if (obj == null)
            {
                return defaultValue;
            }

            if (obj is byte)
            {
                return (byte)obj;
            }

            try
            {
                return byte.Parse(obj.ToString());
            }
            catch (Exception)
            {
                return defaultValue;
            }
        }

        #endregion

        #region 短整型转换

        public static short ShortValue(object obj)
        {
            return ShortValue(obj, (short)0);
        }

        public static short ShortValue(object obj, short defaultValue)
        {
            if (obj == null)
            {
                return defaultValue;
            }

            if (obj is short)
            {
                return (short)obj;
            }

            try
            {
                return short.Parse(obj.ToString());
            }
            catch (Exception)
            {
                return defaultValue;
            }
        }

        #endregion

        #region 整型转换

        public static int IntValue(object obj)
        {
            return IntValue(obj, 0);
        }

        public static int IntValue(object obj, int defaultValue)
        {
            if (obj == null)
            {
                return defaultValue;
            }

            if (obj is int)
            {
                return (int)obj;
            }

            try
            {
                return int.Parse(obj.ToString());
            }
            catch (Exception)
            {
                return defaultValue;
            }
        }

        #endregion

        #region 长整型转换

        public static long LongValue(object obj)
        {
            return LongValue(obj, 0L);
        }

        public static long LongValue(object obj, long defaultValue)
        {
            if (obj == null)
            {
                return defaultValue;
            }

            if (obj is long)
            {
                return (long)obj;
            }

            try
            {
                return long.Parse(obj.ToString());
            }
            catch (Exception)
            {
                return defaultValue;
            }
        }

        #endregion

        #region 双精度浮点型转换

        public static double DoubleValue(object obj)
        {
            return DoubleValue(obj, 0.0);
        }

        public static double DoubleValue(object obj, double defaultValue)
        {
            if (obj == null)
            {
                return defaultValue;
            }

            if (obj is double)
            {
                return (double)obj;
            }

            try
            {
                return double.Parse(obj.ToString());
            }
            catch (Exception)
            {
                return defaultValue;
            }
        }

        #endregion
    }
}