package symbol

import (
	"github.com/pkg/errors"
)

var (
	// ErrDefine will be raised by panic in case of invalid definition
	ErrDefine = errors.New("invalid definition")
)
