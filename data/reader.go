package data

import "io"

type DataReader interface {
	Read(io.Reader, interface{}) ([]interface{}, error)
}
