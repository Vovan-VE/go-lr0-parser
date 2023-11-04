package lr0

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

const stateFormatContext = 30

// State describes an immutable state of reading the underlying buffer at the
// given position
type State struct {
	source []byte
	at     int
}

// NewState creates new State for the given buffer `input` pointing to its start
func NewState(input []byte) *State {
	return &State{
		source: input,
	}
}

// to returns new State for the same buffer pointing to the given position `pos`
func (s *State) to(pos int) *State {
	pos = fixOffset(pos, len(s.source))
	if pos == s.at {
		return s
	}
	return &State{
		source: s.source,
		at:     pos,
	}
}

// IsEOF checks if the position is at EOF
func (s *State) IsEOF() bool {
	return s.at >= len(s.source)
}

// Len returns length of the underlying buffer
func (s *State) Len() int {
	return len(s.source)
}

// Offset returns the current offset
func (s *State) Offset() int {
	return s.at
}

// RestLen returns the rest unread length of the underlying buffer
func (s *State) RestLen() int {
	return len(s.source) - s.at
}

// RestBytes returns slice of rest bytes from the current position
func (s *State) RestBytes() []byte {
	return s.source[s.at:]
}

// BytesTo returns slice of underlying buffer from current position to the given
// State position
func (s *State) BytesTo(to *State) []byte {
	return s.BytesToOffset(to.at)
}

// BytesToOffset returns slice of underlying buffer from current position to the
// given position
func (s *State) BytesToOffset(offset int) []byte {
	to := fixOffset(offset, len(s.source))
	if to < s.at {
		panic(errors.Wrapf(ErrNegativeOffset, "from offset %v to backward offset %v", s.at, to))
	}
	return s.source[s.at:to]
}

// Byte returns a byte from current position
//
// panics at EOF
func (s *State) Byte() byte {
	if s.IsEOF() {
		panic(io.EOF)
	}
	return s.source[s.at]
}

// Rune returns a rune from the current position
//
// panics at EOF
func (s *State) Rune() (r rune, n int) {
	if s.IsEOF() {
		panic(io.EOF)
	}
	cut := s.cutUpTo(4)
	r, n = utf8.DecodeRune(cut)
	return
}

// FF returns new State for the same buffer at the position next n bytes to
// current
//
// n can also be negative, but will panic if refers to negative position
func (s *State) FF(n int) *State {
	return s.to(s.at + n)
}

// TakeByte returns a byte from the current position and new State with next
// position
//
// A combination of Byte() and FF(1).
func (s *State) TakeByte() (*State, byte) {
	b := s.Byte()
	return s.FF(1), b
}

// TakeByteFunc returns a byte from the current position and new State with next
// position only if the current byte match the given callback.
// If no match, returns `nil`
func (s *State) TakeByteFunc(valid func(byte) bool) (*State, byte) {
	b := s.Byte()
	if !valid(b) {
		return nil, 0
	}
	return s.FF(1), b
}

// TakeBytes returns a slice of bytes up to n length from the current position,
// truncated by EOF
func (s *State) TakeBytes(n int) (*State, []byte) {
	next := s.FF(n)
	return next, s.source[s.at:next.at]
}

// TakeBytesFunc return next State and slice of bytes which are valid by the
// result of the given func.
//
// At EOF the `valid` will not be called.
//
// If none valid bytes found, the result is the input state itself and nil
// slice.
func (s *State) TakeBytesFunc(valid func(byte) bool) (*State, []byte) {
	next := s
	for !next.IsEOF() && valid(next.Byte()) {
		next = next.FF(1)
	}
	if next.at == s.at {
		return s, nil
	}
	return next, s.source[s.at:next.at]
}

// TakeRune returns a rune from the current position and new State with next
// position
//
// A combination of Rune() and FF(n), where n is size of the rune read.
func (s *State) TakeRune() (*State, rune) {
	r, n := s.Rune()
	return s.FF(n), r
}

// TakeRunes returns a slice of runes up to n length from the current position,
// truncated by EOF
func (s *State) TakeRunes(n int) (*State, []rune) {
	next := s
	var rr []rune
	var r rune
	for i := 0; !next.IsEOF() && i < n; i++ {
		next, r = next.TakeRune()
		rr = append(rr, r)
	}
	return next, rr
}

// TakeRunesFunc return next State and slice of runes which are valid by the
// result of the given func.
//
// At EOF the `valid` will not be called.
//
// If none valid runes found, the result is the input state itself and nil
// slice.
func (s *State) TakeRunesFunc(valid func(rune) bool) (*State, []rune) {
	ret := s
	var rr []rune
	var r rune
	var n int
	for !ret.IsEOF() {
		r, n = ret.Rune()
		if !valid(r) {
			break
		}
		rr = append(rr, r)
		ret = ret.FF(n)
	}
	return ret, rr
}

// ExpectByte returns next state only if the given slice of bytes will match
// all corresponding bytes in the current position.
// Returns `nil` when no full match.
func (s *State) ExpectByte(b ...byte) *State {
	next, ok := s.ExpectByteOk(b...)
	if !ok {
		return nil
	}
	return next
}

// ExpectByteOk checks if bytes in the current position will match the given
// slice of bytes.
// The returned `State` is the last State checked.
// The returned `bool` will `true` only if all bytes matched.
func (s *State) ExpectByteOk(b ...byte) (next *State, ok bool) {
	next = s
	for !next.IsEOF() && len(b) != 0 && next.Byte() == b[0] {
		next = next.FF(1)
		b = b[1:]
	}
	ok = len(b) == 0
	return
}

// TODO: ExpectRune(r... rune) *State
// TODO: ExpectRuneOk(r... rune) (next *State, ok bool)

// cutUpTo returns a slice of buffer starting from current position and with
// size up to n bytes, truncated by EOF
func (s *State) cutUpTo(n int) []byte {
	return s.BytesToOffset(s.at + n)
}

func fixOffset(offset, length int) int {
	if offset < 0 {
		panic(ErrNegativeOffset)
	}
	if offset > length {
		return length
	}
	return offset
}

func (s *State) getBefore() string {
	from := s.at
	rest := stateFormatContext
	for from > 0 && rest > 0 {
		from--
		if utf8.RuneStart(s.source[from]) {
			rest--
		}
	}
	return string(s.source[from:s.at])
}
func (s *State) getAfter() string {
	to := s.at
	L := len(s.source)
	started := 0
	for to < L {
		if utf8.RuneStart(s.source[to]) {
			started++
			if started > stateFormatContext {
				break
			}
		}
		to++
	}
	return string(s.BytesToOffset(to))
}

func (s *State) String() string {
	var str string
	if before := ctlCharReplacer.Replace(s.getBefore()); before != "" {
		str += "⟪" + before + "⟫"
	}
	str += "⏵"
	if after := ctlCharReplacer.Replace(s.getAfter()); after != "" {
		str += "⟪" + after + "⟫"
	} else {
		str += "<EOF>"
	}
	return str
}

func (s *State) Format(st fmt.State, verb rune) {
	before := ctlCharReplacer.Replace(s.getBefore())
	after := ctlCharReplacer.Replace(s.getAfter())
	switch verb {
	case 'v':
		if st.Flag('+') {
			if before != "" {
				io.WriteString(st, before)
			}
			if after != "" {
				io.WriteString(st, after)
			} else {
				io.WriteString(st, "<EOF>")
			}
			io.WriteString(st, "\n")
			if b := utf8.RuneCount([]byte(before)); b > 0 {
				io.WriteString(st, strings.Repeat("-", b))
			}
			io.WriteString(st, "^\n")
			return
		}
		fallthrough
	case 's', 'q':
		if before != "" {
			io.WriteString(st, "⟪"+before+"⟫")
		}
		io.WriteString(st, "⏵")
		if after != "" {
			io.WriteString(st, "⟪"+after+"⟫")
		} else {
			io.WriteString(st, "<EOF>")
		}
	}
}

var ctlCharReplacer = strings.NewReplacer(
	"\x00", "␀", "\x01", "␁", "\x02", "␂", "\x03", "␃",
	"\x04", "␄", "\x05", "␅", "\x06", "␆", "\x07", "␇",
	"\x08", "␈", "\x09", "␉", "\x0A", "␊", "\x0B", "␋",
	"\x0C", "␌", "\x0D", "␍", "\x0E", "␎", "\x0F", "␏",
	"\x10", "␐", "\x11", "␑", "\x12", "␒", "\x13", "␓",
	"\x14", "␔", "\x15", "␕", "\x16", "␖", "\x17", "␗",
	"\x18", "␘", "\x19", "␙", "\x1A", "␚", "\x1B", "␛",
	"\x1C", "␜", "\x1D", "␝", "\x1E", "␞", "\x1F", "␟",
	" ", "␠", "\x7F", "␡",
)
