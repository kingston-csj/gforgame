using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using UnityEngine;
using UnityEngine.Networking;

namespace Frame.Net.Http
{
    /// <summary>
    /// http请求工具，提供GET、POST请求
    /// </summary>
    public class HttpClient
    {
        /// <summary>
        /// GET 请求
        /// </summary>
        public static string Get(string url, Dictionary<string, object> paramsMap)
        {
            return Get(url, null, paramsMap);
        }

        /// <summary>
        /// GET 请求（带 Header）
        /// </summary>
        public static string Get(string url, Dictionary<string, string> headers, Dictionary<string, object> paramsMap)
        {
            if (paramsMap != null && paramsMap.Count > 0)
            {
                url += BuildQueryString(paramsMap);
            }

            using (var request = UnityWebRequest.Get(url))
            {
                if (headers != null)
                {
                    foreach (var header in headers)
                    {
                        request.SetRequestHeader(header.Key, header.Value);
                    }
                }

                // 同步发送（Unity 可用）
                var operation = request.SendWebRequest();
                while (!operation.isDone)
                {
                }

                if (request.result != UnityWebRequest.Result.Success)
                    throw new IOException(request.error);

                return request.downloadHandler.text;
            }
        }

        /// <summary>
        /// POST 请求（自动 JSON 解析返回值）
        /// </summary>
        public static T Post<T>(string url, Dictionary<string, object> paramsMap, Type responseClazz)
        {
            return Post<T>(url, null, paramsMap, responseClazz);
        }

        /// <summary>
        /// POST 请求（带 Header + JSON 解析返回）
        /// </summary>
        public static T Post<T>(string url, Dictionary<string, string> headers, Dictionary<string, object> paramsMap,
            Type responseClazz)
        {
            string json = JsonUtility.ToJson(paramsMap);

            using (var request = new UnityWebRequest(url, "POST"))
            {
                byte[] body = Encoding.UTF8.GetBytes(json);
                request.uploadHandler = new UploadHandlerRaw(body);
                request.downloadHandler = new DownloadHandlerBuffer();

                // 默认 JSON 请求头
                request.SetRequestHeader("Content-Type", "application/json");

                if (headers != null)
                {
                    foreach (var header in headers)
                    {
                        request.SetRequestHeader(header.Key, header.Value);
                    }
                }

                var operation = request.SendWebRequest();
                while (!operation.isDone)
                {
                }

                if (request.result != UnityWebRequest.Result.Success)
                    throw new IOException(request.error);

                string response = request.downloadHandler.text;
                return (T)JsonUtility.FromJson(response, responseClazz);
            }
        }

        #region 构建 GET 参数

        private static string BuildQueryString(Dictionary<string, object> parameters)
        {
            StringBuilder sb = new StringBuilder();
            sb.Append("?");

            foreach (var p in parameters)
            {
                sb.Append(p.Key);
                sb.Append("=");
                sb.Append(UnityWebRequest.EscapeURL(p.Value?.ToString() ?? ""));
                sb.Append("&");
            }

            if (sb.Length > 1)
                sb.Length--;

            return sb.ToString();
        }

        #endregion
    }
}