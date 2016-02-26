package buildr

import (
	"fmt"
	"strings"
	"os"
	"time"
)

const (
	SHORTNAME = 50
)

func short(name string) string {
	ln := len(name)
	if ln <= SHORTNAME {
		return name
	}
	nm := name[SHORTNAME-4:]
	return nm + "..."
}

type targetI interface {
	modifiedSince(tm time.Time) bool
	name() string
	Run() bool
}

type FileTarget struct {
	deptab map[string]targetI
	depends []targetI
	files   []string
	makef   func(... targetI) bool
}

func File(path string) *FileTarget {
	return &FileTarget {
		deptab: map[string]targetI{},
		depends: []targetI{},
		files: []string{path},
	}
}

func Files(paths... string) *FileTarget {
	return &FileTarget {
		deptab: map[string]targetI{},
		depends: []targetI{},
		files: paths,
	}
}

func (t *FileTarget) Depends(targets... targetI) *FileTarget {
	t.depends = append(t.depends, targets...)
	for _, targ := range targets {
		t.deptab[targ.name()] = targ
	}
	return t
}

func (t *FileTarget) Make(run func(... targetI) bool) *FileTarget {
	t.makef = run
	return t
}

func (t *FileTarget) modifiedSince(tm time.Time) bool {
	for _, f := range t.files {
		info, err := os.Stat(f)
		if err != nil {
			return true
		}
		if info.ModTime().After(tm) {
			return true
		}
	}
	return false
}

func (t *FileTarget) name() string {
	return strings.Join(t.files, " ")
}

func (t *FileTarget) modTime() time.Time {
	tm := time.Unix(0, 0)
	for _, f := range t.files {
		info, err := os.Stat(f)
		if err != nil {
			return time.Now()
		}
		tm_ := info.ModTime()
		if tm_.After(tm) {
			tm = tm_
		}
	}
	return tm
}

func (t *FileTarget) Run() bool {
	tm := t.modTime()
	modified := false
	fmt.Println("[buildr] Make target `" + short(t.name()) + "`...")
	defer fmt.Println("[buildr] Done `" + short(t.name()) + "`")

	for name, targ := range t.deptab {
		if targ.modifiedSince(tm) {
			if !targ.Run() {
				return false
			}
			modified = true
		}
	}

	if modified {
		return t.makef(t.depends...)
	}
	return true
}

func (t *FileTarget) RunTarget(name string) bool {
	if targ, ok := t.deptab[name]; !ok {
		fmt.Println("[buildr] Cannot find target `" + short(name) + "`...")
		return false
	} else {
		return targ.Run()
	}
}
