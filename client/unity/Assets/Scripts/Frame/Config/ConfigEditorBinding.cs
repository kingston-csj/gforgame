using UnityEditor;
using UnityEngine;

namespace Frame.Config
{
    [CreateAssetMenu(
        fileName = "AssetConfigBinding", // 新建 .asset 文件的默认名称
        menuName = "Game/资源工厂", // 右键菜单路径（Assets → Game → 资源工厂）
        order = 10 // 菜单排序（越小越靠上，避免和其他菜单冲突）
    )]
    public class ConfigEditorBinding :ScriptableObject
    {

        private static ConfigEditorBinding _cfg;
        
        /// <summary>
        /// 文件后缀，根据文件类型选择对应的文件解析工具
        /// </summary>
        public string Suffix = ".xlsx";
        
        /// <summary>
        /// 配置文件的路径
        /// </summary>
         public string SourceDirectory = "Config/";
        
        /// <summary>
        /// json文件的输出路径
        /// </summary>
        public string OutputDirectory = "Output/";
        
        public static ConfigEditorBinding ins
        {
            get
            {
                if (_cfg == null)
                {
                    _cfg = AssetDatabase.LoadAssetAtPath<ConfigEditorBinding>("Assets/Editor/AssetConfigBinding.asset");
                }
                return _cfg;
            }

        }
    }
}