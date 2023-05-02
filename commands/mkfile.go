package commands

import (
	"regexp"
	"strings"
)

type Mkfile struct {
	PathFile, IdPartition, ContFile string
	StartString                     int
	PathDisk                        string
	rFile                           bool
	SizeFile                        int
}

func (mkfile *Mkfile) Execute() {

}

func (tmp *Mkfile) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
