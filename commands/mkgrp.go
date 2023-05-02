package commands

import (
	"encoding/binary"
	"fmt"
	"github.com/mia/proyecto2/structs"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Mkgrp struct {
	NameGrup, PathFile string
	StartParition      int
}

func (mkgrp *Mkgrp) Execute() {
	mkgrp.NameGrup = mkgrp.ReturnValueWithoutMarks(mkgrp.NameGrup)
	mkgrp.PathFile = mkgrp.ReturnValueWithoutMarks(mkgrp.PathFile)

	file, err := os.OpenFile(mkgrp.PathFile, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo:", err, "\033[0m")
		return
	}
	defer file.Close()

	superBloc := structs.SuperBlock{}

	file.Seek(int64(mkgrp.StartParition), 0)
	err = binary.Read(file, binary.LittleEndian, &superBloc)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer un superBloc en el archivo:", err, "\033[0m")
		return
	}

	//realizamos la lectura de el inodo user
	inodeA := structs.InodeTable{}
	file.Seek(int64(superBloc.S_inode_start+int32(unsafe.Sizeof(structs.InodeTable{}))), 0)
	err = binary.Read(file, binary.LittleEndian, &inodeA)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer un Inode en el archivo:", err, "\033[0m")
		return
	}
	//vamos a realizar la lectura de los bloques
	tmpString := ""
	blockFile := structs.FileBlock{}
	for i := 0; i < 15; i++ { //recordar colocar 16 falla en mkfs
		if inodeA.I_block[i] != -1 {
			pos := superBloc.S_block_start + inodeA.I_block[i]*int32(unsafe.Sizeof(structs.DirBlock{}))
			file.Seek(int64(pos), 0)
			err = binary.Read(file, binary.LittleEndian, &blockFile)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer un blockFile en el archivo:", err, "\033[0m")
				return
			}
			tmp1 := string(blockFile.B_content[:])
			tmpString += tmp1
		}
	}
	//procesamos la informacion del archivo
	res1 := strings.Split(tmpString, "\n")
	fmt.Println(res1)
	finalString := ""
	numGrupExist := 0
	for i := 0; i < len(res1); i++ {
		splitWithComma := strings.Split(res1[i], ",")
		if len(splitWithComma) > 2 && splitWithComma[1] == "G" {
			numGrup, err := strconv.Atoi(splitWithComma[0])
			if err != nil {
				log.Fatal(err)
			}
			numGrupExist = numGrup
			finalString += res1[i] + "\n"
		} else if len(splitWithComma) > 2 {
			finalString += res1[i] + "\n"
		}
	}
	numGrupExist++
	stringIsert := strconv.Itoa(numGrupExist) + ",G," + mkgrp.NameGrup
	finalString += stringIsert
	if len(finalString) > 64 {
		CantM := len(finalString) / 64
		numBlocks := math.Floor(float64(CantM))
		for int(numBlocks*64) < len(finalString) {
			numBlocks++
		}
		//verificamos el bipmap de bloques
		//para ver los bloques libres
		arrBlocksFree := make([]int, int64(numBlocks))
		file.Seek(int64(superBloc.S_bm_block_start+64), 0)
		i, z := 0, 0
		for i < int(numBlocks) {
			var statusByte byte
			err = binary.Read(file, binary.LittleEndian, &statusByte)
			if err != nil {
				log.Fatal(err)
			}
			if statusByte == '0' {
				i++
				arrBlocksFree = append(arrBlocksFree, z)
			}
			z++
			if (z + int(superBloc.S_bm_block_start)) >= int(superBloc.S_block_start) {
				break
			}
		}
		//reescribimos el inodo archivo
		inodeA.I_size = int32(len(finalString))
		inodeA.I_mtime = mkgrp.ReturnDate8Bytes()

	} else {
		blockFile.B_content = mkgrp.ReturnValueArr64Bytes(finalString)
		file.Seek(int64(superBloc.S_block_start), 0)
		err = binary.Write(file, binary.LittleEndian, &blockFile)
		if err != nil {
			log.Println(err)
		}
		inodeA.I_size = int32(len(finalString))
		inodeA.I_mtime = mkgrp.ReturnDate8Bytes()
		file.Seek(int64(superBloc.S_inode_start+int32(unsafe.Sizeof(structs.InodeTable{}))), 0)
		err = binary.Write(file, binary.LittleEndian, &inodeA)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println("Grupo creado")
}
func (mkfs *Mkgrp) ReturnValueArr64Bytes(value string) [64]byte {
	var tmp [64]byte
	for i := 0; i < 64; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}

func (mkfs *Mkgrp) ReturnDate8Bytes() [8]byte {
	t := string(time.Now().Format("02012006"))
	tmpT := []byte(t)
	return [8]byte(tmpT)
}

func (tmp *Mkgrp) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
