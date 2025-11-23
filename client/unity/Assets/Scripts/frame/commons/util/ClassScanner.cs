using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text.RegularExpressions;
using UnityEngine;

namespace Nova.Commons.Util
{
    /// <summary>
    /// 类扫描器工具类（仅扫描 Unity Assets 目录下的 C# 脚本，返回 HashSet<Type>）
    /// 支持扫描指定包下的类、查找子类、查找带注解的类
    /// </summary>
    public static class ClassScanner
    {
        /// <summary>
        /// 默认过滤器（匹配所有类）
        /// </summary>
        private static readonly Func<Type, bool> EmptyFilter = _ => true;

        // 正则表达式（解析脚本中的类名和命名空间）
        private static readonly Regex NamespaceRegex =
            new Regex(@"namespace\s+([\w\.]+)\s*\{", RegexOptions.Compiled | RegexOptions.Multiline);

        private static readonly Regex ClassRegex =
            new Regex(@"(public|private|protected|internal)?\s*(static)?\s*(abstract)?\s*class\s+([\w<>]+)",
                RegexOptions.Compiled | RegexOptions.Singleline);

        #region 核心接口（完全对齐 Java）

        /// <summary>
        /// 扫描指定包下的所有类（无过滤）
        /// </summary>
        /// <param name="scanPackage">搜索的包根路径（如 "mgame.module"）</param>
        /// <returns>所有匹配的类类型集合</returns>
        public static HashSet<Type> ListClasses(string scanPackage)
        {
            return ListClasses(scanPackage, EmptyFilter);
        }

        /// <summary>
        /// 查找指定包下的所有子类（不包括抽象类）
        /// </summary>
        /// <param name="scanPackage">搜索的包根路径</param>
        /// <param name="parentType">父类/接口类型</param>
        /// <returns>所有非抽象子类类型集合</returns>
        public static HashSet<Type> ListAllSubclasses(string scanPackage, Type parentType)
        {
            return ListClasses(scanPackage, type =>
                parentType.IsAssignableFrom(type) && !type.IsAbstract && type.IsClass);
        }

        /// <summary>
        /// 查找指定包下所有带指定注解的类
        /// </summary>
        /// <typeparam name="TAnnotation">注解类型</typeparam>
        /// <param name="scanPackage">搜索的包根路径</param>
        /// <returns>所有带指定注解的类类型集合</returns>
        public static HashSet<Type> ListClassesWithAnnotation<TAnnotation>(string scanPackage)
            where TAnnotation : Attribute
        {
            return ListClasses(scanPackage, type =>
                type.GetCustomAttribute<TAnnotation>(inherit: false) != null);
        }

        /// <summary>
        /// 扫描指定包下的类（支持自定义过滤规则）
        /// </summary>
        /// <param name="scanPackage">搜索的包根路径</param>
        /// <param name="filter">类过滤规则</param>
        /// <returns>符合规则的类类型集合</returns>
        public static HashSet<Type> ListClasses(string scanPackage, Func<Type, bool> filter)
        {
            HashSet<Type> result = new HashSet<Type>();
            if (string.IsNullOrEmpty(scanPackage))
            {
                Debug.LogError("ClassScanner: 扫描包路径不能为空");
                return result;
            }

            try
            {
                // 1. 扫描 Assets 下目标包对应的目录，收集所有完整类名
                HashSet<string> targetClassNames = CollectTargetClassNames(scanPackage);
                if (targetClassNames.Count == 0)
                {
                    Debug.LogWarning($"ClassScanner: 未找到任何 .cs 脚本 - 包路径: {scanPackage}");
                    return result;
                }

                // 2. 从运行时程序集中查找对应的 Type（仅匹配 Assets 下的脚本类）
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
                        if (!string.IsNullOrEmpty(type.FullName) &&
                            targetClassNames.Contains(type.FullName) &&
                            filter(type))
                        {
                            result.Add(type);
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                Debug.LogError($"ClassScanner: 扫描失败 - 包路径: {scanPackage}, 错误: {ex.Message}");
            }

            Debug.Log($"ClassScanner: 扫描完成 - 包路径: {scanPackage}, 找到类数: {result.Count}");
            return result;
        }

        #endregion

        #region 辅助方法

        /// <summary>
        /// 扫描 Assets 目录下的脚本，收集目标包的完整类名
        /// </summary>
        private static HashSet<string> CollectTargetClassNames(string scanPackage)
        {
            HashSet<string> classNames = new HashSet<string>();
            string rootDir = Path.Combine(Application.dataPath, scanPackage.Replace('.', Path.DirectorySeparatorChar));
            if (!Directory.Exists(rootDir)) return classNames;

            HashSet<string> processedFiles = new HashSet<string>();
            ScanDirectory(rootDir, processedFiles, content =>
            {
                string @namespace = ExtractNamespace(content);
                MatchCollection classMatches = ClassRegex.Matches(content);
                foreach (Match match in classMatches)
                {
                    string className = match.Groups[4].Value;
                    if (string.IsNullOrEmpty(@namespace))
                    {
                        classNames.Add(className);
                    }
                    else
                    {
                        classNames.Add($"{@namespace}.{className}");
                    }
                }
            });

            return classNames;
        }

        /// <summary>
        /// 递归扫描目录下的 .cs 文件
        /// </summary>
        private static void ScanDirectory(string dirPath, HashSet<string> processedFiles,
            Action<string> fileContentHandler)
        {
            try
            {
                // 扫描 .cs 文件
                foreach (string file in Directory.GetFiles(dirPath, "*.cs"))
                {
                    if (processedFiles.Contains(file)) continue;
                    processedFiles.Add(file);

                    try
                    {
                        string content = File.ReadAllText(file);
                        fileContentHandler(content);
                    }
                    catch
                    {
                        /* 忽略读取失败的文件 */
                    }
                }

                // 递归扫描子目录（跳过无用目录）
                foreach (string subDir in Directory.GetDirectories(dirPath))
                {
                    string dirName = Path.GetFileName(subDir);
                    if (dirName != "bin" && dirName != "obj" && dirName != "Library" && dirName != "Packages")
                    {
                        ScanDirectory(subDir, processedFiles, fileContentHandler);
                    }
                }
            }
            catch
            {
                /* 忽略目录访问异常 */
            }
        }

        /// <summary>
        /// 提取脚本的命名空间
        /// </summary>
        private static string ExtractNamespace(string content)
        {
            Match match = NamespaceRegex.Match(content);
            return match.Success ? match.Groups[1].Value : string.Empty;
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

        #endregion
    }
}