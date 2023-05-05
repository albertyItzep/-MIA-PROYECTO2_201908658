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

type Mkfile struct {
	PathNewFile, PathDisk, PathContent string
	StartPartition, SizeNewFile        int
	CreatePrevius                      bool
	RootNode                           int
}

func (mkfile *Mkfile) Execute() string {
	mkfile.PathNewFile = mkfile.ReturnValueWithoutMarks(mkfile.PathNewFile)
	mkfile.PathDisk = mkfile.ReturnValueWithoutMarks(mkfile.PathDisk)

	file, err := os.OpenFile(mkfile.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo " + "\033[0m")
		return ""
	}
	defer file.Close()

	superBlock := structs.SuperBlock{}
	//nos movemos al inicio de la particion y leemos el super Bloque
	file.Seek(int64(mkfile.StartPartition), 0)
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

	arrayStrings := strings.Split(mkfile.PathNewFile, "/")
	tmp := inodePrincipal
	cadena := ""
	podInode := 0
	for i := 1; i < len(arrayStrings); i++ {
		cadena = arrayStrings[i]
		siguienteInodo, existeValor := mkfile.ReturnDirExist(file, &superBlock, tmp, cadena)
		if existeValor {
			if i != len(arrayStrings)-1 {
				posicion := superBlock.S_inode_start + int32(siguienteInodo)*108
				podInode = siguienteInodo
				file.Seek(int64(posicion), 0)
				err := binary.Read(file, binary.LittleEndian, &tmp)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
					return ""
				}
			} else {
				return "El archivo ya existe"
			}
		} else {
			if i != len(arrayStrings)-1 {
				if mkfile.CreatePrevius {
					//creamos la carpeta
					fmt.Println("creamos carpeta", " "+cadena)
					inodeN := mkfile.CreateDir(file, &superBlock, &tmp, cadena, podInode)
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
				fmt.Println("creamos el archivo que deseamos", " "+cadena)
				mkfile.CreateFile(file, &superBlock, &tmp, cadena, podInode)
			}
		}
	}
	return "Path creada con exito"
}

func (mkfile *Mkfile) CreateFile(file *os.File, superbloque *structs.SuperBlock, inode *structs.InodeTable, value string, posInode int) {
	//primero validamos la cantidad de bloques que vamos a utilizar
	resultado := 0
	if mkfile.SizeNewFile > 0 {
		resultado = mkfile.SizeNewFile / 64
		if mkfile.SizeNewFile%64 != 0 {
			resultado++
		}
	}
	tmp := mkfile.RetunSpaceInBlockFree(file, superbloque, inode)
	if tmp[0] != -1 && tmp[1] != -1 {
		dirBloc := structs.DirBlock{}
		pos := superbloque.S_block_start + inode.I_block[tmp[0]]*64
		file.Seek(int64(pos), 0)
		err := binary.Read(file, binary.LittleEndian, &dirBloc)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return
		}
		dirBloc.B_Content[tmp[1]].B_name = mkfile.ReturnValueArr12Bytes(value)
		inodeNew := structs.InodeTable{}
		inodeNew.I_uid = 1
		inodeNew.I_gid = 1
		inodeNew.I_size = 0
		inodeNew.I_atime = mkfile.ReturnDate8Bytes()
		inodeNew.I_ctime = mkfile.ReturnDate8Bytes()
		inodeNew.I_mtime = mkfile.ReturnDate8Bytes()
		inodeNew.I_type = 1
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
			return
		}
		//escribimos el inode
		pos = superbloque.S_inode_start + superbloque.S_firts_ino*108
		inodePos := pos
		file.Seek(int64(pos), 0)
		err = binary.Write(file, binary.LittleEndian, &inodeNew)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return
		}
		mkfile.RootNode = int(superbloque.S_firts_ino)
		//escribimos en el bitmap inodes
		mkfile.WriteInodeBipmapUsed(int(superbloque.S_firts_ino), file, *superbloque)

		//modificamos el superbloque
		superbloque.S_firts_ino = int32(mkfile.ReturnInodeFreeBipmap())
		superbloque.S_free_inodes_count = superbloque.S_free_inodes_count - 1
		file.Seek(int64(mkfile.StartPartition), 0)
		err = binary.Write(file, binary.LittleEndian, superbloque)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al escribir el archivo " + "\033[0m")
			return
		}

		//creamos los bloques con contenido
		//debemos crear n cantidad de bloques y los llenamos de contenido
		if mkfile.SizeNewFile > 0 {
			araCharacters := [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
			cont := 0
			file.Seek(int64(inodePos), 0)
			inodeTmp := structs.InodeTable{}
			err = binary.Read(file, binary.LittleEndian, &inodeTmp)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return
			}
			for i := 0; i < resultado; i++ {
				newBlockFile := structs.FileBlock{}
				for x := 0; x < 64; x++ {
					if cont == 9 {
						cont = 0
					}
					if mkfile.SizeNewFile < 64 && x == mkfile.SizeNewFile {
						break
					}
					newBlockFile.B_content[x] = araCharacters[cont]
					cont++
				}
				cont = 0
				posicionInsercion := mkfile.ReturnBlockFreeInInode(inodeTmp)
				inodeTmp.I_block[posicionInsercion] = superbloque.S_first_blo
				pos := superbloque.S_block_start + superbloque.S_first_blo
				file.Seek(int64(pos), 0)
				err = binary.Write(file, binary.LittleEndian, &newBlockFile)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
					return
				}
				//Toca actualizar el superbloque
				mkfile.WriteBlockBipmapUsed(int(superbloque.S_first_blo), file, *superbloque)
				superbloque.S_free_blocks_count = superbloque.S_free_blocks_count - 1

				superbloque.S_first_blo = int32(mkfile.ReturnBlockFreeBipmap())
				file.Seek(int64(mkfile.StartPartition), 0)
				err = binary.Write(file, binary.LittleEndian, superbloque)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al escribir el archivo " + "\033[0m")
					return
				}
			}
			file.Seek(int64(inodePos), 0)
			err = binary.Write(file, binary.LittleEndian, &inodeTmp)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return
			}
		}
	} else if tmp[0] != -1 && tmp[1] == -1 {
		newBlock := structs.DirBlock{}
		newInode := structs.InodeTable{}
		//insertamos la informacion en el bloque
		newBlock.B_Content[0].B_inodp = int32(mkfile.RootNode)
		newBlock.B_Content[0].B_name[0] = '.'

		newBlock.B_Content[1].B_inodp = superbloque.S_first_blo
		newBlock.B_Content[1].B_name[0] = '.'
		newBlock.B_Content[1].B_name[0] = '.'

		newBlock.B_Content[2].B_inodp = superbloque.S_firts_ino
		newBlock.B_Content[2].B_name = mkfile.ReturnValueArr12Bytes(value)

		newBlock.B_Content[3].B_name[0] = '.'

		//escribimos el bloque en fisico
		pos := superbloque.S_block_start + superbloque.S_first_blo*64
		file.Seek(int64(pos), 0)
		err := binary.Write(file, binary.LittleEndian, &newBlock)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return
		}
		//actualizamos el bloque en donde nos encontramos para que apunte al nuevo bloque
		inode.I_block[tmp[0]] = superbloque.S_first_blo
		posicion := superbloque.S_inode_start + int32(posInode)*108
		file.Seek(int64(posicion), 0)
		err = binary.Write(file, binary.LittleEndian, inode)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return
		}
		//actualizamos el nuevo inodo
		newInode.I_uid = 1
		newInode.I_gid = 1
		newInode.I_size = 0
		newInode.I_atime = mkfile.ReturnDate8Bytes()
		newInode.I_ctime = mkfile.ReturnDate8Bytes()
		newInode.I_mtime = mkfile.ReturnDate8Bytes()
		newInode.I_type = 1
		newInode.I_perm = 664
		for i := 0; i < 16; i++ {
			newInode.I_block[i] = -1
		}

		pos = superbloque.S_inode_start + superbloque.S_firts_ino*108
		inodePos := pos
		file.Seek(int64(pos), 0)
		err = binary.Write(file, binary.LittleEndian, &newInode)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return
		}
		mkfile.RootNode = int(superbloque.S_firts_ino)

		mkfile.WriteInodeBipmapUsed(int(superbloque.S_firts_ino), file, *superbloque)
		mkfile.WriteBlockBipmapUsed(int(superbloque.S_first_blo), file, *superbloque)

		superbloque.S_firts_ino = int32(mkfile.ReturnInodeFreeBipmap())
		superbloque.S_first_blo = int32(mkfile.ReturnBlockFreeBipmap())

		superbloque.S_free_inodes_count = superbloque.S_free_inodes_count - 1
		superbloque.S_free_blocks_count = superbloque.S_free_blocks_count - 1
		file.Seek(int64(mkfile.StartPartition), 0)
		err = binary.Write(file, binary.LittleEndian, superbloque)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al escribir el archivo " + "\033[0m")
			return
		}
		if mkfile.SizeNewFile > 0 {
			araCharacters := [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
			cont := 0
			file.Seek(int64(inodePos), 0)
			inodeTmp := structs.InodeTable{}
			err = binary.Read(file, binary.LittleEndian, &inodeTmp)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return
			}
			for i := 0; i < resultado; i++ {
				newBlockFile := structs.FileBlock{}
				for x := 0; x < 64; x++ {
					if cont == 9 {
						cont = 0
					}
					if mkfile.SizeNewFile < 64 && x == mkfile.SizeNewFile {
						break
					}
					newBlockFile.B_content[x] = araCharacters[cont]
					cont++
				}
				cont = 0
				posicionInsercion := mkfile.ReturnBlockFreeInInode(inodeTmp)
				inodeTmp.I_block[posicionInsercion] = superbloque.S_first_blo
				pos := superbloque.S_block_start + superbloque.S_first_blo
				file.Seek(int64(pos), 0)
				err = binary.Write(file, binary.LittleEndian, &newBlockFile)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
					return
				}
				//Toca actualizar el superbloque
				mkfile.WriteBlockBipmapUsed(int(superbloque.S_first_blo), file, *superbloque)
				superbloque.S_free_blocks_count = superbloque.S_free_blocks_count - 1

				superbloque.S_first_blo = int32(mkfile.ReturnBlockFreeBipmap())
				file.Seek(int64(mkfile.StartPartition), 0)
				err = binary.Write(file, binary.LittleEndian, superbloque)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al escribir el archivo " + "\033[0m")
					return
				}
			}
			file.Seek(int64(inodePos), 0)
			err = binary.Write(file, binary.LittleEndian, &inodeTmp)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
				return
			}
		}
	}
}
func (mkfile *Mkfile) ReturnBlockFreeInInode(inode structs.InodeTable) int {
	for i := 0; i < 16; i++ {
		if inode.I_block[i] == -1 {
			return i
		}
	}
	return -1
}
func (mkfile *Mkfile) ReturnDirExist(file *os.File, superbloqu *structs.SuperBlock, Inode structs.InodeTable, value string) (int, bool) {
	arrayDirectorios := strings.Split(value, "/")
	for i := 0; i < len(arrayDirectorios); i++ {
		siguienteInodo, existDir := mkfile.ReturnExistValueInInode(file, superbloqu, &Inode, arrayDirectorios[i])
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

func (mkfile *Mkfile) CreateDir(file *os.File, superbloque *structs.SuperBlock, inode *structs.InodeTable, value string, posInode int) int {
	//primero verificamos si el inodo enviado tiene espacio libre
	tmp := mkfile.RetunSpaceInBlockFree(file, superbloque, inode)
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
		dirBloc.B_Content[tmp[1]].B_name = mkfile.ReturnValueArr12Bytes(value)
		//debemos crear el directorio
		//primero tomamos el inodo libre en el superbloque
		inodeNew := structs.InodeTable{}
		inodeNew.I_uid = 1
		inodeNew.I_gid = 1
		inodeNew.I_size = 0
		inodeNew.I_atime = mkfile.ReturnDate8Bytes()
		inodeNew.I_ctime = mkfile.ReturnDate8Bytes()
		inodeNew.I_mtime = mkfile.ReturnDate8Bytes()
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
		mkfile.RootNode = int(superbloque.S_firts_ino)
		//escribimos en el bitmap inodes
		mkfile.WriteInodeBipmapUsed(int(superbloque.S_firts_ino), file, *superbloque)

		//modificamos el superbloque
		superbloque.S_firts_ino = int32(mkfile.ReturnInodeFreeBipmap())
		superbloque.S_free_inodes_count = superbloque.S_free_inodes_count - 1
		file.Seek(int64(mkfile.StartPartition), 0)
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
		newBlock.B_Content[0].B_inodp = int32(mkfile.RootNode)
		newBlock.B_Content[0].B_name[0] = '.'

		newBlock.B_Content[1].B_inodp = superbloque.S_first_blo
		newBlock.B_Content[1].B_name[0] = '.'
		newBlock.B_Content[1].B_name[0] = '.'

		newBlock.B_Content[2].B_inodp = superbloque.S_firts_ino
		newBlock.B_Content[2].B_name = mkfile.ReturnValueArr12Bytes(value)

		newBlock.B_Content[3].B_name[0] = '.'

		//escribimos el bloque en fisico
		pos := superbloque.S_block_start + superbloque.S_first_blo*64
		file.Seek(int64(pos), 0)
		err := binary.Write(file, binary.LittleEndian, &newBlock)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al leer el archivo " + "\033[0m")
			return -1
		}
		//actualizamos el bloque en donde nos encontramos para que apunte al nuevo bloque
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
		newInode.I_atime = mkfile.ReturnDate8Bytes()
		newInode.I_ctime = mkfile.ReturnDate8Bytes()
		newInode.I_mtime = mkfile.ReturnDate8Bytes()
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
		mkfile.RootNode = int(superbloque.S_firts_ino)

		mkfile.WriteInodeBipmapUsed(int(superbloque.S_firts_ino), file, *superbloque)
		mkfile.WriteBlockBipmapUsed(int(superbloque.S_first_blo), file, *superbloque)

		superbloque.S_firts_ino = int32(mkfile.ReturnInodeFreeBipmap())
		superbloque.S_first_blo = int32(mkfile.ReturnBlockFreeBipmap())

		superbloque.S_free_inodes_count = superbloque.S_free_inodes_count - 1
		superbloque.S_free_blocks_count = superbloque.S_free_blocks_count - 1
		file.Seek(int64(mkfile.StartPartition), 0)
		err = binary.Write(file, binary.LittleEndian, superbloque)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al escribir el archivo " + "\033[0m")
			return -1
		}
		return int(newBlock.B_Content[2].B_inodp)
	}
	return -1
}

func (mkfile *Mkfile) ReturnExistValueInInode(file *os.File, superBlock *structs.SuperBlock, inode *structs.InodeTable, nameValue string) (int, bool) {
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
			nextInode, existDir := mkfile.ReturnExistNameInBlock(&dirBloc, nameValue)
			if existDir {
				return nextInode, true
			}
		}
	}
	return -1, false
}

func (mkfile *Mkfile) WriteInodeBipmapUsed(byteUsed int, file *os.File, superbloque structs.SuperBlock) {
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
func (mkfile *Mkfile) WriteBlockBipmapUsed(byteUsed int, file *os.File, superbloque structs.SuperBlock) {
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
func (mkfile *Mkfile) ReturnExistNameInBlock(bloc *structs.DirBlock, nameValue string) (int, bool) {
	// El bloque contiene 4 elementos content que tienen nombre y el inodo referencia
	// vamos a validar si existe el nombre
	tmp := mkfile.ReturnValueArr12BytesOFString(bloc.B_Content[2].B_name[:])
	if tmp == nameValue {
		return int(bloc.B_Content[2].B_inodp), true
	}
	tmp = mkfile.ReturnValueArr12BytesOFString(bloc.B_Content[3].B_name[:])
	if tmp == nameValue {
		return int(bloc.B_Content[3].B_inodp), true
	}
	return -1, false
}

func (mkfile *Mkfile) RetunSpaceInBlockFree(file *os.File, superBloque *structs.SuperBlock, inode *structs.InodeTable) [2]int {
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
			nameDirBlock := mkfile.ReturnValueArr12BytesOFString(dirBloc.B_Content[2].B_name[:])
			if nameDirBlock == "" {
				tmp[0] = i
				tmp[1] = 2
				return tmp
			}
			nameDirBlock = mkfile.ReturnValueArr12BytesOFString(dirBloc.B_Content[3].B_name[:])
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

func (mkfile *Mkfile) ReturnValueArr12Bytes(value string) [12]byte {
	var tmp [12]byte
	for i := 0; i < 12; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}

func (mkfile *Mkfile) ReturnValueArr12BytesOFString(value []byte) string {
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

func (mkfile *Mkfile) ReturnInodeFreeBipmap() int {
	//abrimos el archivo
	file, err := os.OpenFile(mkfile.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo " + "\033[0m")
		return -1
	}
	defer file.Close()
	superBlock := structs.SuperBlock{}
	//nos movemos al inicio de la particion
	file.Seek(int64(mkfile.StartPartition), 0)
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

func (mkfile *Mkfile) ReturnBlockFreeBipmap() int {
	//abrimos el archivo
	file, err := os.OpenFile(mkfile.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo " + "\033[0m")
		return -1
	}
	defer file.Close()
	superBlock := structs.SuperBlock{}
	//nos movemos al inicio de la particion
	file.Seek(int64(mkfile.StartPartition), 0)
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

func (mkfile *Mkfile) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
func (mkfile *Mkfile) ReturnDate8Bytes() [8]byte {
	t := string(time.Now().Format("02012006"))
	tmpT := []byte(t)
	return [8]byte(tmpT)
}
