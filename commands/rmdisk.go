package commands

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Rmdisk struct {
	Path string
}

func (tmp *Rmdisk) Execute() string {
	tmp.Path = tmp.ReturnValueWithoutMarks(tmp.Path)
	err := os.Remove(tmp.Path)
	if err != nil {
		fmt.Println(err)
		return "Error al eliminar el disco"
	}
	return "Eliminacion Exitosa"
}
func (tmp *Rmdisk) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
