using System;
using System.Reflection;
using frame.Assets;
using Game.Confi;
using Nova.Data;
using UnityEngine;
using AppContext = Game.Core.AppContext;

namespace Game.Configs
{
    public class DataManager
    {
        public ConfigItemContainer configItemContainer;
        public ConfigCommonContainer configCommonContainer;

        private const string CONFIG_ROOT_PATH = "config/";

        /// <summary>
        /// 自动初始化所有 ConfigContainer<> 子类字段
        /// </summary>
        public void AutoInit()
        {
            Debug.Log("开始自动初始化配置容器...");

            // 反射获取 DataManager 所有字段（包含 public/private，可按需调整）
            FieldInfo[] fields = typeof(DataManager).GetFields(
                BindingFlags.Public | BindingFlags.Instance | BindingFlags.NonPublic);

            foreach (FieldInfo field in fields)
            {
                Type fieldType = field.FieldType;
                Debug.Log($"检查字段：{field.Name}，类型：{fieldType.FullName}");

                // 判断字段是否是 ConfigContainer<> 的子类
                if (IsSubclassOfGeneric(fieldType, typeof(ConfigContainer<>)))
                {
                    try
                    {
                        // 从父类获取泛型参数 E（如 ConfigItemData）
                        Type genericParentType = fieldType.BaseType;
                        Type[] genericArguments = genericParentType.GetGenericArguments();
                        if (genericArguments.Length == 0)
                        {
                            Debug.LogError($"字段 {field.Name} 的父类未指定泛型参数！");
                            continue;
                        }

                        Type configDataType = genericArguments[0]; // 拿到 E = ConfigItemData
                        Debug.Log($"字段 {field.Name} 对应的配置数据类型：{configDataType.Name}");
                        // 加载对应的 TextAsset 配置文件（约定：Configs/泛型参数类名.json）
                        // 名字转为itemData
                        string fileName = configDataType.Name.Replace("Config", "");
                        // 首字母小写
                        fileName = char.ToLowerInvariant(fileName[0]) + fileName.Substring(1);
                        TextAsset textAsset =
                            AppContext.assetResourceFactory.GetTextAsset(AssetResourceGroup.JsonConfig, fileName);
                        if (textAsset == null)
                        {
                            Debug.LogError(
                                $"未找到配置文件：{fileName}（请检查 Resources/config/ 目录下是否存在 {fileName}.json）");
                            continue;
                        }

                        // 找到容器的构造函数（参数为 TextAsset）
                        ConstructorInfo constructor = fieldType.GetConstructor(new Type[] { typeof(TextAsset) });
                        if (constructor == null)
                        {
                            Debug.LogError($"容器类 {fieldType.Name} 缺少「接收 TextAsset 的构造函数」！");
                            continue;
                        }

                        // 实例化容器
                        object containerInstance = constructor.Invoke(new object[] { textAsset });
                        if (containerInstance == null)
                        {
                            Debug.LogError($"实例化字段 {field.Name} 失败：构造函数返回 null！");
                            continue;
                        }

                        //  给 DataManager 字段赋值
                        field.SetValue(this, containerInstance);


                        Debug.Log($"实例化容器：{field.Name}（配置文件：{fileName}）");
                    }
                    catch (Exception e)
                    {
                        Debug.LogError($"实例化字段 {field.Name} 失败：{e.Message}\n{e.StackTrace}");
                    }
                }
            }
        }

        /// <summary>
        /// 辅助方法：判断类型是否是某个泛型类型的子类（如 ConfigItemContainer → ConfigContainer<>）
        /// </summary>
        private bool IsSubclassOfGeneric(Type type, Type genericBaseType)
        {
            if (type == null || genericBaseType == null || !genericBaseType.IsGenericTypeDefinition)
                return false;

            Type currentType = type;
            while (currentType != typeof(object) && currentType != null)
            {
                // 对比泛型定义（如 ConfigContainer<ConfigItemData> → ConfigContainer<>）
                if (currentType.IsGenericType && currentType.GetGenericTypeDefinition() == genericBaseType)
                    return true;

                currentType = currentType.BaseType;
            }

            return false;
        }
    }
}