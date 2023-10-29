package actions

import (
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

type Actions interface {
}

type ruleRef struct {
	subject symbol.Id
	tag     symbol.Tag
}
