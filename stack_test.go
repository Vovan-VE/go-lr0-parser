package lr0

import (
	"testing"
)

// Goal : Sum $
// Sum  : Sum plus Val
// Sum  : Sum minus Val
// Sum  : Val
// Val  : zero
// Val  : one

// grammar & table copy-paste from table tests
//	row 0 ---------
//		EOF: ACCEPT
//		terminals:
//			zero -> 1
//			one  -> 2
//		goto:
//			Val -> 3
//			Sum -> 4
//		rule: -
//	row 1 ---------
//		EOF: ACCEPT
//		terminals: -
//		goto: -
//		rule:
//			Val : zero
//	row 2 ---------
//		EOF: ACCEPT
//		terminals: -
//		goto: -
//		rule:
//			Val : one
//	row 3 ---------
//		EOF: ACCEPT
//		terminals: -
//		goto: -
//		rule:
//			Sum : Val
//	row 4 ---------
//		EOF: -
//		terminals:
//			plus  -> 5
//			minus -> 6
//		goto: -
//		rule:
//			Goal : Sum $
//	row 5 ---------
//		EOF: ACCEPT
//		terminals:
//			zero -> 1
//			one  -> 2
//		goto:
//			Val -> 7
//		rule: -
//	row 6 ---------
//		EOF: ACCEPT
//		terminals:
//			zero -> 1
//			one -> 2
//		goto:
//			Val -> 8
//		rule: -
//	row 7 ---------
//		EOF: ACCEPT
//		terminals: -
//		goto: -
//		rule:
//			Sum : Sum plus Val
//	row 8 ---------
//		EOF: ACCEPT
//		terminals: -
//		goto: -
//		rule:
//			Sum : Sum minus Val

func TestStack(t *testing.T) {
	testTable := newTable(newGrammar(
		[]Terminal{
			NewTerm(tZero, "zero").Str("0"),
			NewTerm(tOne, "one").Str("1"),
			NewTerm(tPlus, `"+"`).Str("+"),
			NewTerm(tMinus, `"-"`).Str("-"),
		},
		[]NonTerminalDefinition{
			NewNT(nGoal, "Goal").
				Main().
				Is(nSum),
			NewNT(nSum, "Sum").
				Is(nSum, tPlus, nVal).Do(calc3StrTrace).
				Is(nSum, tMinus, nVal).Do(calc3StrTrace).
				Is(nVal),
			NewNT(nVal, "Val").
				Is(tZero).
				Is(tOne),
		},
	))

	st := newStack(testTable)
	if st.si != 0 {
		t.Fatal("initial state ", st.si)
	}

	row1, ok1 := testTable.Row(0).TerminalAction(tZero)
	row2, ok2 := testTable.Row(0).TerminalAction(tOne)
	row3, ok3 := testTable.Row(0).GotoAction(nVal)
	row4, ok4 := testTable.Row(0).GotoAction(nSum)
	row5, ok5 := testTable.Row(row4).TerminalAction(tPlus)
	row6, ok6 := testTable.Row(row4).TerminalAction(tMinus)
	row7, ok7 := testTable.Row(row5).GotoAction(nVal)
	row8, ok8 := testTable.Row(row6).GotoAction(nVal)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
		t.Fatalf("not all ok: %#v", []bool{ok1, ok2, ok3, ok4, ok5, ok6, ok7, ok8})
	}

	// 1+0-1
	// ^
	// in 0 "1" shifts to 2
	st.Shift(row2, tOne, "1")
	// [{2 one "1"}]
	if len(st.items) != 1 {
		t.Errorf("a: items: %#v", st.items)
	}

	// 1+0-1
	// -^
	// in 2 "+" not expected - reduce
	if ok, err := st.Reduce(); !ok || err != nil {
		t.Fatalf("b: reduce %v, %v", ok, err)
	}
	// [{3 Val "1"}]
	if len(st.items) != 1 {
		t.Errorf("b: items: %#v", st.items)
	}
	if st.items[0] != (stackItem{state: row3, node: nVal, value: "1"}) {
		t.Errorf("b: state is %v", st.items[0])
	}
	// in 3 "+" not expected - reduce
	if ok, err := st.Reduce(); !ok || err != nil {
		t.Fatalf("c: reduce %v, %v", ok, err)
	}
	// [{4 Sum "1"}]
	if len(st.items) != 1 {
		t.Errorf("c: items: %#v", st.items)
	}
	if st.items[0] != (stackItem{state: row4, node: nSum, value: "1"}) {
		t.Errorf("c: state is %v", st.items[0])
	}
	// in 4 "+" shifts to 5
	st.Shift(row5, tPlus, "+")
	// [{4 Sum "1"}, {5 plus "+"}]
	if len(st.items) != 2 {
		t.Errorf("d: items: %#v", st.items)
	}

	// 1+0-1
	// --^
	// in 5 "0" shifts to 1
	st.Shift(row1, tZero, "0")
	// [{4 Sum "1"}, {5 plus "+"}, {1 zero "0"}]
	if len(st.items) != 3 {
		t.Errorf("e: items: %#v", st.items)
	}

	// 1+0-1
	// ---^
	// in 1 "-" not expected - reduce
	if ok, err := st.Reduce(); !ok || err != nil {
		t.Fatalf("f: reduce %v, %v", ok, err)
	}
	// [{4 Sum "1"}, {5 plus "+"}, {7 Val "0"}]
	if len(st.items) != 3 {
		t.Errorf("f: items: %#v", st.items)
	}
	if st.items[2] != (stackItem{state: row7, node: nVal, value: "0"}) {
		t.Errorf("f: state is %v", st.items[2])
	}
	// in 7 "-" not expected - reduce
	if ok, err := st.Reduce(); !ok || err != nil {
		t.Fatalf("g: reduce %v, %v", ok, err)
	}
	// [{4 Sum "(1 + 0)"}]
	if len(st.items) != 1 {
		t.Errorf("g: items: %#v", st.items)
	}
	if st.items[0] != (stackItem{state: row4, node: nSum, value: "(1 + 0)"}) {
		t.Errorf("g: state is %v", st.items[0])
	}
	// in 4 "-" shifts to 6
	st.Shift(row6, tMinus, "-")
	// [{4 Sum "(1 + 0)"}, {6 minus "-"}]
	if len(st.items) != 2 {
		t.Errorf("h: items: %#v", st.items)
	}

	// 1+0-1
	// ----^
	// in 6 "1" shifts to 2
	st.Shift(row2, tOne, "1")
	// [{4 Sum "(1 + 0)"}, {6 minus "-"}, {2 one "1"}]
	if len(st.items) != 3 {
		t.Errorf("i: items: %#v", st.items)
	}

	// 1+0-1
	// -----^
	// eof
	// in 2 no terminals expected and EOF valid - reduce
	if ok, err := st.Reduce(); !ok || err != nil {
		t.Fatalf("j: reduce %v, %v", ok, err)
	}
	// [{4 Sum "(1 + 0)"}, {6 minus "-"}, {8 Val "1"}]
	if len(st.items) != 3 {
		t.Errorf("j: items: %#v", st.items)
	}
	if st.items[2] != (stackItem{state: row8, node: nVal, value: "1"}) {
		t.Errorf("j: state is %v", st.items[2])
	}
	// in 8 no terminals expected and EOF valid - reduce and done
	if ok, err := st.Reduce(); !ok || err != nil {
		t.Fatalf("k: reduce %v, %v", ok, err)
	}
	// [{4 Sum "((1 + 0) - 1)"}]
	if len(st.items) != 1 {
		t.Errorf("k: items: %#v", st.items)
	}
	if st.items[0] != (stackItem{state: row4, node: nSum, value: "((1 + 0) - 1)"}) {
		t.Errorf("k: state is %v", st.items[0])
	}
	result := st.Done()
	if result != "((1 + 0) - 1)" {
		t.Errorf("result wrong: %q", result)
	}
	if len(st.items) != 0 {
		t.Errorf("result items: %#v", st.items)
	}
}
