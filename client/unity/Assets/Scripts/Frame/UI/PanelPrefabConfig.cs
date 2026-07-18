namespace Nova.Ui
{
    /// <summary>
    /// 预制体资源项
    /// </summary>
    public class PanelPrefabConfig
    {
        /// <summary>
        /// 资源名称
        /// </summary>
        public string PanelName;

        /// <summary>
        /// 资源路径
        /// </summary>
        public string PrefabPath;


        public LayerIds Layer;
        
        public PanelPrefabConfig(string panelName, string prefabPath, LayerIds layer)
        {
            PanelName = panelName;
            PrefabPath = prefabPath;
            Layer = layer;
        }
    }
}