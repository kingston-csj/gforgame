using System;
using Nova.Commons.Convert;

namespace Nova.Editor.ConfigExporter
{
    using UnityEngine;
    using UnityEditor; // 必须引入Editor命名空间
    using System.IO;
    using System.Collections.Generic;

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

                // 2. Excel文件路径（Assets/Config下的RoleConfig.xlsx）
                string excelPath = Path.Combine(Application.dataPath, "Config/common.xlsx");

                // 校验文件是否存在
                if (!File.Exists(excelPath))
                {
                    EditorUtility.DisplayDialog("错误", $"配置表不存在：{excelPath}", "确定");
                    return;
                }

                // 3. 读取并解析Excel
                List<CommonConfig> configs;
                using (var stream = new FileStream(excelPath, FileMode.Open, FileAccess.Read))
                {
                    configs = excelReader.Read<CommonConfig>(stream, typeof(CommonConfig));
                }

                // 4. 提示导出结果（用Editor弹窗更友好）
                string resultMsg = $"成功导出{configs.Count}条配置数据！\n";
                foreach (var cfg in configs)
                {
                    resultMsg += $"ID：{cfg.id}，名称：{cfg.value}，等级：{cfg.key}\n";
                }

                EditorUtility.DisplayDialog("导出成功", resultMsg, "确定");
                Debug.Log(resultMsg);
            }
            catch (System.Exception ex)
            {
                // 异常弹窗提示（比Console更直观）
                EditorUtility.DisplayDialog("导出失败", $"错误信息：{ex.Message}\n堆栈：{ex.StackTrace}", "确定");
                Debug.LogError($"配置导出失败：{ex.Message}");
            }
        }
    }

    public class CommonConfig
    {
        public int id;
        public string value;
        public string desc;
        public string key;
    }
}