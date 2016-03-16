package buildr

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

type TargetI interface {
	modifiedSince(tm time.Time) bool
	name() string
	Build() bool
}

//FileTarget

type FileTarget struct {
	deptab  map[string]TargetI
	depends []TargetI
	files   []string
	makef   func(...TargetI) bool
}

func File(path string) *FileTarget {
	return &FileTarget{
		deptab:  map[string]TargetI{},
		depends: []TargetI{},
		files:   []string{path},
		makef:   func(...TargetI) bool { return true },
	}
}

func Files(paths ...string) *FileTarget {
	return &FileTarget{
		deptab:  map[string]TargetI{},
		depends: []TargetI{},
		files:   paths,
		makef:   func(...TargetI) bool { return true },
	}
}

func (t *FileTarget) Depends(targets ...TargetI) *FileTarget {
	t.depends = append(t.depends, targets...)
	for _, targ := range targets {
		t.deptab[targ.name()] = targ
	}
	return t
}

func (t *FileTarget) Make(run func(...TargetI) bool) *FileTarget {
	t.makef = run
	return t
}

func (t *FileTarget) modifiedSince(tm time.Time) bool {
	for _, f := range t.files {
		info, err := os.Stat(f)
		if err != nil {
			log.Println("[buildr]", err)
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
			log.Println("[buildr]", err)
			return tm
		}
		tm_ := info.ModTime()
		if tm_.After(tm) {
			tm = tm_
		}
	}
	return tm
}

func (t *FileTarget) Build() bool {
	tm := t.modTime()
	modified := false
	if len(t.depends) == 0 {
		modified = true
	}

	for _, targ := range t.deptab {
		if !targ.Build() {
			return false
		}
		if targ.modifiedSince(tm) {
			modified = true
		}
	}

	if modified {
		fmt.Println("[buildr] Make target `" + short(t.name()) + "`...")
		defer fmt.Println("[buildr] Done `" + short(t.name()) + "`")
		return t.makef(t.depends...)
	}
	return true
}

func (t *FileTarget) BuildTarget(name string) bool {
	if targ, ok := t.deptab[name]; !ok {
		fmt.Println("[buildr] Cannot find target `" + short(name) + "`...")
		return false
	} else {
		return targ.Build()
	}
}

//GlobTarget

type GlobTarget struct {
	deptab  map[string]TargetI
	depends []TargetI
	mask    string
	makef   func(...TargetI) bool
}

func Glob(mask string) *GlobTarget {
	return &GlobTarget{
		deptab:  map[string]TargetI{},
		depends: []TargetI{},
		mask:    mask,
		makef:   func(...TargetI) bool { return true },
	}
}

func (t *GlobTarget) Depends(targets ...TargetI) *GlobTarget {
	t.depends = append(t.depends, targets...)
	for _, targ := range targets {
		t.deptab[targ.name()] = targ
	}
	return t
}

func (t *GlobTarget) Make(run func(...TargetI) bool) *GlobTarget {
	t.makef = run
	return t
}

func (t *GlobTarget) modifiedSince(tm time.Time) bool {
	files, err := filepath.Glob(t.mask)
	if err != nil {
		log.Println("[buildr]", err)
		return true
	}
	if len(files) == 0 {
		return true
	}
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			log.Println("[buildr]", err)
			return true
		}
		if info.ModTime().After(tm) {
			return true
		}
	}
	return false
}

func (t *GlobTarget) name() string {
	return t.mask
}

func (t *GlobTarget) modTime() time.Time {
	tm := time.Unix(0, 0)
	files, err := filepath.Glob(t.mask)
	if err != nil {
		log.Println("[buildr]", err)
		return tm
	}
	if len(files) == 0 {
		return tm
	}
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			log.Println("[buildr]", err)
			return tm
		}
		tm_ := info.ModTime()
		if tm_.After(tm) {
			tm = tm_
		}
	}
	return tm
}

func (t *GlobTarget) Build() bool {
	tm := t.modTime()
	modified := false
	if len(t.depends) == 0 {
		modified = true
	}

	for _, targ := range t.deptab {
		if !targ.Build() {
			return false
		}
		if targ.modifiedSince(tm) {
			modified = true
		}
	}

	if modified {
		fmt.Println("[buildr] Make target `" + short(t.name()) + "`...")
		defer fmt.Println("[buildr] Done `" + short(t.name()) + "`")
		return t.makef(t.depends...)
	}
	return true
}

func (t *GlobTarget) BuildTarget(name string) bool {
	if targ, ok := t.deptab[name]; !ok {
		fmt.Println("[buildr] Cannot find target `" + short(name) + "`...")
		return false
	} else {
		return targ.Build()
	}
}
