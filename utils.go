package buildr

import (
	"os"
	"log"

	"github.com/codeskyblue/go-sh"
)

func Mkdir(d string) bool {
	if err := os.Mkdir(d, os.ModePerm); err != nil {
		log.Fatalln(err)
		return false
	}
	return true
}

func Exists(f string) bool {
	_, err := os.Stat(f)
	return !os.IsNotExist(err)
}

func CreateIfNotExists(fname string) (*os.File, bool) {
	f, err := os.Create(fname)
	if os.IsExist(err) {
		f, err = os.Open(fname)
	}
	if err != nil {
		log.Fatalln(err)
		return nil, false
	}
	return f, true
}

func Cmd(cmd *sh.Session) bool {
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
		return false
	}
	return true
}

func Check(err error) bool {
	if err != nil {
		log.Fatalln(err)
		return false
	}
	return true
}
