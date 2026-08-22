package consume

import (
	"reflect"

	"github.com/forfun/gforgame/data"
	domainconsume "github.com/forfun/gforgame/internal/domain/consume"
)

func init() {
	data.ConvertRegistryInstance.RegisterFunc(
		reflect.TypeOf((*domainconsume.Consume)(nil)).Elem(),
		func(source string) (any, error) {
			return ParseConsume(source), nil
		})
}
