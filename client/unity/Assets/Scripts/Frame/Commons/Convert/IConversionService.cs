using System;

namespace Nova.Commons.Convert
{
    /// <summary>
    /// 类型转换服务接口，借鉴java spring框架
    /// </summary>
    public interface IConversionService
    {

        /// <summary>
        /// 将字符串转换为目标类型
        /// </summary>
        /// <param name="sourceValue">原始字符串值</param>
        /// <param name="targetType">目标类型</param>
        /// <returns>转换后的值</returns>
        object Convert(string sourceValue, Type targetType);
    }
}