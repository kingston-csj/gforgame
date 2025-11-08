using Game.Configs;
using Nova.Data;
using UnityEngine;

namespace Game.Confi
{
    public class ConfigItemContainer:ConfigContainer<ConfigItemData>
    {
        
        public ConfigItemContainer(TextAsset textAsset):base(textAsset)
        {
            
        }
    }
}