package protocoldoc

import (
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/forfun/gforgame/internal/protos"
	exporter "github.com/forfun/gforgame/internal/tools/protocol"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

type Type struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Ref     string `json:"ref,omitempty"`
	Element *Type  `json:"element,omitempty"`
	Key     *Type  `json:"key,omitempty"`
}

type Field struct {
	Name         string `json:"name"`
	Type         Type   `json:"type"`
	Desc         string `json:"desc"`
	Nullable     bool   `json:"nullable"`
	HasDefault   bool   `json:"hasDefault"`
	DefaultValue any    `json:"defaultValue"`
	JSON         bool   `json:"json"`
}

type Schema struct {
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
	Fields []Field `json:"fields"`
}

type Message struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Module    string `json:"module"`
	Cmd       *int32 `json:"cmd"`
	Direction string `json:"direction"`
	Desc      string `json:"desc"`
	Schema    string `json:"schema"`
	Testable  bool   `json:"testable"`
	Response  string `json:"response,omitempty"`
}

type Module struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Protocols []Message `json:"protocols"`
}

type Connection struct {
	URL         string   `json:"url"`
	Port        string   `json:"port"`
	Path        string   `json:"path"`
	ServerID    uint32   `json:"serverId"`
	TestEnabled bool     `json:"testEnabled"`
	Codecs      []string `json:"codecs"`
}

type Metadata struct {
	Version     int               `json:"version"`
	GeneratedAt string            `json:"generatedAt"`
	Connection  Connection        `json:"connection"`
	Modules     []Module          `json:"modules"`
	Schemas     map[string]Schema `json:"schemas"`
}

func sourceDocs() map[string]exporter.StructInfo {
	generator := &exporter.BaseGenerator{}
	infos, err := generator.ParseSources(protos.DocumentationSources)
	if err != nil {
		panic(err)
	}
	docs := make(map[string]exporter.StructInfo, len(infos))
	for _, info := range infos {
		docs[info.Name] = info
	}
	return docs
}

// Build derives IDs from the live registry and request/response pairs from installed handlers.
func Build(connection Connection, routes *network.MessageRoute) Metadata {
	m := Metadata{Version: 1, GeneratedAt: time.Now().UTC().Format(time.RFC3339), Connection: connection, Modules: []Module{}, Schemas: map[string]Schema{}}
	docs := sourceDocs()
	types := map[string]reflect.Type{}
	var describe func(reflect.Type) Type
	describe = func(t reflect.Type) Type {
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		switch t.Kind() {
		case reflect.String:
			return Type{Kind: "string", Name: "string"}
		case reflect.Bool:
			return Type{Kind: "boolean", Name: "boolean"}
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			return Type{Kind: "long", Name: "long"}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			return Type{Kind: "integer", Name: "int"}
		case reflect.Float32:
			return Type{Kind: "number", Name: "float"}
		case reflect.Float64:
			return Type{Kind: "number", Name: "double"}
		case reflect.Array, reflect.Slice:
			e := describe(t.Elem())
			return Type{Kind: "list", Name: "List<" + e.Name + ">", Element: &e}
		case reflect.Map:
			k, e := describe(t.Key()), describe(t.Elem())
			return Type{Kind: "map", Name: "Map<" + k.Name + ", " + e.Name + ">", Key: &k, Element: &e}
		case reflect.Struct:
			id := t.PkgPath() + "." + t.Name()
			if _, exists := m.Schemas[id]; !exists {
				types[id] = t
				m.Schemas[id] = Schema{Name: t.Name(), Fields: []Field{}}
				schema := Schema{Name: t.Name(), Desc: docs[t.Name()].Comment, Fields: []Field{}}
				for _, parsed := range docs[t.Name()].Fields {
					f, ok := t.FieldByName(parsed.GoName)
					if !ok {
						continue
					}
					name := parsed.Name
					ft := describe(f.Type)
					value := reflect.Zero(f.Type).Interface()
					if ft.Kind == "long" {
						value = "0"
					}
					schema.Fields = append(schema.Fields, Field{Name: name, Type: ft, Desc: parsed.Comment, Nullable: f.Type.Kind() == reflect.Ptr || f.Type.Kind() == reflect.Map || f.Type.Kind() == reflect.Slice, HasDefault: true, DefaultValue: value, JSON: true})
				}
				m.Schemas[id] = schema
			}
			return Type{Kind: "object", Name: t.Name(), Ref: id}
		}
		return Type{Kind: "unknown", Name: t.String()}
	}
	for _, cmd := range protocol.GetMsgName2IdMapper() {
		t, _ := protocol.GetMessageType(cmd)
		describe(t)
	}
	groups := map[string][]Message{}
	for id, t := range types {
		name := t.Name()
		module := docs[name].SourceFile
		if module == "" {
			module = "common"
		}
		msg := Message{ID: id, Name: name, Module: module, Schema: id, Desc: docs[name].Comment, Direction: "VO"}
		for _, dir := range []string{"Req", "Res", "Push"} {
			if strings.HasPrefix(name, dir) {
				msg.Direction = strings.ToUpper(dir)
			}
		}
		if cmd, err := protocol.GetMessageCmdFromType(reflect.PointerTo(t)); err == nil {
			msg.Cmd = &cmd
			if routes != nil && module != "cross" {
				if handler, ok := routes.Handlers[cmd]; ok {
					msg.Testable = msg.Direction == "REQ"
					mt := handler.Method.Type
					if mt.NumOut() > 0 {
						res := mt.Out(0)
						for res.Kind() == reflect.Ptr {
							res = res.Elem()
						}
						if res.Kind() == reflect.Struct {
							msg.Response = res.PkgPath() + "." + res.Name()
						}
					}
				}
			}
		}
		groups[module] = append(groups[module], msg)
	}
	for id, messages := range groups {
		sort.Slice(messages, func(i, j int) bool {
			if (messages[i].Cmd == nil) != (messages[j].Cmd == nil) {
				return messages[i].Cmd != nil
			}
			if messages[i].Cmd != nil && messages[j].Cmd != nil && *messages[i].Cmd != *messages[j].Cmd {
				return *messages[i].Cmd < *messages[j].Cmd
			}
			return messages[i].Name < messages[j].Name
		})
		m.Modules = append(m.Modules, Module{ID: id, Name: id, Protocols: messages})
	}
	sort.Slice(m.Modules, func(i, j int) bool { return m.Modules[i].ID < m.Modules[j].ID })
	return m
}
