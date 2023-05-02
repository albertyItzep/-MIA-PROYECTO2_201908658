package commands

import (
	"encoding/binary"
	"os"
	"regexp"
	"strings"
	"unsafe"

	"github.com/mia/proyecto2/structs"
)

type LoginUser struct {
	PathFile, IdPartition string
	StartPartition        int
	User, Pwd             string
	Loged                 bool
}

func (login *LoginUser) Execute() string {
	login.PathFile = login.ReturnValueWithoutMarks(login.PathFile)
	file, err := os.OpenFile(login.PathFile, os.O_RDWR, 0644)
	if err != nil {
		return "\033[31m[Error] > Al abrir el archivo " + "\033[0m"
	}
	defer file.Close()

	superBloc := structs.SuperBlock{}

	file.Seek(int64(login.StartPartition), 0)
	err = binary.Read(file, binary.LittleEndian, &superBloc)
	if err != nil {
		return "\033[31m[Error] > Al leer un superBloc en el archivo " + "\033[0m"
	}
	//realizamos la lectura de el inodo user
	inodeA := structs.InodeTable{}
	file.Seek(int64(superBloc.S_inode_start+int32(unsafe.Sizeof(structs.InodeTable{}))), 0)
	err = binary.Read(file, binary.LittleEndian, &inodeA)
	if err != nil {
		return "\033[31m[Error] > Al leer un Inode en el archivo " + "\033[0m"
	}
	//vamos a realizar la lectura de los bloques
	tmpString := ""
	for i := 0; i < 15; i++ { //recordar colocar 16 falla en mkfs
		if inodeA.I_block[i] != -1 {
			blockFile := structs.FileBlock{}
			pos := superBloc.S_block_start + inodeA.I_block[i]*int32(unsafe.Sizeof(structs.DirBlock{}))
			file.Seek(int64(pos), 0)
			err = binary.Read(file, binary.LittleEndian, &blockFile)
			if err != nil {
				return "\033[31m[Error] > Al leer un blockFile en el archivo" + "\033[0m"
			}
			tmp1 := string(blockFile.B_content[:])
			tmpString += tmp1
		}
	}
	//procesamos la informacion del archivo
	res1 := strings.Split(tmpString, "\n")
	for i := 0; i < len(res1); i++ {
		res2 := strings.Split(res1[i], ",")
		if len(res2) > 2 && res2[1] == "U" {
			if res2[3] == login.User {
				if res2[4] == login.Pwd {
					login.Loged = true
					return "ok"
				} else {
					return "pi"
				}
			}
		}
	}
	return "ui"
}

/*The function return if th user is loged*/
func (login *LoginUser) LogedUser() bool {
	return login.Loged
}

func (tmp *LoginUser) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
