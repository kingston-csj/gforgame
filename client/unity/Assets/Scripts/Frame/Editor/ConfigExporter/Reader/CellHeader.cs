using System.Reflection;

namespace Nova.Editor.ConfigExporter
{
    /// <summary>
    /// 单元格表头信息
    /// </summary>
    public class CellHeader
    {
        /// <summary>
            /// 列名（对应配置表header行的字段名）
            /// </summary>
            public string Column { get; set; }
        
            /// <summary>
            /// 对应实体类的字段信息
            /// </summary>
            public FieldInfo Field { get; set; }
    }
}