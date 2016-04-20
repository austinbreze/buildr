package buildr

import (
	"bytes"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type Func struct {
	Start, End int
	Head       string
	Code       string
}

func makeTable(src string) map[string]*Func {
	lines := strings.Split(src, "\n")
	cur := ""
	tab := map[string]*Func{}
	for i, line := range lines {
		line = strings.TrimLeft(line, "\t\n\r")
		if !strings.HasPrefix(line, "func ") {
			continue
		}
		if len(cur) > 0 {
			f := tab[cur]
			f.End = i - 1
			f.Code = strings.Join(lines[f.Start:f.End], "\n")
		}
		cur = strings.Split(line, "{")[0]
		tab[cur] = &Func{
			Start: i - 1,
			Head:  cur,
		}
	}
	if len(cur) > 0 {
		f := tab[cur]
		f.End = len(lines)
		f.Code = strings.Join(lines[f.Start:f.End], "\n")
	}
	return tab
}

func ExtendBlank(blank, ext string) (string, bool) {
	fext, err := format.Source([]byte("package main\n" + ext))
	if err != nil {
		log.Fatalln(err)
		return "", false
	}

	etab := makeTable(string(fext))
	btab := makeTable(blank)
	rslt := ""

	for h, f := range etab {
		if _, ok := btab[h]; !ok {
			rslt += "\n" + f.Code
		}
	}
	return rslt, true
}

func ExtendBlankFile(path string, f func(io.Writer) bool) bool {
	blank, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	ext := &bytes.Buffer{}
	if !f(ext) {
		return false
	}
	tail, ok := ExtendBlank(string(blank), string(ext.Bytes()))
	if !ok {
		return false
	}
	if len(strings.TrimSpace(tail)) == 0 {
		return true
	}
	writeTail := func(w io.Writer) bool {
		_, err := w.Write([]byte(tail))
		if err != nil {
			log.Fatalln(err)
			return false
		}
		return true
	}
	if !AppendFile(path, writeTail) {
		return false
	}
	return true
}
