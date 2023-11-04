package lr0

import (
	"testing"
)

var testTableItemsetGrammar = newGrammar(
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
			Is(nSum, tPlus, nVal).Do(calc3AnyNil).
			Is(nSum, tMinus, nVal).Do(calc3AnyNil).
			Is(nVal),
		NewNT(nVal, "Val").
			Is(tZero).
			Is(tOne),
	},
)
var (
	testTableItemsetRuleGoal     = testTableItemsetGrammar.Rule(0)
	testTableItemsetRuleSumPlus  = testTableItemsetGrammar.Rule(1)
	testTableItemsetRuleSumMinus = testTableItemsetGrammar.Rule(2)
	testTableItemsetRuleSumVal   = testTableItemsetGrammar.Rule(3)
	testTableItemsetRuleValZero  = testTableItemsetGrammar.Rule(4)
	testTableItemsetRuleValOne   = testTableItemsetGrammar.Rule(5)
)

func TestExpandAllPossibleTableItems(t *testing.T) {
	in0 := []tableItem{newTableItem(testTableItemsetGrammar.MainRule())}
	out0 := expandAllPossibleTableItems(in0, testTableItemsetGrammar)
	if len(out0) != 6 {
		t.Errorf("0: len is %v", len(out0))
	}

	var in1 []tableItem
	for _, r := range testTableItemsetGrammar.RulesFor(nSum) {
		in1 = append(in1, newTableItem(r))
	}
	out1 := expandAllPossibleTableItems(in1, testTableItemsetGrammar)
	if len(out1) != 5 {
		t.Errorf("1: len is %v", len(out1))
	}

	var in2 []tableItem
	for _, r := range testTableItemsetGrammar.RulesFor(nVal) {
		in2 = append(in2, newTableItem(r))
	}
	out2 := expandAllPossibleTableItems(in2, testTableItemsetGrammar)
	if len(out2) != 2 {
		t.Errorf("2: len is %v", len(out2))
	}
}

func TestTableItemset_HasFinalItem(t *testing.T) {
	in0 := []tableItem{
		newTableItem(testTableItemsetGrammar.MainRule()),
	}
	set0 := newTableItemset(in0, testTableItemsetGrammar)
	if set0.HasFinalItem() {
		t.Error("it must not have")
	}

	var in1 []tableItem
	for _, r := range testTableItemsetGrammar.RulesFor(nSum) {
		in1 = append(in1, newTableItem(r))
	}
	set1 := newTableItemset(in1, testTableItemsetGrammar)
	if set1.HasFinalItem() {
		t.Error("it must not have")
	}

	in2 := []tableItem{
		newTableItem(testTableItemsetGrammar.MainRule()).Shift(),
	}
	set2 := newTableItemset(in2, testTableItemsetGrammar)
	if !set2.HasFinalItem() {
		t.Error("it must have")
	}
}

func TestTableItemset_IsEqual(t *testing.T) {
	l := newLexer(
		NewTerm(tZero, "zero").Str("0"),
		NewTerm(tPlus, "plus").Str("+"),
		NewTerm(tMinus, "minus").Str("-"),
	)
	rulesSum := NewNT(nSum, "Sum").
		Is(nSum, tPlus, nVal).Do(calc3AnyNil).
		Is(nSum, tMinus, nVal).Do(calc3AnyNil).
		Is(nVal).
		GetRules(l)
	rulesVal := NewNT(nVal, "Val").
		Is(tZero).
		GetRules(l)
	r1 := rulesSum[0]
	r2 := rulesSum[1]
	r3 := rulesSum[2]
	r4 := rulesVal[0]

	s0 := tableItemset{items: []tableItem{
		newTableItem(r1),
		newTableItem(r2),
		newTableItem(r3),
		newTableItem(r4),
	}}
	s0same := tableItemset{items: []tableItem{
		newTableItem(r3),
		newTableItem(r2),
		newTableItem(r4),
		newTableItem(r1),
	}}
	if !s0.IsEqual(s0same) {
		t.Error("0: not same")
	}

	if s0.IsEqual(tableItemset{}) {
		t.Error("0a: same")
	}

	s1 := tableItemset{items: []tableItem{
		newTableItem(r1),
		newTableItem(r2).Shift(),
		newTableItem(r3),
		newTableItem(r4),
	}}
	if s0.IsEqual(s1) {
		t.Error("1: same")
	}
}

func TestTableItemset_GetNextItemsets(t *testing.T) {
	in0 := []tableItem{
		newTableItem(testTableItemsetGrammar.MainRule()),
	}
	set0 := newTableItemset(in0, testTableItemsetGrammar)

	set1 := set0.GetNextItemsets(testTableItemsetGrammar)
	if len(set1) != 4 {
		t.Fatalf("set1 len is %v", len(set1))
	}

	set1zero, ok := set1[tZero]
	if !ok {
		t.Fatal("set 1: no zero")
	}
	if !set1zero.IsEqual(
		tableItemset{[]tableItem{
			newTableItem(testTableItemsetRuleValZero).Shift(),
		}},
	) {
		t.Error("set1zero wrong")
	}
	if len(set1zero.GetNextItemsets(testTableItemsetGrammar)) != 0 {
		t.Error("set1zero next wrong")
	}

	set1one, ok := set1[tOne]
	if !ok {
		t.Fatal("set 1: no one")
	}
	if !set1one.IsEqual(
		tableItemset{[]tableItem{
			newTableItem(testTableItemsetRuleValOne).Shift(),
		}},
	) {
		t.Error("set1one wrong")
	}
	if len(set1one.GetNextItemsets(testTableItemsetGrammar)) != 0 {
		t.Error("set1one next wrong")
	}

	set1val, ok := set1[nVal]
	if !ok {
		t.Fatal("set 1: no val")
	}
	if !set1val.IsEqual(
		tableItemset{[]tableItem{
			newTableItem(testTableItemsetRuleSumVal).Shift(),
		}},
	) {
		t.Error("set1val wrong")
	}
	if len(set1val.GetNextItemsets(testTableItemsetGrammar)) != 0 {
		t.Error("set1val next wrong")
	}

	set1sum, ok := set1[nSum]
	if !ok {
		t.Fatal("set 1: no sum")
	}
	if !set1sum.IsEqual(
		tableItemset{[]tableItem{
			newTableItem(testTableItemsetRuleGoal).Shift(),
			newTableItem(testTableItemsetRuleSumPlus).Shift(),
			newTableItem(testTableItemsetRuleSumMinus).Shift(),
		}},
	) {
		t.Error("set1sum wrong")
	}

	set1sumNext := set1sum.GetNextItemsets(testTableItemsetGrammar)
	if len(set1sumNext) != 2 {
		t.Error("set1sumNext wrong")
	}

	set1sumPlus, ok := set1sumNext[tPlus]
	if !ok {
		t.Fatal("set1sumNext no plus")
	}
	if !set1sumPlus.IsEqual(
		tableItemset{[]tableItem{
			newTableItem(testTableItemsetRuleSumPlus).Shift().Shift(),
			newTableItem(testTableItemsetRuleValZero),
			newTableItem(testTableItemsetRuleValOne),
		}},
	) {
		t.Error("set1sumPlus wrong")
	}

	set1sumMinus, ok := set1sumNext[tMinus]
	if !ok {
		t.Fatal("set1sumNext no minus")
	}
	if !set1sumMinus.IsEqual(
		tableItemset{[]tableItem{
			newTableItem(testTableItemsetRuleSumMinus).Shift().Shift(),
			newTableItem(testTableItemsetRuleValZero),
			newTableItem(testTableItemsetRuleValOne),
		}},
	) {
		t.Error("set1sumMinus wrong")
	}
}
