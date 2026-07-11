namespace Nova.Commons.Util
{
    using System;
    using System.IO;
    using System.IO.Compression;
    using System.Text;

    /// <summary>
    /// 字符串压缩解压工具
    /// </summary>
    public static class CompressUtility
    {
        /// <summary>
        /// 压缩字符串 → 直接返回可存入JSON的Base64字符串
        /// </summary>
        public static string Compress(string text)
        {
            if (string.IsNullOrEmpty(text))
                return string.Empty;

            byte[] data = Encoding.UTF8.GetBytes(text);
            using (var ms = new MemoryStream())
            {
                using (var gzip = new GZipStream(ms, CompressionMode.Compress))
                {
                    gzip.Write(data, 0, data.Length);
                }

                // 压缩后直接转Base64
                return Convert.ToBase64String(ms.ToArray());
            }
        }

        /// <summary>
        /// 从Base64解压回原始字符串
        /// </summary>
        public static string Decompress(string base64Compressed)
        {
            if (string.IsNullOrEmpty(base64Compressed))
                return string.Empty;

            byte[] data = Convert.FromBase64String(base64Compressed);
            using (var ms = new MemoryStream(data))
            using (var gzip = new GZipStream(ms, CompressionMode.Decompress))
            using (var reader = new StreamReader(gzip, Encoding.UTF8))
            {
                return reader.ReadToEnd();
            }
        }
    }
}