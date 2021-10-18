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
	"sync"

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
	logCh    chan int
	fInfo    []FuncInfo
	wg       sync.WaitGroup
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
		logCh:    make(chan int, 100),
	}
}

func (histories *Histories) Dump(Fmt string, v ...interface{}) {

	h := History{
		Name: fmt.Sprintf(Fmt, v...),
		Dump: histories.DumpFn(histories.Variable)}

	histories.logs = append(histories.logs, h)
	histories.logCh <- len(histories.logs) - 1
}

func (histories *Histories) DumpWithFlag(flag bool, Fmt string, v ...interface{}) {

	if !flag {
		return
	}
	histories.Dump(Fmt, v...)
}

func (histories *Histories) Strings() (result []string) {

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

func (histories *Histories) String() string {
	var b strings.Builder

	for _, log := range histories.Strings() {
		fmt.Fprint(&b, log+"\n")
	}
	return b.String()
}

func (histories *Histories) StringBy(i int) string {
	var b strings.Builder
	log := histories.logs[i]

	if i == 0 {
		fmt.Fprintf(&b, "%s\n%s\n", log.Name, log.Dump)
		return b.String()
	}

	fmt.Fprintln(&b, log.Name)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(histories.logs[i-1].Dump, histories.logs[i].Dump, false)
	fmt.Fprint(&b, "\n"+dmp.DiffText2(diffs)+"\n")
	return b.String()

}

func (histories *Histories) Finish() {

	histories.logCh <- -1

	histories.wg.Wait()
	close(histories.logCh)

}

func (histories *Histories) StreamOut(fn func(string)) {
	histories.wg.Add(1)

	go histories.streamOut(fn, &histories.wg)

}

func (histories *Histories) streamOut(wFn func(string), wg *sync.WaitGroup) {
	for i := range histories.logCh {
		if i == -1 {
			break
		}
		wFn(histories.StringBy(i))
	}
	wg.Done()
}
