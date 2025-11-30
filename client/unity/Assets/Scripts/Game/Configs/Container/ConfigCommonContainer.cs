using System.Collections.Generic;
using Game.Configs;
using Nova.Data;
using UnityEngine;

namespace Game.Confi
{
    public class ConfigCommonContainer : ConfigContainer<ConfigCommonData>
    {
        // common.json 的映射表
        private Dictionary<string, string> mapper = new();

        public ConfigCommonContainer(TextAsset textAsset) : base(textAsset)
        {
        }

        protected override void AfterLoad()
        {
            foreach (var item in GetItems())
            {
                mapper.Add(item.key, item.value);
            }
        }

        public string GetValue(string key)
        {
            return mapper[key];
        }
    }
}