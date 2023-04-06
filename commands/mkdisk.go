package commands

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"encoding/binary"
	"strings"
)

type Mkdisk struct {
	Path      string
	Fit, Unit byte
	Size      int
}

func (tmp *Mkdisk) Execute() {
	tmp.Path = tmp.ReturnValueWithoutMarks(tmp.Path)
	tmp.Size = tmp.ReturnSize(tmp.Size, tmp.Unit)
	//obtenemos la ruta sin el nombre del archivo
	pathTmp := tmp.ReturnPathWithoutFileName(tmp.Path)

	err := os.MkdirAll(pathTmp, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	err = os.Chmod(pathTmp, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	file, err2 := os.Create(tmp.Path)
	defer file.Close()
	if err2 != nil && !os.IsExist(err) {
		log.Fatal(err2)
	}
	err = os.Chmod(tmp.Path, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	// escribimos el archivo en tamano bytes
	var c [1024]byte
	// abrimos el archivo y cargamos informacion
	for i := 0; i < tmp.Size; i++ {
		var bin\_buf bytes.Buffer
		binary.Write(&bin\_buf,binary.BigEndian,c)
	}
}
func (tmp *Mkdisk) ReturnSize(sizeI int, unit byte) int {
	var size int
	switch {
	case unit == 'k':
		size = sizeI * 1024
	case unit == 'm' || tmp.Unit == 'o':
		size = sizeI * 1024 * 1024
	}
	return size
}
func (tmp *Mkdisk) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
func (tmp *Mkdisk) ReturnPathWithoutFileName(value string) string {
	var tmpString string
	regPahtMkdir := regexp.MustCompile("/[a-zA-Z0-9]+.dsk")
	tmpString = regPahtMkdir.ReplaceAllString(value, "")
	return tmpString
}
