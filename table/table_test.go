package table

import (
	"testing"

	"github.com/vovan-ve/go-lr0-parser/symbol"
)

func TestNew(t *testing.T) {
	const rowsCount = 9

	tbl := New(testGrammar).(*table)
	if len(tbl.rows) != rowsCount {
		t.Errorf("rows are %v", len(tbl.rows))
	}

	// row 0 (Goal : > Sum) ==================================================

	row0 := tbl.rows[0]
	if row0.AcceptEof() {
		t.Error("row0 eof")
	}
	if row0.ReduceRule() != nil {
		t.Error("row0 reduce rule")
	}
	if len(row0.terminals) != 2 {
		t.Errorf("row0.terminals: %#v", row0.terminals)
	}
	var (
		ok        bool
		next0zero StateIndex
		next0one  StateIndex
		next0val  StateIndex
		next0sum  StateIndex
	)
	if next0zero, ok = row0.terminals[tZero]; !ok || next0zero < 0 || next0zero >= rowsCount {
		t.Errorf("row0.terminals zero: %v, %v", next0zero, ok)
	}
	if next0one, ok = row0.terminals[tOne]; !ok || next0one < 0 || next0one >= rowsCount {
		t.Errorf("row0.terminals one: %v, %v", next0one, ok)
	}
	if len(row0.gotos) != 2 {
		t.Errorf("row0.gotos: %#v", row0.gotos)
	}
	if next0val, ok = row0.gotos[nVal]; !ok || next0val < 0 || next0val >= rowsCount {
		t.Errorf("row0.gotos val: %v, %v", next0val, ok)
	}
	if next0sum, ok = row0.gotos[nSum]; !ok || next0sum < 0 || next0sum >= rowsCount {
		t.Errorf("row0.gotos sum: %v, %v", next0sum, ok)
	}
	if len(map[StateIndex]struct{}{next0zero: {}, next0one: {}, next0val: {}, next0sum: {}}) != 4 {
		t.Errorf("row0 actions are not uniq: %v, %v, %v, %v", next0zero, next0one, next0val, next0sum)
	}

	// row 1 (Val : zero >) ==========================================

	rowValZero := tbl.rows[next0zero]
	if !rowValZero.IsReduceOnly() {
		t.Errorf("rowValZero not reduce only: %#v", rowValZero)
	}
	if rowValZero.ReduceRule() != testRuleImplValZero {
		t.Error("rowValZero reduce rule wrong")
	}

	// row 2 (Val : one >) ==========================================

	rowValOne := tbl.rows[next0one]
	if !rowValOne.IsReduceOnly() {
		t.Errorf("rowValOne not reduce only: %#v", rowValOne)
	}
	if rowValOne.ReduceRule() != testRuleImplValOne {
		t.Error("rowValOne reduce rule wrong")
	}

	// row 3 (Sum : Val >) ==========================================

	rowSumVal := tbl.rows[next0val]
	if !rowSumVal.IsReduceOnly() {
		t.Errorf("rowValOne not reduce only: %#v", rowSumVal)
	}
	if rowSumVal.ReduceRule() != testRuleImplSumVal {
		t.Error("rowValOne reduce rule wrong")
	}

	// row 4 (Goal : Sum > $) =======================================

	rowGoal := tbl.rows[next0sum]
	if !rowGoal.AcceptEof() {
		t.Error("rowGoal not eof")
	}
	if rowGoal.ReduceRule() != testRuleImplGoal {
		t.Error("rowGoal reduce rule wrong")
	}

	var (
		nextGoalPlus  StateIndex
		nextGoalMinus StateIndex
	)
	if len(rowGoal.terminals) != 2 {
		t.Errorf("rowGoal.terminals: %#v", rowGoal.terminals)
	}
	if nextGoalPlus, ok = rowGoal.terminals[tPlus]; !ok || nextGoalPlus < 0 || nextGoalPlus >= rowsCount {
		t.Errorf("rowGoal.terminals plus: %v, %v", nextGoalPlus, ok)
	}
	if nextGoalMinus, ok = rowGoal.terminals[tMinus]; !ok || nextGoalMinus < 0 || nextGoalMinus >= rowsCount {
		t.Errorf("rowGoal.terminals minus: %v, %v", nextGoalMinus, ok)
	}
	if len(rowGoal.gotos) != 0 {
		t.Errorf("rowGoal.gotos: %#v", rowGoal.gotos)
	}
	if len(map[StateIndex]struct{}{nextGoalPlus: {}, nextGoalMinus: {}}) != 2 {
		t.Errorf("rowGoal actions are not uniq: %v, %v", nextGoalPlus, nextGoalMinus)
	}

	// row 5 (Sum : Sum plus > Val) =====================================

	rowSumPlus := tbl.rows[nextGoalPlus]
	if rowSumPlus.AcceptEof() {
		t.Error("rowSumPlus eof")
	}
	if rowSumPlus.ReduceRule() != nil {
		t.Error("rowSumPlus reduce rule wrong")
	}
	var (
		nextPlusZero StateIndex
		nextPlusOne  StateIndex
		nextPlusVal  StateIndex
	)
	if len(rowSumPlus.terminals) != 2 {
		t.Errorf("rowSumPlus.terminals: %#v", rowSumPlus.terminals)
	}
	if nextPlusZero, ok = rowSumPlus.terminals[tZero]; !ok || nextPlusZero < 0 || nextPlusZero >= rowsCount {
		t.Errorf("rowSumPlus.terminals zero: %v, %v", nextPlusZero, ok)
	}
	if nextPlusOne, ok = rowSumPlus.terminals[tOne]; !ok || nextPlusOne < 0 || nextPlusOne >= rowsCount {
		t.Errorf("rowSumPlus.terminals one: %v, %v", nextPlusOne, ok)
	}
	if len(rowSumPlus.gotos) != 1 {
		t.Errorf("rowSumPlus.gotos: %#v", rowSumPlus.gotos)
	}
	if nextPlusVal, ok = rowSumPlus.gotos[nVal]; !ok || nextPlusVal < 0 || nextPlusVal >= rowsCount {
		t.Errorf("rowSumPlus.gotos val: %v, %v", nextPlusVal, ok)
	}
	if len(map[StateIndex]struct{}{nextPlusZero: {}, nextPlusOne: {}, nextPlusVal: {}}) != 3 {
		t.Errorf("rowSumPlus actions are not uniq: %v, %v, %v", nextPlusZero, nextPlusOne, nextPlusVal)
	}

	// row 6 (Sum : Sum minus > Val) =====================================

	rowSumMinus := tbl.rows[nextGoalMinus]
	if rowSumMinus.AcceptEof() {
		t.Error("rowSumMinus eof")
	}
	if rowSumMinus.ReduceRule() != nil {
		t.Error("rowSumMinus reduce rule wrong")
	}
	var (
		nextMinusZero StateIndex
		nextMinusOne  StateIndex
		nextMinusVal  StateIndex
	)
	if len(rowSumMinus.terminals) != 2 {
		t.Errorf("rowSumMinus.terminals: %#v", rowSumMinus.terminals)
	}
	if nextMinusZero, ok = rowSumMinus.terminals[tZero]; !ok || nextMinusZero < 0 || nextMinusZero >= rowsCount {
		t.Errorf("rowSumMinus.terminals zero: %v, %v", nextMinusZero, ok)
	}
	if nextMinusOne, ok = rowSumMinus.terminals[tOne]; !ok || nextMinusOne < 0 || nextMinusOne >= rowsCount {
		t.Errorf("rowSumMinus.terminals one: %v, %v", nextMinusOne, ok)
	}
	if len(rowSumMinus.gotos) != 1 {
		t.Errorf("rowSumMinus.gotos: %#v", rowSumMinus.gotos)
	}
	if nextMinusVal, ok = rowSumMinus.gotos[nVal]; !ok || nextMinusVal < 0 || nextMinusVal >= rowsCount {
		t.Errorf("rowSumMinus.gotos val: %v, %v", nextMinusVal, ok)
	}
	if len(map[StateIndex]struct{}{nextMinusZero: {}, nextMinusOne: {}, nextMinusVal: {}}) != 3 {
		t.Errorf("rowSumMinus actions are not uniq: %v, %v, %v", nextMinusZero, nextMinusOne, nextMinusVal)
	}

	// row 7 (Sum : Sum plus Val >) =============================

	rowSumPlusVal := tbl.rows[nextPlusVal]
	if !rowSumPlusVal.IsReduceOnly() {
		t.Errorf("rowSumPlusVal not reduce only: %#v", rowSumPlusVal)
	}
	if rowSumPlusVal.ReduceRule() != testRuleImplSumPlus {
		t.Error("rowSumPlusVal reduce rule wrong")
	}

	// row 8 (Sum : Sum minus Val >) =============================

	rowSumMinusVal := tbl.rows[nextMinusVal]
	if !rowSumMinusVal.IsReduceOnly() {
		t.Errorf("rowSumMinusVal not reduce only: %#v", rowSumMinusVal)
	}
	if rowSumMinusVal.ReduceRule() != testRuleImplSumMinus {
		t.Error("rowSumMinusVal reduce rule wrong")
	}

	if d := tbl.dump(testGrammar.(symbol.Registry)); d != expectTableDump {
		t.Error("table dump is:\n", d)
	}
}

const expectTableDump = `====[ table ]====
row 0 ---------
	EOF: ACCEPT
	terminals:
		zero -> 1
		one -> 2
	goto:
		Val -> 3
		Sum -> 4
	rule: -
row 1 ---------
	EOF: ACCEPT
	terminals: -
	goto: -
	rule:
		Val : zero
row 2 ---------
	EOF: ACCEPT
	terminals: -
	goto: -
	rule:
		Val : one
row 3 ---------
	EOF: ACCEPT
	terminals: -
	goto: -
	rule:
		Sum : Val
row 4 ---------
	EOF: -
	terminals:
		"+" -> 5
		"-" -> 6
	goto: -
	rule:
		Goal : Sum $
row 5 ---------
	EOF: ACCEPT
	terminals:
		zero -> 1
		one -> 2
	goto:
		Val -> 7
	rule: -
row 6 ---------
	EOF: ACCEPT
	terminals:
		zero -> 1
		one -> 2
	goto:
		Val -> 8
	rule: -
row 7 ---------
	EOF: ACCEPT
	terminals: -
	goto: -
	rule:
		Sum : Sum "+" Val
row 8 ---------
	EOF: ACCEPT
	terminals: -
	goto: -
	rule:
		Sum : Sum "-" Val
=================
`
