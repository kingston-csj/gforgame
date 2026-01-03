using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using UnityEngine;

namespace Nova.Commons.Util
{
    /// <summary>
    /// 类扫描器工具类
    /// 支持扫描目标类，包括查找子类、查找带特性的类
    /// </summary>
    public static class ClassScanner
    {
        /// <summary>
        /// 默认过滤器（匹配所有类）
        /// </summary>
        private static readonly Func<Type, bool> EmptyFilter = _ => true;


        /// <summary>
        /// 查找指定类的所有子类（不包括抽象类）
        /// </summary>
        /// <param name="parentType">父类/接口类型</param>
        /// <returns>所有非抽象子类类型集合</returns>
        public static HashSet<Type> ListAllSubclasses(Type parentType)
        {
            return ListClasses(type =>
                parentType.IsAssignableFrom(type) && !type.IsAbstract && type.IsClass);
        }

        /// <summary>
        /// 查找指定包下所有带指定特性的类
        /// </summary>
        /// <typeparam name="TAnnotation">特性类型</typeparam>
        /// <returns>所有带指定特性的类类型集合</returns>
        public static HashSet<Type> ListClassesWithAttribution<TAnnotation>()
            where TAnnotation : Attribute
        {
            return ListClasses(type =>
                type.GetCustomAttribute<TAnnotation>(inherit: false) != null);
        }

        /// <summary>
        /// 扫描所有程序集中的类（支持自定义过滤规则）
        /// </summary>
        /// <param name="filter">类过滤规则</param>
        /// <returns>符合规则的类类型集合</returns>
        public static HashSet<Type> ListClasses(Func<Type, bool> filter)
        {
            HashSet<Type> result = new HashSet<Type>();
            try
            {
                // 从运行时程序集中查找对应的 Type（仅匹配 Assets 下的脚本类）
                Assembly[] assemblies = AppDomain.CurrentDomain.GetAssemblies()
                    .Where(asm => !IsSystemOrUnityAssembly(asm))
                    .ToArray();

                foreach (Assembly assembly in assemblies)
                {
                    Type[] types;
                    try
                    {
                        types = assembly.GetTypes();
                    }
                    catch (ReflectionTypeLoadException ex)
                    {
                        types = ex.Types.Where(t => t != null).ToArray();
                        Debug.LogWarning($"ClassScanner: 程序集 {assembly.GetName().Name} 部分类型加载失败");
                    }

                    // 匹配：完整类名在目标集合中 + 符合过滤规则
                    foreach (Type type in types)
                    {
                        if (!string.IsNullOrEmpty(type.FullName) && filter(type))
                        {
                            result.Add(type);
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                Debug.LogError($"ClassScanner: 扫描失败 - 错误: {ex.Message}");
            }

            Debug.Log($"ClassScanner: 扫描完成 - 找到类数: {result.Count}");
            return result;
        }

        /// <summary>
        /// 判断是否为系统/Unity 内置程序集（跳过）
        /// </summary>
        private static bool IsSystemOrUnityAssembly(Assembly assembly)
        {
            string name = assembly.GetName().Name;
            return name.StartsWith("System.") || name.StartsWith("Unity.") || name.StartsWith("UnityEngine.") ||
                   name.StartsWith("UnityEditor.");
        }
    }
}