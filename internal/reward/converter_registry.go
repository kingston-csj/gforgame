package reward

import (
	"reflect"

	"github.com/forfun/gforgame/data"
	domainreward "github.com/forfun/gforgame/internal/domain/reward"
)

func init() {
	data.ConvertRegistryInstance.RegisterFunc(
		reflect.TypeOf((*domainreward.Reward)(nil)).Elem(),
		func(source string) (any, error) {
			return ParseReward(source), nil
		})
}
