package config

import (
	"github.com/forfun/gforgame/internal/domain/attr"
)

type HeroStageData struct {
	// 属性
	attr.AttributeItems `json:"-" excel:"-"`
	Id                  int32 `json:"id" excel:"id"`
	Stage               int32 `json:"stage" excel:"stage"`
	MaxLevel            int32 `json:"max_level" excel:"max_level"`
	Cost                int32 `json:"cost" excel:"cost"`
	Hp                  int32 `json:"hp" excel:"hp"`
	Attack              int32 `json:"attack" excel:"attack"`
	Defense             int32 `json:"defense" excel:"defense"`
	Speed               int32 `json:"speed" excel:"speed"`
}
