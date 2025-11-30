// 自动生成的 C# 代码，请勿手动修改
namespace Game.Net.Message
{
    using Nova.Net.Socket;

    /// <summary>
    /// {{.StructComment}}
    /// </summary>
    {{if .Cmd}} [MessageMeta(Cmd = {{.Cmd}})] {{end}}
    public class {{.StructName}} {{if .Cmd}}: Message {{end}}
    {
        {{range .Fields}}
        /// <summary>
        /// {{.Comment}}
        /// </summary>
        public {{.FieldType}} {{.Name}};
        {{end}}
    }
}