namespace Nova.Editor.ConfigExporter
{
    public class CellColumn
    {
        /// <summary>
        /// 列头信息
        /// </summary>
        public CellHeader Header { get; set; }

        /// <summary>
        /// 单元格原始字符串值
        /// </summary>
        public string Value { get; set; }
    }
}