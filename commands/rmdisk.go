package commands

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Rmdisk struct {
	Path string
}

func (tmp *Rmdisk) Execute() {
	tmp.Path = tmp.ReturnValueWithoutMarks(tmp.Path)
	var i int
	fmt.Println("Confirma la eliminacion del disco")
	fmt.Println("1. Si")
	fmt.Println("2. No")
	fmt.Scanln(&i)
	if i == 1 {
		err := os.Remove(tmp.Path)
		if err != nil {
			log.Fatal(err)
		}
	} else if i == 2 {
		fmt.Println("Eliminacion cancelada")
	} else {
		fmt.Println("Opcion Incorrecta")
	}
}
func (tmp *Rmdisk) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
