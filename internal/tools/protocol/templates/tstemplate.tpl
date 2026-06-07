/**
 * {{.ClassComment}}
 */
export default class {{.ClassName}} {
    {{if .Cmd}} public static cmd: number =  {{.Cmd}}; {{end}}
        {{range .Fields}}
        /** {{.Comment}} */
        public  {{.Name}} : {{.FieldType}};
        {{end}}
    
}