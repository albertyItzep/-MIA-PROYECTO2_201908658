package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mia/proyecto2/structs"
)

type Mkdir struct {
	PathNewDir, PathDisk     string
	CreatePrevius            bool
	StartPartition, RootNode int
}

func (mkdir *Mkdir) Execute() string {
	mkdir.PathNewDir = mkdir.ReturnValueWithoutMarks(mkdir.PathNewDir)
	mkdir.PathDisk = mkdir.ReturnValueWithoutMarks(mkdir.PathDisk)

	file, err := os.OpenFile(mkdir.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo " + "\033[0m")
		return ""
	}
	defer file.Close()

	superBlock := structs.SuperBlock{}
	//nos movemos al inicio de la particion y leemos el super Bloque
	file.Seek(int64(mkdir.StartPartition), 0)
	err = binary.Read(file, binary.LittleEndian, &superBlock)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
		return ""
	}

	//leemos el inodo inicial en la tabla de inodos
	inodePrincipal := structs.InodeTable{}
	file.Seek(int64(superBlock.S_inode_start), 0)
	err = binary.Read(file, binary.LittleEndian, &inodePrincipal)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
		return ""
	}

	arrayStrings := strings.Split(mkdir.PathNewDir, "/")
	tmp := inodePrincipal
	cadena := ""
	podInode := 0
	for i := 1; i < len(arrayStrings); i++ {
		cadena = arrayStrings[i]
		siguienteInodo, existeValor := mkdir.ReturnDirExist(file, &superBlock, tmp, cadena)
		if existeValor {
			posicion := superBlock.S_inode_start + int32(siguienteInodo)*108
			podInode = siguienteInodo
			file.Seek(int64(posicion), 0)
			err := binary.Read(file, binary.LittleEndian, &tmp)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return ""
			}
		} else {
			if i != len(arrayStrings)-1 {
				if mkdir.CreatePrevius {
					//creamos la carpeta
					fmt.Println("creamos carpeta", " "+cadena)
					inodeN := mkdir.CreateDir(file, &superBlock, &tmp, cadena, podInode)
					posicion := superBlock.S_inode_start + int32(inodeN)*108
					podInode = inodeN
					file.Seek(int64(posicion), 0)
					err := binary.Read(file, binary.LittleEndian, &tmp)
					if err != nil {
						fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
						return ""
					}
				} else {
					fmt.Println("Error No se puede crear carpeta carpeta padre no existe", " "+cadena)
					return "No Existen los directorios Padre"
				}
			} else {
				fmt.Println("creamos carpeta", " "+cadena)
				inodeN := mkdir.CreateDir(file, &superBlock, &tmp, cadena, podInode)
				posicion := superBlock.S_inode_start + int32(inodeN)*108
				podInode = inodeN
				file.Seek(int64(posicion), 0)
				err := binary.Read(file, binary.LittleEndian, &tmp)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
					return ""
				}
			}
		}
	}
	return "Path creada con exito"
}

func (mkdir *Mkdir) ReturnDirExist(file *os.File, superbloqu *structs.SuperBlock, Inode structs.InodeTable, value string) (int, bool) {
	arrayDirectorios := strings.Split(value, "/")
	for i := 0; i < len(arrayDirectorios); i++ {
		siguienteInodo, existDir := mkdir.ReturnExistValueInInode(file, superbloqu, &Inode, arrayDirectorios[i])
		if existDir {
			posicion := superbloqu.S_inode_start + int32(siguienteInodo)*108
			file.Seek(int64(posicion), 0)
			err := binary.Read(file, binary.LittleEndian, &Inode)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return -1, false
			}
			if i == len(arrayDirectorios)-1 {
				return siguienteInodo, true
			}
		}
	}
	return -1, false
}

func (mkdir *Mkdir) CreateDir(file *os.File, superbloque *structs.SuperBlock, inode *structs.InodeTable, value string, posInode int) int {
	//primero verificamos si el inodo enviado tiene espacio libre
	tmp := mkdir.RetunSpaceInBlockFree(file, superbloque, inode)
	if tmp[0] != -1 && tmp[1] != -1 {
		//solo escribimos en el bloque el nuevo valor y creamos un inodo dir
		dirBloc := structs.DirBlock{}
		pos := superbloque.S_block_start + inode.I_block[tmp[0]]*64
		file.Seek(int64(pos), 0)
		err := binary.Read(file, binary.LittleEndian, &dirBloc)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		dirBloc.B_Content[tmp[1]].B_name = mkdir.ReturnValueArr12Bytes(value)
		//debemos crear el directorio
		//primero tomamos el inodo libre en el superbloque
		inodeNew := structs.InodeTable{}
		inodeNew.I_uid = 1
		inodeNew.I_gid = 1
		inodeNew.I_size = 0
		inodeNew.I_atime = mkdir.ReturnDate8Bytes()
		inodeNew.I_ctime = mkdir.ReturnDate8Bytes()
		inodeNew.I_mtime = mkdir.ReturnDate8Bytes()
		inodeNew.I_type = 0
		inodeNew.I_perm = 664
		for i := 0; i < 16; i++ {
			inodeNew.I_block[i] = -1
		}
		dirBloc.B_Content[tmp[1]].B_inodp = superbloque.S_firts_ino
		//actualizamos el bloque
		file.Seek(int64(pos), 0)
		err = binary.Write(file, binary.LittleEndian, &dirBloc)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		//escribimos el inode
		pos = superbloque.S_inode_start + superbloque.S_firts_ino*108
		file.Seek(int64(pos), 0)
		err = binary.Write(file, binary.LittleEndian, &inodeNew)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		mkdir.RootNode = int(superbloque.S_firts_ino)
		//escribimos en el bitmap inodes
		mkdir.WriteInodeBipmapUsed(int(superbloque.S_firts_ino), file, *superbloque)

		//modificamos el superbloque
		superbloque.S_firts_ino = int32(mkdir.ReturnInodeFreeBipmap())
		superbloque.S_free_inodes_count = superbloque.S_free_inodes_count - 1
		file.Seek(int64(mkdir.StartPartition), 0)
		err = binary.Write(file, binary.LittleEndian, superbloque)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al escribir el archivo " + "\033[0m")
			return -1
		}
		fmt.Println("carpeta creada con exito ", value)
		return int(dirBloc.B_Content[tmp[1]].B_inodp)
	} else if tmp[0] != -1 && tmp[1] == -1 {
		//debemos crear un nuevo bloque y un nuevo inode
		newBlock := structs.DirBlock{}
		newInode := structs.InodeTable{}
		//insertamos la informacion en el bloque
		newBlock.B_Content[0].B_inodp = int32(mkdir.RootNode)
		newBlock.B_Content[0].B_name[0] = '.'

		newBlock.B_Content[1].B_inodp = superbloque.S_first_blo
		newBlock.B_Content[1].B_name[0] = '.'
		newBlock.B_Content[1].B_name[0] = '.'

		newBlock.B_Content[2].B_inodp = superbloque.S_firts_ino
		newBlock.B_Content[2].B_name = mkdir.ReturnValueArr12Bytes(value)

		newBlock.B_Content[3].B_name[0] = '.'

		//escribimos el bloque en fisico
		pos := superbloque.S_block_start + superbloque.S_first_blo*64
		file.Seek(int64(pos), 0)
		err := binary.Write(file, binary.LittleEndian, &newBlock)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		inode.I_block[tmp[0]] = superbloque.S_first_blo
		posicion := superbloque.S_inode_start + int32(posInode)*108
		file.Seek(int64(posicion), 0)
		err = binary.Write(file, binary.LittleEndian, inode)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		//actualizamos el nuevo inodo
		newInode.I_uid = 1
		newInode.I_gid = 1
		newInode.I_size = 0
		newInode.I_atime = mkdir.ReturnDate8Bytes()
		newInode.I_ctime = mkdir.ReturnDate8Bytes()
		newInode.I_mtime = mkdir.ReturnDate8Bytes()
		newInode.I_type = 0
		newInode.I_perm = 664
		for i := 0; i < 16; i++ {
			newInode.I_block[i] = -1
		}

		pos = superbloque.S_inode_start + superbloque.S_firts_ino*108
		file.Seek(int64(pos), 0)
		err = binary.Write(file, binary.LittleEndian, &newInode)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		mkdir.RootNode = int(superbloque.S_firts_ino)

		mkdir.WriteInodeBipmapUsed(int(superbloque.S_firts_ino), file, *superbloque)
		mkdir.WriteBlockBipmapUsed(int(superbloque.S_first_blo), file, *superbloque)

		superbloque.S_firts_ino = int32(mkdir.ReturnInodeFreeBipmap())
		superbloque.S_first_blo = int32(mkdir.ReturnBlockFreeBipmap())

		superbloque.S_free_inodes_count = superbloque.S_free_inodes_count - 1
		superbloque.S_free_blocks_count = superbloque.S_free_blocks_count - 1
		file.Seek(int64(mkdir.StartPartition), 0)
		err = binary.Write(file, binary.LittleEndian, superbloque)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al escribir el archivo " + "\033[0m")
			return -1
		}
		return int(newBlock.B_Content[2].B_inodp)
	}
	return -1
}

func (mkdir *Mkdir) ReturnExistValueInInode(file *os.File, superBlock *structs.SuperBlock, inode *structs.InodeTable, nameValue string) (int, bool) {
	if inode.I_type == 1 {
		return -1, false
	}
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			dirBloc := structs.DirBlock{}
			posicion := superBlock.S_block_start + inode.I_block[i]*64
			file.Seek(int64(posicion), 0)
			err := binary.Read(file, binary.LittleEndian, &dirBloc)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return -1, false
			}
			nextInode, existDir := mkdir.ReturnExistNameInBlock(&dirBloc, nameValue)
			if existDir {
				return nextInode, true
			}
		}
	}
	return -1, false
}

func (mkdir *Mkdir) WriteInodeBipmapUsed(byteUsed int, file *os.File, superbloque structs.SuperBlock) {
	//nos movemos a el punto inicio Bitmap + byteUsed
	pos := superbloque.S_bm_inode_start + int32(byteUsed)
	var buffer1 byte
	buffer1 = '1'
	file.Seek(int64(pos), 0)
	err := binary.Write(file, binary.LittleEndian, &buffer1)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
		return
	}
}
func (mkdir *Mkdir) WriteBlockBipmapUsed(byteUsed int, file *os.File, superbloque structs.SuperBlock) {
	//nos movemos a el punto inicio Bitmap + byteUsed
	pos := superbloque.S_bm_block_start + int32(byteUsed)
	var buffer1 byte
	buffer1 = '1'
	file.Seek(int64(pos), 0)
	err := binary.Write(file, binary.LittleEndian, &buffer1)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
		return
	}
}
func (mkdir *Mkdir) ReturnExistNameInBlock(bloc *structs.DirBlock, nameValue string) (int, bool) {
	// El bloque contiene 4 elementos content que tienen nombre y el inodo referencia
	// vamos a validar si existe el nombre
	tmp := mkdir.ReturnValueArr12BytesOFString(bloc.B_Content[2].B_name[:])
	if tmp == nameValue {
		return int(bloc.B_Content[2].B_inodp), true
	}
	tmp = mkdir.ReturnValueArr12BytesOFString(bloc.B_Content[3].B_name[:])
	if tmp == nameValue {
		return int(bloc.B_Content[3].B_inodp), true
	}
	return -1, false
}

func (mkdir *Mkdir) RetunSpaceInBlockFree(file *os.File, superBloque *structs.SuperBlock, inode *structs.InodeTable) [2]int {
	var tmp [2]int
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			dirBloc := structs.DirBlock{}
			posicion := superBloque.S_block_start + inode.I_block[i]*64
			file.Seek(int64(posicion), 0)
			err := binary.Read(file, binary.LittleEndian, &dirBloc)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return tmp
			}
			nameDirBlock := mkdir.ReturnValueArr12BytesOFString(dirBloc.B_Content[2].B_name[:])
			if nameDirBlock == "" {
				tmp[0] = i
				tmp[1] = 2
				return tmp
			}
			nameDirBlock = mkdir.ReturnValueArr12BytesOFString(dirBloc.B_Content[3].B_name[:])
			if nameDirBlock == "" || nameDirBlock == "." {
				tmp[0] = i
				tmp[1] = 3
				return tmp
			}
		}
	}
	for i := 0; i < 16; i++ {
		if inode.I_block[i] == -1 {
			tmp[0] = i
			tmp[1] = -1
			return tmp
		}
	}
	tmp[0] = -1
	tmp[1] = -1
	return tmp
}

func (mkdir *Mkdir) ReturnValueArr12Bytes(value string) [12]byte {
	var tmp [12]byte
	for i := 0; i < 12; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}

func (mkdir *Mkdir) ReturnValueArr12BytesOFString(value []byte) string {
	var tmp string
	var charTmp byte
	for i := 0; i < len(value); i++ {
		if value[i] != charTmp {
			tmp += string(value[i])
		} else {
			break
		}
	}
	return tmp
}

func (Mkdir *Mkdir) ReturnInodeFreeBipmap() int {
	//abrimos el archivo
	file, err := os.OpenFile(Mkdir.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo " + "\033[0m")
		return -1
	}
	defer file.Close()
	superBlock := structs.SuperBlock{}
	//nos movemos al inicio de la particion
	file.Seek(int64(Mkdir.StartPartition), 0)
	err = binary.Read(file, binary.LittleEndian, &superBlock)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
		return -1
	}

	var bitInode byte
	for i := 0; i < int(superBlock.S_inodes_count); i++ {
		file.Seek(int64(int(superBlock.S_bm_inode_start)+i), 0)
		err = binary.Read(file, binary.LittleEndian, &bitInode)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		if bitInode == '0' {
			return i
		}
	}
	return -1
}

func (Mkdir *Mkdir) ReturnBlockFreeBipmap() int {
	//abrimos el archivo
	file, err := os.OpenFile(Mkdir.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo " + "\033[0m")
		return -1
	}
	defer file.Close()
	superBlock := structs.SuperBlock{}
	//nos movemos al inicio de la particion
	file.Seek(int64(Mkdir.StartPartition), 0)
	err = binary.Read(file, binary.LittleEndian, &superBlock)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
		return -1
	}

	var bitInode byte
	for i := 0; i < int(superBlock.S_blocks_count); i++ {
		file.Seek(int64(int(superBlock.S_bm_block_start)+i), 0)
		err = binary.Read(file, binary.LittleEndian, &bitInode)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		if bitInode == '0' {
			return i
		}
	}
	return -1
}

func (tmp *Mkdir) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
func (mkdir *Mkdir) ReturnDate8Bytes() [8]byte {
	t := string(time.Now().Format("02012006"))
	tmpT := []byte(t)
	return [8]byte(tmpT)
}
