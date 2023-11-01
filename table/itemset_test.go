package table

import (
	"testing"

	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

const (
	tZero symbol.Id = iota + 1
	tOne
	tPlus
	tMinus

	nVal
	nSum
	nGoal
)

var testGrammar = grammar.New(
	[]lexer.Terminal{
		lexer.NewTerm(tZero, "zero").Str("0"),
		lexer.NewTerm(tOne, "one").Str("1"),
		lexer.NewTerm(tPlus, `"+"`).Str("+"),
		lexer.NewTerm(tMinus, `"-"`).Str("-"),
	},
	[]grammar.NonTerminalDefinition{
		grammar.NewNT(nGoal, "Goal").
			Main().
			Is(nSum),
		grammar.NewNT(nSum, "Sum").
			Is(nSum, tPlus, nVal).Do(calcStub3).
			Is(nSum, tMinus, nVal).Do(calcStub3).
			Is(nVal),
		grammar.NewNT(nVal, "Val").
			Is(tZero).
			Is(tOne),
	},
)
var (
	testRuleImplGoal     = testGrammar.RuleImpl(0)
	testRuleImplSumPlus  = testGrammar.RuleImpl(1)
	testRuleImplSumMinus = testGrammar.RuleImpl(2)
	testRuleImplSumVal   = testGrammar.RuleImpl(3)
	testRuleImplValZero  = testGrammar.RuleImpl(4)
	testRuleImplValOne   = testGrammar.RuleImpl(5)
)

func TestGetAllPossibleItems(t *testing.T) {
	in0 := []item{newItem(testGrammar.MainRule())}
	out0 := getAllPossibleItems(in0, testGrammar)
	if len(out0) != 6 {
		t.Errorf("0: len is %v", len(out0))
	}

	var in1 []item
	for _, r := range testGrammar.RulesFor(nSum) {
		in1 = append(in1, newItem(r))
	}
	out1 := getAllPossibleItems(in1, testGrammar)
	if len(out1) != 5 {
		t.Errorf("1: len is %v", len(out1))
	}

	var in2 []item
	for _, r := range testGrammar.RulesFor(nVal) {
		in2 = append(in2, newItem(r))
	}
	out2 := getAllPossibleItems(in2, testGrammar)
	if len(out2) != 2 {
		t.Errorf("2: len is %v", len(out2))
	}
}

func TestItemset_HasFinalItem(t *testing.T) {
	in0 := []item{
		newItem(testGrammar.MainRule()),
	}
	set0 := newItemset(in0, testGrammar)
	if set0.HasFinalItem() {
		t.Error("it must not have")
	}

	var in1 []item
	for _, r := range testGrammar.RulesFor(nSum) {
		in1 = append(in1, newItem(r))
	}
	set1 := newItemset(in1, testGrammar)
	if set1.HasFinalItem() {
		t.Error("it must not have")
	}

	in2 := []item{
		newItem(testGrammar.MainRule()).Shift(),
	}
	set2 := newItemset(in2, testGrammar)
	if !set2.HasFinalItem() {
		t.Error("it must have")
	}
}

func TestItemset_IsEqual(t *testing.T) {
	l := lexer.New(
		lexer.NewTerm(tZero, "zero").Str("0"),
		lexer.NewTerm(tPlus, "plus").Str("+"),
		lexer.NewTerm(tMinus, "minus").Str("-"),
	)
	rulesSum := grammar.NewNT(nSum, "Sum").
		Is(nSum, tPlus, nVal).Do(calcStub3).
		Is(nSum, tMinus, nVal).Do(calcStub3).
		Is(nVal).
		GetRules(l)
	rulesVal := grammar.NewNT(nVal, "Val").
		Is(tZero).
		GetRules(l)
	r1 := rulesSum[0]
	r2 := rulesSum[1]
	r3 := rulesSum[2]
	r4 := rulesVal[0]

	s0 := itemset{items: []item{
		newItem(r1),
		newItem(r2),
		newItem(r3),
		newItem(r4),
	}}
	s0same := itemset{items: []item{
		newItem(r3),
		newItem(r2),
		newItem(r4),
		newItem(r1),
	}}
	if !s0.IsEqual(s0same) {
		t.Error("0: not same")
	}

	if s0.IsEqual(itemset{}) {
		t.Error("0a: same")
	}

	s1 := itemset{items: []item{
		newItem(r1),
		newItem(r2).Shift(),
		newItem(r3),
		newItem(r4),
	}}
	if s0.IsEqual(s1) {
		t.Error("1: same")
	}
}

func TestItemset_GetNextItemsets(t *testing.T) {
	in0 := []item{
		newItem(testGrammar.MainRule()),
	}
	set0 := newItemset(in0, testGrammar)

	set1 := set0.GetNextItemsets(testGrammar)
	if len(set1) != 4 {
		t.Fatalf("set1 len is %v", len(set1))
	}

	set1zero, ok := set1[tZero]
	if !ok {
		t.Fatal("set 1: no zero")
	}
	if !set1zero.IsEqual(
		itemset{[]item{
			newItem(testRuleImplValZero).Shift(),
		}},
	) {
		t.Error("set1zero wrong")
	}
	if len(set1zero.GetNextItemsets(testGrammar)) != 0 {
		t.Error("set1zero next wrong")
	}

	set1one, ok := set1[tOne]
	if !ok {
		t.Fatal("set 1: no one")
	}
	if !set1one.IsEqual(
		itemset{[]item{
			newItem(testRuleImplValOne).Shift(),
		}},
	) {
		t.Error("set1one wrong")
	}
	if len(set1one.GetNextItemsets(testGrammar)) != 0 {
		t.Error("set1one next wrong")
	}

	set1val, ok := set1[nVal]
	if !ok {
		t.Fatal("set 1: no val")
	}
	if !set1val.IsEqual(
		itemset{[]item{
			newItem(testRuleImplSumVal).Shift(),
		}},
	) {
		t.Error("set1val wrong")
	}
	if len(set1val.GetNextItemsets(testGrammar)) != 0 {
		t.Error("set1val next wrong")
	}

	set1sum, ok := set1[nSum]
	if !ok {
		t.Fatal("set 1: no sum")
	}
	if !set1sum.IsEqual(
		itemset{[]item{
			newItem(testRuleImplGoal).Shift(),
			newItem(testRuleImplSumPlus).Shift(),
			newItem(testRuleImplSumMinus).Shift(),
		}},
	) {
		t.Error("set1sum wrong")
	}

	set1sumNext := set1sum.GetNextItemsets(testGrammar)
	if len(set1sumNext) != 2 {
		t.Error("set1sumNext wrong")
	}

	set1sumPlus, ok := set1sumNext[tPlus]
	if !ok {
		t.Fatal("set1sumNext no plus")
	}
	if !set1sumPlus.IsEqual(
		itemset{[]item{
			newItem(testRuleImplSumPlus).Shift().Shift(),
			newItem(testRuleImplValZero),
			newItem(testRuleImplValOne),
		}},
	) {
		t.Error("set1sumPlus wrong")
	}

	set1sumMinus, ok := set1sumNext[tMinus]
	if !ok {
		t.Fatal("set1sumNext no minus")
	}
	if !set1sumMinus.IsEqual(
		itemset{[]item{
			newItem(testRuleImplSumMinus).Shift().Shift(),
			newItem(testRuleImplValZero),
			newItem(testRuleImplValOne),
		}},
	) {
		t.Error("set1sumMinus wrong")
	}
}
