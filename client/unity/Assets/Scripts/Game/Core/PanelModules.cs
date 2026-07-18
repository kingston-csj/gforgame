using System.Collections.Generic;
using Nova.Ui;

namespace Game.Core
{
    public class PanelModules
    {
        // 在这里，注册所有的面板配置
        public static PanelPrefabConfig Login =
            new PanelPrefabConfig(R.Panels.Login, "gui/LoginPanel", LayerIds.layer5);

        public static List<PanelPrefabConfig> GetAllPanelConfigs()
        {
            List<PanelPrefabConfig> configs = new List<PanelPrefabConfig>();
            // 遍历所有属性，自动添加
            foreach (var item in typeof(PanelModules).GetFields())
            {
                if (item.FieldType == typeof(PanelPrefabConfig))
                {
                    PanelPrefabConfig config = item.GetValue(null) as PanelPrefabConfig;
                    configs.Add(config);
                }
            }

            return configs;
        }
    }
}