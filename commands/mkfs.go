package commands

import (
	"encoding/binary"
	"fmt"
	"github.com/mia/proyecto2/structs"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

type Mkfs struct {
	IdMkfs, TypeMkfs string
	SizeOfPartition  int
	StartPartition   int
	PathFile         string
}

func (mkfs *Mkfs) Execute() {
	mkfs.PathFile = mkfs.ReturnValueWithoutMarks(mkfs.PathFile)

	superBlock := structs.SuperBlock{}
	block := structs.FileBlock{}
	inode := structs.InodeTable{}

	//calculamos el numero total de inodos para la particion
	var n, div float64
	n = float64(mkfs.SizeOfPartition - int(unsafe.Sizeof(superBlock)))
	div = 4 + float64(unsafe.Sizeof(inode)) + 3*float64(unsafe.Sizeof(block))
	n = n / div
	n = math.Floor(n)

	//insertamos valores iniciales al superbloque
	superBlock.S_filesystem_type = 1
	superBlock.S_inodes_count = int32(n)
	superBlock.S_blocks_count = int32(3 * n)
	superBlock.S_free_inodes_count = int32(n) - 2
	superBlock.S_free_blocks_count = int32(3*n) - 2
	superBlock.S_mtime = mkfs.ReturnDate8Bytes()
	superBlock.S_mnt_count = 1
	superBlock.S_magic = 0xEF53
	superBlock.S_inode_size = int32(unsafe.Sizeof(structs.InodeTable{}))
	superBlock.S_block_size = int32(unsafe.Sizeof(structs.FileBlock{}))

	startBitmapInodes := mkfs.StartPartition + int(unsafe.Sizeof(structs.SuperBlock{}))
	startBitmapBloks := startBitmapInodes + int(n)
	superBlock.S_bm_inode_start = int32(startBitmapInodes)
	superBlock.S_bm_block_start = int32(startBitmapBloks)

	firstInodeFree := startBitmapBloks + int(3*n)
	firstBlockFree := firstInodeFree + int(n)*int(unsafe.Sizeof(structs.SuperBlock{}))
	superBlock.S_firts_ino = int32(2)
	superBlock.S_firts_ino = int32(2)

	superBlock.S_inode_start = int32(firstInodeFree)
	superBlock.S_block_start = int32(firstBlockFree)
	//abrimos el archivo
	file, err := os.OpenFile(mkfs.PathFile, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo:", err, "\033[0m")
		return
	}
	defer file.Close()

	//escribimos el superBlock
	file.Seek(int64(mkfs.StartPartition), 0)
	err = binary.Write(file, binary.LittleEndian, &superBlock)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir el superbloque:", err, "\033[0m")
		return
	}

	//escribimos el bitmapInodes
	buffer := '0'
	buffer2 := '1'
	for i := 0; i < int(n); i++ {
		file.Seek(int64(int(superBlock.S_bm_inode_start)+i), 0)
		err = binary.Write(file, binary.LittleEndian, &buffer)
		if err != nil {
			fmt.Println("\033[31m[Error] > Error al escribir el bitmapInodes:", err, "\033[0m")
			return
		}
	}
	//insertamos los inodos usados para user.txt
	file.Seek(int64(superBlock.S_bm_inode_start), 0)
	err = binary.Write(file, binary.LittleEndian, &buffer2)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir un inodo en bitmapInodes:", err, "\033[0m")
		return
	}
	file.Seek(int64(superBlock.S_bm_inode_start+1), 0)
	err = binary.Write(file, binary.LittleEndian, &buffer2)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir un inodo en bitmapInodes:", err, "\033[0m")
		return
	}

	//escribimos el bitmapBloques
	for i := 0; i < int(3*n); i++ {
		file.Seek(int64(int(superBlock.S_bm_block_start)+i), 0)
		err = binary.Write(file, binary.LittleEndian, &buffer)
		if err != nil {
			fmt.Println("\033[31m[Error] > Error al escribir el bitmapBloques:", err, "\033[0m")
			return
		}
	}
	//insertamos los bloques usados por user.txt
	file.Seek(int64(superBlock.S_bm_block_start), 0)
	err = binary.Write(file, binary.LittleEndian, &buffer2)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir  un inodo en bitmapBlocks:", err, "\033[0m")
		return
	}
	file.Seek(int64(superBlock.S_bm_block_start+1), 0)
	err = binary.Write(file, binary.LittleEndian, &buffer2)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir un inodo en bitmapBlocks:", err, "\033[0m")
		return
	}
	// inodo para la carperta root
	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_size = 27
	inode.I_atime = mkfs.ReturnDate8Bytes()
	inode.I_ctime = mkfs.ReturnDate8Bytes()
	inode.I_mtime = mkfs.ReturnDate8Bytes()
	inode.I_block[0] = 0
	inode.I_type = '0'
	inode.I_perm = 664
	for i := 1; i < 16; i++ {
		inode.I_block[i] = -1
	}
	file.Seek(int64(superBlock.S_inode_start), 0)
	err = binary.Write(file, binary.LittleEndian, &inode)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir un inodo:", err, "\033[0m")
		return
	}
	// bloque para la carpeta root
	blocUser := structs.DirBlock{}
	blocUser.B_Content[0].B_inodp = 0
	blocUser.B_Content[0].B_name[0] = '.'

	blocUser.B_Content[1].B_inodp = 0
	blocUser.B_Content[1].B_name[0] = '.'
	blocUser.B_Content[1].B_name[1] = '.'

	blocUser.B_Content[2].B_inodp = 1
	blocUser.B_Content[2].B_name = mkfs.ReturnValueArr12Bytes("users.txt")

	blocUser.B_Content[3].B_inodp = -1
	blocUser.B_Content[3].B_name[0] = '.'

	file.Seek(int64(superBlock.S_block_start), 0)
	err = binary.Write(file, binary.LittleEndian, &blocUser)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir un bloque:", err, "\033[0m")
		return
	}
	//inodo de tipo archivo para users.txt
	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_size = 27
	inode.I_atime = mkfs.ReturnDate8Bytes()
	inode.I_ctime = mkfs.ReturnDate8Bytes()
	inode.I_mtime = mkfs.ReturnDate8Bytes()
	inode.I_block[0] = 1
	inode.I_type = '1'
	inode.I_perm = 755
	for i := 1; i < 16; i++ {
		inode.I_block[i] = -1
	}
	file.Seek(int64(superBlock.S_inode_start+int32(unsafe.Sizeof(structs.InodeTable{}))), 0)
	err = binary.Write(file, binary.LittleEndian, &inode)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir un inodo:", err, "\033[0m")
		return
	}

	//bloque archivo para users.txt
	blocUserT := structs.FileBlock{}
	blocUserT.B_content = mkfs.ReturnValueArr64Bytes("1,G,root\n1,U,root,root,123\n")
	file.Seek(int64(superBlock.S_block_start+int32(unsafe.Sizeof(structs.FileBlock{}))), 0)
	err = binary.Write(file, binary.LittleEndian, &blocUserT)
	if err != nil {
		fmt.Println("\033[31m[Error] > Error al escribir un inodo:", err, "\033[0m")
		return
	}
	fmt.Println("Formating Ext2")
	fmt.Println("...")
	fmt.Println("Formato exitoso")
}
func (tmp *Mkfs) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}

/*The function return the date of mounted the sistem*/
func (mkfs *Mkfs) ReturnDate8Bytes() [8]byte {
	t := string(time.Now().Format("02012006"))
	tmpT := []byte(t)
	return [8]byte(tmpT)
}

/*The function return the value in an array of 12 bytes*/
func (mkfs *Mkfs) ReturnValueArr12Bytes(value string) [12]byte {
	var tmp [12]byte
	for i := 0; i < 12; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}

/*The function return the value in an array of 12 bytes*/
func (mkfs *Mkfs) ReturnValueArr64Bytes(value string) [64]byte {
	var tmp [64]byte
	for i := 0; i < 64; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}
