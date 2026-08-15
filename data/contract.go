package data

import "reflect"

type CellHeader struct {
	Column string
	Field  reflect.Value
}

type CellColumn struct {
	Header CellHeader
	Value  string
}
