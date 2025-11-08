namespace Nova.Commons.Util
{
    using System;
    using Newtonsoft.Json;
    using UnityEngine;

    /// <summary>
    /// 模仿 Unity JsonUtility API 的 Newtonsoft.Json 工具类
    /// 优势：支持根数组、大小写不敏感、忽略多余字段、兼容复杂类型
    /// </summary>
    public static class JsonUtil
    {
        /// <summary>
        /// 默认序列化配置（解决 Unity 开发常见痛点）
        /// </summary>
        private static readonly JsonSerializerSettings DefaultSettings = new JsonSerializerSettings
        {
            // 忽略 JSON 中不存在的字段（配置迭代时不会报错）
            MissingMemberHandling = MissingMemberHandling.Ignore,
            // 忽略空值字段（序列化时不输出 null 值）
            NullValueHandling = NullValueHandling.Ignore,
            // 支持循环引用（避免复杂对象序列化报错）
            ReferenceLoopHandling = ReferenceLoopHandling.Ignore
        };


        /// <summary>
        /// 解析 JSON 字符串到指定类型（模仿 JsonUtility.FromJson<T>）
        /// 支持：普通对象、根数组（直接用 List<T> 接收）、复杂嵌套类型
        /// </summary>
        /// <typeparam name="T">目标类型（可是类、结构体、List<T>、数组）</typeparam>
        /// <param name="json">JSON 字符串</param>
        /// <returns>解析后的对象</returns>
        public static T FromJson<T>(string json)
        {
            try
            {
                if (string.IsNullOrEmpty(json))
                {
                    Debug.LogWarning("JsonUtil.FromJson：JSON 字符串为空！");
                    return default;
                }

                // 用默认配置解析（解决常见痛点）
                return JsonConvert.DeserializeObject<T>(json, DefaultSettings);
            }
            catch (Exception e)
            {
                Debug.LogError($"JsonUtil.FromJson 解析失败（类型：{typeof(T).Name}）：{e.Message}\n{e.StackTrace}");
                return default;
            }
        }

        /// <summary>
        /// 重载：解析 JSON 字符串到指定类型（非泛型版本，模仿 JsonUtility.FromJson）
        /// </summary>
        /// <param name="json">JSON 字符串</param>
        /// <param name="type">目标类型</param>
        /// <returns>解析后的对象</returns>
        public static object FromJson(string json, Type type)
        {
            try
            {
                if (string.IsNullOrEmpty(json) || type == null)
                {
                    Debug.LogWarning("JsonUtil.FromJson：JSON 为空或类型为 null！");
                    return null;
                }

                return JsonConvert.DeserializeObject(json, type, DefaultSettings);
            }
            catch (Exception e)
            {
                Debug.LogError($"JsonUtil.FromJson 解析失败（类型：{type.Name}）：{e.Message}\n{e.StackTrace}");
                return null;
            }
        }
    }
    
}