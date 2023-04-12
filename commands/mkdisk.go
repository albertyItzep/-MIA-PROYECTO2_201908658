package commands

import (
	"encoding/binary"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mia/proyecto2/structs"
)

type Mkdisk struct {
	Path      string
	Fit, Unit byte
	Size      int
}

func (tmp *Mkdisk) Execute() {
	tmp.Path = tmp.ReturnValueWithoutMarks(tmp.Path)
	tamDiskK := tmp.ReturnSizeRep(tmp.Size, tmp.Unit)
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
	if err2 != nil && !os.IsExist(err) {
		log.Fatal(err2)
	}
	defer file.Close()

	err = os.Chmod(tmp.Path, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// escribimos el archivo en tamano bytes
	var c [1024]byte
	for i := 0; i < tamDiskK; i++ {
		err2 = binary.Write(file, binary.LittleEndian, &c)
		if err2 != nil && !os.IsExist(err) {
			log.Fatal(err2)
		}
	}
	// escribimos el mbr dentro de el archivo
	t := string(time.Now().Format("02012006"))
	tmpT := []byte(t)
	tmpS := uint32(rand.Intn(101))
	mbr := structs.MBR{Mbr_tamano: uint32(tmp.Size), Mbr_fecha_creacion: [8]byte(tmpT), Mbr_dsk_signature: tmpS}

	file.Seek(0, 0)
	err2 = binary.Write(file, binary.LittleEndian, &mbr)
	if err2 != nil && !os.IsExist(err) {
		log.Fatal(err2)
	}
}

/*Return the size in bytes*/
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

/*Return the cant of reps for the write*/
func (tmp *Mkdisk) ReturnSizeRep(sizeI int, unit byte) int {
	var size int
	switch {
	case unit == 'k':
		size = sizeI
	case unit == 'm' || tmp.Unit == 'o':
		size = sizeI * 1024
	}
	return size
}

/*Return value without marks*/
func (tmp *Mkdisk) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}

/*Return path without FileName for the Mkdisk*/
func (tmp *Mkdisk) ReturnPathWithoutFileName(value string) string {
	var tmpString string
	regPahtMkdir := regexp.MustCompile("/[a-zA-Z0-9]+.dsk")
	tmpString = regPahtMkdir.ReplaceAllString(value, "")
	return tmpString
}
