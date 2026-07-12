namespace Frame.Commons.Utils
{
    public static class StringExtensions
    {
        /// <summary>
        /// 判断字符串是否为空
        /// </summary>
        /// <param name="str"></param>
        /// <returns></returns>
        public static bool IsEmpty(this string str)
        {
            return string.IsNullOrEmpty(str);
        }


        ///
        /// 将字符串转换为 Base64 编码
        /// 
        public static string ToBase64(this string str)
        {
            if (str.IsEmpty()) return string.Empty;
            return System.Convert.ToBase64String(System.Text.Encoding.UTF8.GetBytes(str));
        }

        ///
        /// 将 Base64 编码转换为字符串
        /// 
        public static string FromBase64(this string base64Str)
        {
            if (base64Str.IsEmpty()) return string.Empty;
            return System.Text.Encoding.UTF8.GetString(System.Convert.FromBase64String(base64Str));
        }
    }
}