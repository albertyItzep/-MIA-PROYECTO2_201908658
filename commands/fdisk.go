package commands

import (
	"regexp"
	"strings"

	"github.com/mia/proyecto2/structs"
)

type Fdisk struct {
	Name, Path      string
	Fit, Type, Unit byte
	Size            uint32
	MbrFdisk        structs.MBR
}

func (tmp *Fdisk) Execute() {
	tmp.Name = tmp.ReturnValueWithoutMarks(tmp.Name)
	tmp.Path = tmp.ReturnValueWithoutMarks(tmp.Path)

}
func (tmp *Fdisk) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
