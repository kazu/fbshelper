// Package dump ... hexdump with historical diff
//
// Copyright (c) 20121-2021 Kazuhisa TAKEI. All rights reserved.
// https://github.com/kazu/fbshelper
// See the included LICENSE file for license details.
//
package dump

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// TypeInfo ... type of DumpInfomation
type TypeInfo byte

const (
	// TypeVariable is caller variable
	TypeVariable TypeInfo = 1 << iota
	// TypeMethod is called method
	TypeMethod
	// TypeParam is parameter of method
	TypeParam
)

// DumpFunc ... function to hexDump
type DumpFunc func(interface{}) string

// Histories ... colliection of Histories
type Histories struct {
	Variable interface{}
	DumpFn   DumpFunc
	logs     []History
	fInfo    []FuncInfo
}

// History ... one history for dump
type History struct {
	Name string
	Dump string
}

// FuncInfo ... function info
type FuncInfo struct {
	Name  string
	TInfo TypeInfo
}

// New ... make Histories
func New(variable interface{}, dump DumpFunc, info []FuncInfo) Histories {
	return Histories{
		Variable: variable,
		DumpFn:   dump,
		fInfo:    info,
	}
}

func (histories *Histories) Dump(Fmt string, v ...interface{}) {

	histories.logs = append(histories.logs, History{
		Name: fmt.Sprintf(Fmt, v...),
		Dump: histories.DumpFn(histories.Variable),
	})
}

func (histories Histories) Strings() (result []string) {

	var b strings.Builder

	var variable string
	var method string
	params := []string{}

	result = []string{}

	for _, info := range histories.fInfo {

		switch info.TInfo {
		case TypeVariable:
			variable = info.Name
		case TypeMethod:
			method = info.Name
		case TypeParam:
			params = append(params, info.Name)
		}
	}

	fmt.Fprintf(&b, "%s.%s(%+v)\n", variable, method, params)

	for i := range histories.logs {
		log := histories.logs[i]
		var b strings.Builder
		if i == 0 {
			fmt.Fprintf(&b, "%s\n%s\n", log.Name, log.Dump)
			result = append(result, b.String())
			continue
		}
		fmt.Fprintln(&b, histories.logs[i].Name)
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(histories.logs[i-1].Dump, histories.logs[i].Dump, false)
		fmt.Fprint(&b, "\n"+dmp.DiffText2(diffs)+"\n")
		result = append(result, b.String())
	}

	return

}

func (histories Histories) String() string {
	var b strings.Builder

	for _, log := range histories.Strings() {
		fmt.Fprint(&b, log)
	}
	return b.String()
}
