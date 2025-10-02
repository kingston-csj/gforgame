using System;
using System.Reflection;

namespace Nova.Commons.Util
{
    /// <summary>
    /// 反射工具类
    /// 提供静态方法调用、属性访问等反射操作
    /// </summary>
    public class ReflectUtil
    {
        /// <summary>
        /// 调用静态方法
        /// </summary>
        /// <param name="type"></param>
        /// <param name="methodName"></param>
        /// <param name="args"></param>
        /// <returns></returns>
        public static object CallStaticMethod(Type type, string methodName, params object[] args)
        {
            MethodInfo method = type.GetMethod(methodName, BindingFlags.Public | BindingFlags.Static);

            return method.Invoke(null, args);
        }
    }
    
}