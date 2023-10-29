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

var (
	testRuleGoal     = grammar.NewRuleMain(nGoal, []symbol.Id{nSum})
	testRuleSumPlus  = grammar.NewRule(nSum, []symbol.Id{nSum, tPlus, nVal})
	testRuleSumMinus = grammar.NewRule(nSum, []symbol.Id{nSum, tMinus, nVal})
	testRuleSumVal   = grammar.NewRule(nSum, []symbol.Id{nVal})
	testRuleValZero  = grammar.NewRule(nVal, []symbol.Id{tZero})
	testRuleValOne   = grammar.NewRule(nVal, []symbol.Id{tOne})
)
var testGrammar = grammar.New(
	[]lexer.Terminal{
		lexer.NewFixedStr(tZero, "0"),
		lexer.NewFixedStr(tOne, "1"),
		lexer.NewFixedStr(tPlus, "+"),
		lexer.NewFixedStr(tMinus, "-"),
	},
	[]grammar.Rule{
		testRuleGoal,
		testRuleSumPlus,
		testRuleSumMinus,
		testRuleSumVal,
		testRuleValZero,
		testRuleValOne,
	},
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
	r1 := grammar.NewRule(nSum, []symbol.Id{nSum, tPlus, nVal})
	r2 := grammar.NewRule(nSum, []symbol.Id{nSum, tMinus, nVal})
	r3 := grammar.NewRule(nSum, []symbol.Id{nVal})
	r4 := grammar.NewRule(nVal, []symbol.Id{tZero})

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
			newItem(testRuleValZero).Shift(),
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
			newItem(testRuleValOne).Shift(),
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
			newItem(testRuleSumVal).Shift(),
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
			newItem(testRuleGoal).Shift(),
			newItem(testRuleSumPlus).Shift(),
			newItem(testRuleSumMinus).Shift(),
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
			newItem(testRuleSumPlus).Shift().Shift(),
			newItem(testRuleValZero),
			newItem(testRuleValOne),
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
			newItem(testRuleSumMinus).Shift().Shift(),
			newItem(testRuleValZero),
			newItem(testRuleValOne),
		}},
	) {
		t.Error("set1sumMinus wrong")
	}
}
