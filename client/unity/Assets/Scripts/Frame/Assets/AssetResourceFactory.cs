using System.Collections.Generic;
using System.Linq;
using UnityEngine;

namespace frame.Assets
{
    /// <summary>
    /// 资源工厂
    /// </summary>
    public class AssetResourceFactory : ScriptableObject
    {
        public List<AssetTextItem> textItems;


        /// <summary>
        /// 获取文本资源
        /// </summary>
        /// <param name="group">指定分组</param>
        /// <param name="name">指定资源名称</param>
        /// <returns></returns>
        public TextAsset GetTextAsset(string group, string name)
        {
            return textItems.Find(item => item.group == group).pool.ToList().Find(item => item.name == name);
        }
    }
}