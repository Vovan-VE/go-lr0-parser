package lr0

import (
	"strconv"

	"github.com/pkg/errors"
)

const (
	tInt Id = iota + 1
	tZero
	tOne
	tPlus
	tMinus
	tMul
	tDiv
	tIdent
	tInc

	nVal
	nProd
	nSum
	nGoal
)

var errDivZero = errors.New("division by zero")

func matchIdentifier(state *State) (next *State, value any) {
	if state.IsEOF() {
		return
	}
	if next, _ = state.TakeByteFunc(isAlpha); next == nil {
		return
	}
	next, _ = next.TakeBytesFunc(isAlphaNum)
	value = string(state.BytesTo(next))
	return
}

func isDigit(b byte) bool        { return b >= '0' && b <= '9' }
func byteIsNotSpace(b byte) bool { return b != ' ' }
func runeIsNotSpace(r rune) bool { return r != ' ' }
func isAlphaNum(b byte) bool     { return isAlpha(b) || isDigit(b) }

func isAlpha(b byte) bool {
	switch {
	case b >= 'A' && b <= 'Z', b >= 'a' && b <= 'z', b == '_':
		return true
	default:
		return false
	}
}

func bytesToInt(b []byte) (int, error) { return strconv.Atoi(string(b)) }

func calc3AnyNil(any, any, any) (any, error)     { return nil, nil }
func calc3StrTrace(a, op, b string) (any, error) { return "(" + a + " " + op + " " + b + ")", nil }
func calc2IntSum(a, b int) int                   { return a + b }
func calc2IntSub(a, b int) int                   { return a - b }
