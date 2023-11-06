# Changelog

## 0.1.0 (2023-11-06)

- Add: A `Terminal` can now be defined with more trivial callbacks. The
  `*TerminalFactory` now has `FuncByte()` and `FuncRune()`:
  ```go
  NewTerm(tInt, "int").FuncByte(isDigit, bytesToInt)
  NewWhitespace().FuncRune(unicode.IsSpace)

  func isDigit(b byte) bool              { return b >= '0' && b <= '9' }
  func bytesToInt(b []byte) (int, error) { return strconv.Atoi(string(b)) }
  ```

## 0.0.1 (2023-11-02)

First release.
