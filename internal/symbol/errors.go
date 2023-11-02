package symbol

import (
	"github.com/pkg/errors"
)

var (
	// ErrDefine will be raised by panic in case of invalid definition
	ErrDefine = errors.New("invalid definition")
	// ErrInternal is an internal error which actually should not be raised, but
	// coded in some panic for debug purpose just in case
	ErrInternal = errors.New("internal error")
)
