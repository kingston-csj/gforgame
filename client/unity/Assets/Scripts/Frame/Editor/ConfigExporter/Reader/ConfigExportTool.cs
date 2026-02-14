using System;
using System.Collections;
using System.IO;
using System.Reflection;
using System.Text;
using Frame.Config;
using Nova.Commons.Convert;
using Nova.Commons.Util;
using Nova.Data;
using Nova.Net.Socket;
using UnityEditor;
using UnityEngine;

namespace Nova.Editor.ConfigExporter
{
    // 必须引入Editor命名空间

    /// <summary>
    /// Unity顶部菜单：Tools→导出配置
    /// </summary>
    public class ConfigExportTool
    {
        // 核心：创建顶部菜单，路径为Tools/导出配置（优先级1000，确保在Tools下显示）
        [MenuItem("Tools/导出配置", false, 1000)]
        public static void ExportConfig()
        {
            try
            {
                // 1. 初始化解析器 
                IConversionService conversionService = new GenericConversionService();
                IDataReader excelReader = new ExcelReader(conversionService);

                ConfigEditorBinding editorCfg = ConfigEditorBinding.ins;
                string resultMsg = "";

                // 扫描所有的配置文件类
                foreach (Type item in ClassScanner.ListAllSubclasses(typeof(AbsConfigData)))
                {
                    // 寻找对应的excel文件
                    DataTable tableMeta = item.GetCustomAttribute(typeof(DataTable)) as DataTable;
                    // 首字母小写
                    string fileName = char.ToLowerInvariant(item.Name[0]) + item.Name.Substring(1) + editorCfg.Suffix;
                    if (tableMeta != null)
                    {
                        fileName = tableMeta.name + editorCfg.Suffix;
                    }

                    //  Excel文件路径
                    string relativePath = Path.Combine(editorCfg.SourceDirectory, fileName); // 先用Path.Combine拼接相对路径
                    string excelPath = Path.Combine(Application.dataPath, relativePath); // 拼接绝对路径
                    excelPath = Path.GetFullPath(excelPath);
                    // 校验文件是否存在
                    if (!File.Exists(excelPath))
                    {
                        EditorUtility.DisplayDialog("错误", $"配置表不存在：{excelPath}", "确定");
                        return;
                    }

                    // 3. 读取并解析Excel
                    IList configs = null;
                    using (var stream = new FileStream(excelPath, FileMode.Open, FileAccess.Read))
                    {
                        // 获取ExcelReader的Read泛型方法
                        MethodInfo readMethod = typeof(ExcelReader)
                            .GetMethod("Read", new[] { typeof(Stream), typeof(Type) });

                        // 把泛型方法的T替换为item（动态类型）
                        MethodInfo genericMethod = readMethod.MakeGenericMethod(item);

                        // 执行方法，结果转IList（所有List<T>都实现了IList）
                        configs = (IList)genericMethod.Invoke(excelReader, new object[] { stream, item });
                        resultMsg += $"成功导出{fileName}，共{configs.Count}条配置数据\n";
                        // 写入JSON文件（UTF8无BOM格式，避免中文乱码）
                        writeJsonFile(fileName.Replace(editorCfg.Suffix, ".json"), configs);
                    }
                }

                Debug.Log(resultMsg);
                EditorUtility.DisplayDialog("导出成功", resultMsg, "确定");
            }
            catch (Exception ex)
            {
                // 异常弹窗提示（比Console更直观）
                EditorUtility.DisplayDialog("导出失败", $"错误信息：{ex.Message}\n堆栈：{ex.StackTrace}", "确定");
                Debug.LogError($"配置导出失败：{ex.Message}");
            }
        }


        private static void writeJsonFile(string fileName, IList configs)
        {
            ConfigEditorBinding editorCfg = ConfigEditorBinding.ins;

            string jsonFileName = Path.GetFileNameWithoutExtension(fileName) + ".json";
            string relativeJsonPath = Path.Combine(editorCfg.OutputDirectory, jsonFileName);
            string absoluteJsonPath = Path.Combine(Application.dataPath, relativeJsonPath);
            // 标准化路径（统一分隔符、处理../等相对路径）
            absoluteJsonPath = Path.GetFullPath(absoluteJsonPath);

            // 确保输出目录存在（不存在则自动创建）
            string jsonOutputDir = Path.GetDirectoryName(absoluteJsonPath);
            Directory.CreateDirectory(jsonOutputDir);
            File.WriteAllText(absoluteJsonPath, JsonUtil.ToJson(configs), Encoding.UTF8);
        }
    }
}