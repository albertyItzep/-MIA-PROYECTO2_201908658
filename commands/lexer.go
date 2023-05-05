package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mia/proyecto2/structs"
)

type Lexer struct {
	CommandString string
	ListDisk      structs.DiskList
	ListPartitio  structs.PartitionList
	ListMount     MountList
	UserLoged     LoginUser
}

/* take a string and searched a command defined, execute a function for the command */
func (tmp *Lexer) GeneralComand(command string) string {
	tmp.CommandString = command
	if matched, _ := regexp.Match("(mkdisk)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandMkdisk()
	} else if matched, _ := regexp.Match("(rmdisk)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandRmdisk()
	} else if matched, _ := regexp.Match("(fdisk)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandFdisk()
	} else if matched, _ := regexp.Match("(mount)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandMount()
	} else if matched, _ := regexp.Match("(mkfs)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandMkfs()
	} else if matched, _ := regexp.Match("(rep)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandRep()
	} else if matched, _ := regexp.Match("(pause)(.*)", []byte(tmp.CommandString)); matched {
		fmt.Println("contiene el comando pause")
	} else if matched, _ := regexp.Match("(login)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandLogin()
	} else if matched, _ := regexp.Match("(logout)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandLogout()
	} else if matched, _ := regexp.Match("(mkgrp)(.*)", []byte(tmp.CommandString)); matched {
		tmp.CommandMkgrp()
	} else if matched, _ := regexp.Match("(rmgrp)(.*)", []byte(tmp.CommandString)); matched {
		fmt.Println("contiene el comando rmgrp")
	} else if matched, _ := regexp.Match("(mkuser)(.*)", []byte(tmp.CommandString)); matched {
		fmt.Println("contiene el comando mkuser")
	} else if matched, _ := regexp.Match("(rmuser)(.*)", []byte(tmp.CommandString)); matched {
		fmt.Println("contiene el comando rmuser")
	} else if matched, _ := regexp.Match("(mkfile)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandMkfile()
	} else if matched, _ := regexp.Match("(mkdir)(.*)", []byte(tmp.CommandString)); matched {
		return tmp.CommandMkdir()
	}
	return "Error"
}

/* This method is used for make file of types binaries, with the structure implemented*/
func (tmp *Lexer) CommandMkdisk() string {
	pathMkdir := tmp.PathParameter(true)
	size := tmp.SizeParameter(true)
	fit := tmp.FitParameter(false)
	unit := tmp.UnitParameter(false)
	if pathMkdir != "" && size > 0 {
		tmpM := Mkdisk{Path: pathMkdir, Fit: fit, Unit: unit, Size: size}
		tmpM.Execute()
		size = tmpM.ReturnSize(size, unit)
		tmp.ListDisk.InsertNode(pathMkdir, size)
		return "Disco creado exitosamente ..."
	}
	return "Error en Mkdisk"
}

/*This method is used for delete file binari*/
func (tmp *Lexer) CommandRmdisk() string {
	path := tmp.PathParameter(true)
	if path != "" {
		rmdisk := Rmdisk{Path: path}
		return rmdisk.Execute()
	}
	return "Error al elimiar el disco"
}

/*This method is used for modifi ofs patrtitions*/
func (tmp *Lexer) CommandFdisk() string {
	name := tmp.NameParameter(true)
	fit := tmp.FitParameter(false)
	typeFdisk := tmp.TypeParameter(false)
	pathFdisk := tmp.PathParameter(true)
	sizeFdisk := tmp.SizeParameter(true)
	unitfdisk := tmp.UnitParameter(false)
	fdisk := Fdisk{Name: name, Path: pathFdisk, Fit: fit, Type: typeFdisk, Size: uint32(sizeFdisk), Unit: unitfdisk}
	estatusR := fdisk.Execute()
	existDisk := tmp.ListDisk.ExistDiscList(pathFdisk)
	if !existDisk {
		if tmp.ListDisk.ExistFileFisic(pathFdisk) {
			tmp.ListDisk.InsertNode(pathFdisk, tmp.ListDisk.ReturnFileSizeFisic(pathFdisk))
		}
	}
	if estatusR == "Particion creada exitosamente ..." {
		existDisk = tmp.ListDisk.ExistDiscList(pathFdisk)
		sizeFdisk = fdisk.ReturnSize(sizeFdisk, unitfdisk)
		tmp.ListDisk.InsertPartitionDisk(pathFdisk)
		tmp.ListPartitio.InsertNode(pathFdisk, name, sizeFdisk, existDisk)
	}
	return estatusR
}

/*The function execute the mount command*/
func (tmp *Lexer) CommandMount() string {
	name := tmp.NameParameter(true)
	pathFile := tmp.PathParameter(true)
	starPartition := tmp.ListPartitio.ReturnStartPartitionValue(pathFile, name)
	tmp.ListDisk.InsertPartitionDiskMounted(pathFile)
	numParition := tmp.ListDisk.ReturnPartitionsDiskMounted(pathFile)
	idDisk := tmp.ListDisk.ReturnIdOfPartition(pathFile) + 1
	sizePartition := tmp.ListPartitio.ReturnSizePartition(pathFile, name)

	tmp.ListMount.InserMount(pathFile, name, starPartition, sizePartition, numParition, idDisk)
	return "particion montada con exito ..."
}

/*The functio execute the mkfs command*/
func (tmp *Lexer) CommandMkfs() string {
	id := tmp.IdParameter(true)
	typePar := tmp.TypeMkfsParameter(false)
	SizeOfPartition := tmp.ListMount.ReturnSizeWithId(id)
	if id != "no" {
		startPartition := tmp.ListMount.ReturnStartPartitionWithId(id)
		pathFile := tmp.ListMount.ReturnPathitionWithId(id)
		mkfs := Mkfs{IdMkfs: id, TypeMkfs: typePar, SizeOfPartition: SizeOfPartition, StartPartition: startPartition, PathFile: pathFile}
		mkfs.Execute()
		return "Formateo realizado con exito ..."
	}
	return "Error al generar el formateo"
}

/*The function execute the login command*/
func (tmp *Lexer) CommandLogin() string {
	idPartition := tmp.IdParameter(true)
	userLogin := tmp.UserParameter(true)
	passLogin := tmp.PasswordParameter(true)
	if idPartition != "no" && userLogin != "no" && passLogin != "" {
		pathFile := tmp.ListMount.ReturnPathitionWithId(idPartition)
		startPartition := tmp.ListMount.ReturnStartPartitionWithId(idPartition)
		tmp.UserLoged.IdPartition = idPartition
		tmp.UserLoged.User = userLogin
		tmp.UserLoged.Pwd = passLogin
		tmp.UserLoged.StartPartition = startPartition
		tmp.UserLoged.PathFile = pathFile
		if tmp.UserLoged.LogedUser() {
			return "SA"
		} else {
			return tmp.UserLoged.Execute()
		}
	}
	return "EL"
}

/*The function execute the logout user*/
func (tmp *Lexer) CommandLogout() string {
	if tmp.UserLoged.LogedUser() {
		tmp.UserLoged.User = ""
		tmp.UserLoged.IdPartition = ""
		tmp.UserLoged.Pwd = ""
		tmp.UserLoged.PathFile = ""
		tmp.UserLoged.StartPartition = -1
		tmp.UserLoged.Loged = false
		return "Nos vemos pronto"
	} else {
		fmt.Println("No existe sesion activa")
	}
	return "No existe Secion Activa"
}

/*The function execute the make a grup command*/
func (tmp *Lexer) CommandMkgrp() {
	if tmp.UserLoged.LogedUser() {
		if tmp.UserLoged.User == "root" {
			nameGrup := tmp.NameMkgrupParameter(true)
			pathFile := tmp.UserLoged.PathFile
			startPartition := tmp.UserLoged.StartPartition
			mkgrp := Mkgrp{NameGrup: nameGrup, PathFile: pathFile, StartParition: startPartition}
			mkgrp.Execute()
		} else {
			fmt.Println("Permisos no validos, utilice usuario root")
		}
	} else {
		fmt.Println("Sesion invalida")
	}
}

/*the function execute the make a file*/
func (tmp *Lexer) CommandMkfile() string {
	pathNewFile := tmp.PathParameter(true)
	sizeNewFile := tmp.SizeParameter(false)
	startPartition := 133
	PathDisk := "/home/user/disco1.dsk"
	mkfile := Mkfile{PathNewFile: pathNewFile, PathDisk: PathDisk, SizeNewFile: sizeNewFile, StartPartition: startPartition}
	mkfile.Execute()
	return ""
}

/*the function execute the make a dir*/
func (tmp *Lexer) CommandMkdir() string {
	if tmp.UserLoged.LogedUser() {
		pathFile := tmp.PathParameter(true)
		startPartition := tmp.UserLoged.StartPartition
		PathDisk := tmp.ListMount.ReturnPathitionWithId(tmp.UserLoged.IdPartition)
		mkdir := Mkdir{PathNewDir: pathFile, PathDisk: PathDisk, StartPartition: startPartition, CreatePrevius: true}
		return mkdir.Execute()
	}
	return "Login Necesario"
}

/*The function create the code for reports*/
func (tmp *Lexer) CommandRep() string {
	if tmp.UserLoged.LogedUser() {
		fmt.Println("hello")
		nameRep := tmp.NameParameter(true)
		pathRep := tmp.PathParameter(true)
		idRep := tmp.IdParameter(true)
		pathRepFile := tmp.PathParameter(false)
		PathDisk := tmp.ListMount.ReturnPathitionWithId(idRep)
		startPartition := tmp.ListMount.ReturnStartPartitionWithId(idRep)
		sizeDisk := tmp.ListDisk.ReturSizeDisk(PathDisk)
		fmt.Println(nameRep)
		rep := Report{NameReport: nameRep, PathRep: pathRep, idPartition: idRep, PathRepFile: pathRepFile, PathDisk: PathDisk, SizeDisk: sizeDisk, StartParition: startPartition}
		return rep.Execute()
	}
	return "No Login"
}

/*The parameter Name contain the name of partition*/
func (tmp *Lexer) NameParameter(obligatory bool) string {
	cadena := ""
	if matched, _ := regexp.Match(">name=[\"?[a-zA-Z0-9\\_]+\"?", []byte(tmp.CommandString)); matched {
		regeName := regexp.MustCompile(">name=[\"?[a-zA-Z0-9\\_]+\"?")
		content := regeName.FindAllString(tmp.CommandString, -1)
		if len(content) > 0 {
			cadena = content[0]
			remplace := regexp.MustCompile(">name=")
			cadena = remplace.ReplaceAllString(cadena, "")
		}
	} else if obligatory {
		fmt.Println("El parametro size es obligatorio")
		cadena = ""
	} else {
		cadena = ""
	}
	return cadena
}

/*The parameter contain the name of the grup*/
func (tmp *Lexer) NameMkgrupParameter(obligatory bool) string {
	cadena := ""
	if matched, _ := regexp.Match(">name=[\"?[a-zA-Z0-9\\_[:space:]]+\"?", []byte(tmp.CommandString)); matched {
		regeName := regexp.MustCompile(">name=[\"?[a-zA-Z0-9\\_[:space:]]+\"?")
		content := regeName.FindAllString(tmp.CommandString, -1)
		if len(content) > 0 {
			cadena = content[0]
			cadena = strings.Trim(cadena, ">name=")
		}
	} else if obligatory {
		fmt.Println("El parametro size es obligatorio")
		cadena = ""
	} else {
		cadena = ""
	}
	return cadena
}

/*The parameter Type contain the type of partition*/
func (tmp *Lexer) TypeParameter(obligatory bool) byte {
	var typeTmp byte
	if matched, _ := regexp.Match(">type=[a-zA-Z]", []byte(tmp.CommandString)); matched {
		regeType := regexp.MustCompile(">type=[a-zA-Z]")
		content := regeType.FindAllString(tmp.CommandString, -1)
		if len(content) > 0 {
			tmpString := content[0]
			replace := regexp.MustCompile(">type=")
			tmpString = replace.ReplaceAllString(tmpString, "")
			typeTmp = tmpString[0]
		}
	} else if obligatory {
		fmt.Println("El parametro Type es obligatorio")
	} else {
		typeTmp = 'o'
	}
	return typeTmp
}

/* The parameter path contain the path where create file bin*/
func (tmp *Lexer) PathParameter(obligatory bool) string {
	result1 := ""
	if matchedPaht, _ := regexp.Match(">path=", []byte(tmp.CommandString)); matchedPaht {
		//verificamos que la path con comillas
		if matchedPaht1, _ := regexp.Match("\"(/.*)+/[a-zA-Z0-9]+.dsk\"", []byte(tmp.CommandString)); matchedPaht1 {
			regePath := regexp.MustCompile("\"(/.*)+/[a-zA-Z0-9]+.dsk\"")
			content := regePath.FindAllString(tmp.CommandString, -1)
			if len(content) > 0 {
				result1 = content[0]
			}
		} else if matchedPaht1, _ := regexp.Match("(/[a-zA-Z0-9\\.]+)+", []byte(tmp.CommandString)); matchedPaht1 {
			regePath := regexp.MustCompile("(/[a-zA-Z0-9...]+)+")
			content := regePath.FindAllString(tmp.CommandString, -1)
			if len(content) > 0 {
				result1 = content[0]
			}
		}
	} else if obligatory {
		fmt.Println("El parametro Path es obligatorio")
		result1 = ""
	} else if !obligatory {
		result1 = ""
	}
	return result1
}

/* The parameter size contain the size of file or partition or a diferent object */
func (tmp *Lexer) SizeParameter(obligatory bool) int {
	var size int
	var tmpString string
	if matched, _ := regexp.Match(">size=", []byte(tmp.CommandString)); matched {
		regeSize := regexp.MustCompile(">size=[0-9]+")
		content := regeSize.FindAllString(tmp.CommandString, -1)
		if len(content) > 0 {
			tmpString = content[0]
			tmpString = strings.Trim(tmpString, ">size=")
			size, _ = strconv.Atoi(tmpString)
		}
	} else if obligatory {
		fmt.Println("El parametro size es obligatorio")
		size = -1
	} else {
		size = 0
	}
	return size
}

/* This parameter fit contain the configuration of asignation of disk or partition*/
func (tmp *Lexer) FitParameter(obligatory bool) byte {
	var fit byte
	var tmpString string
	if matched, _ := regexp.Match(">fit=", []byte(tmp.CommandString)); matched {
		regexFit := regexp.MustCompile(">fit=(bf|ff|wf)")
		content := regexFit.FindAllString(tmp.CommandString, -1)
		tmpString = content[0]
		remplace := regexp.MustCompile(">fit=")
		tmpString = remplace.ReplaceAllString(tmpString, "")
		switch {
		case tmpString == "ff":
			fit = 'f'
		case tmpString == "bf":
			fit = 'b'
		case tmpString == "wf":
			fit = 'w'
		}
	} else if obligatory {
		fmt.Println("El parametro fit es obligatorio")
		fit = 'n'
	} else if !obligatory { // is optional, i can asigned orther value
		fit = 'o'
	}
	return fit
}

/* The parameter unit contain de information respect the type of storage in the disk or partition*/
func (tmp *Lexer) UnitParameter(obligatory bool) byte {
	var unit byte
	var tmpString string
	if matched, _ := regexp.Match(">unit=", []byte(tmp.CommandString)); matched {
		regexFit := regexp.MustCompile(">unit=(k|m)")
		content := regexFit.FindAllString(tmp.CommandString, -1)
		tmpString = content[0]
		remplace := regexp.MustCompile(">unit=")
		tmpString = remplace.ReplaceAllString(tmpString, "")
		switch {
		case tmpString == "m":
			unit = 'm'
		case tmpString == "k":
			unit = 'k'
		}
	} else if obligatory {
		fmt.Println("El parametro unit es obligatorio")
		unit = 'n'
	} else if !obligatory { // is optional, i can asigned orther value
		unit = 'o'
	}
	return unit
}

/*The id parameter contain the information respect of the partition use*/
func (tmp *Lexer) IdParameter(obligatory bool) string {
	var text string
	if matched, _ := regexp.Match(">id=", []byte(tmp.CommandString)); matched {
		regexId := regexp.MustCompile(">id=[a-zA-Z0-9]+")
		content := regexId.FindAllString(tmp.CommandString, -1)
		text = content[0]
		remplace := regexp.MustCompile(">id=")
		text = remplace.ReplaceAllString(text, "")
	} else if obligatory {
		fmt.Println("El parametro id es obligatorio")
		text = "no"
	} else if !obligatory {
		text = "no"
	}
	return text
}

/*The function contain the password parameter*/
func (tmp *Lexer) PasswordParameter(obligatory bool) string {
	var tmpString string
	passWithoutMarks := regexp.MustCompile(">pwd=[a-zA-Z0-9]+")
	passWithMarks := regexp.MustCompile(">pwd=\"[a-zA-Z0-9[:space:]]+\"")
	remplace := regexp.MustCompile(">pwd=")
	if matched, _ := regexp.Match(">pwd=[a-zA-Z0-9]+", []byte(tmp.CommandString)); matched {
		content := passWithoutMarks.FindAllString(tmp.CommandString, -1)
		tmpString = content[0]
		tmpString = remplace.ReplaceAllString(tmpString, "")
	} else if matched1, _ := regexp.Match(">pwd=\"[a-zA-Z0-9[:space:]]+\"", []byte(tmp.CommandString)); matched1 {
		content := passWithMarks.FindAllString(tmp.CommandString, -1)
		tmpString = content[0]
		tmpString = remplace.ReplaceAllString(tmpString, "")
	} else if obligatory {
		fmt.Println("El parametro id es obligatorio")
		tmpString = "no"
	} else if !obligatory {
		tmpString = "no"
	}
	return tmpString
}

/*The funtion contain the user parameter*/
func (tmp *Lexer) UserParameter(obligatory bool) string {
	var tmpString string
	if matched, _ := regexp.Match(">user=", []byte(tmp.CommandString)); matched {
		regexUserWithoutMarks := regexp.MustCompile(">user=[a-zA-Z0-9]+")
		regexUserWithMarks := regexp.MustCompile(">user=\"[a-zA-Z0-9[:space:]]+\"")
		if matched1, _ := regexp.Match(">user=[a-zA-Z0-9]+", []byte(tmp.CommandString)); matched1 {
			content := regexUserWithoutMarks.FindAllString(tmp.CommandString, -1)
			tmpString = content[0]
			remplace := regexp.MustCompile(">user=")
			tmpString = remplace.ReplaceAllString(tmpString, "")
		} else if matched1, _ := regexp.Match(">user=\"[a-zA-Z0-9[:space:]]+\"", []byte(tmp.CommandString)); matched1 {
			content := regexUserWithMarks.FindAllString(tmp.CommandString, -1)
			tmpString = content[0]
			remplace := regexp.MustCompile(">user=")
			tmpString = remplace.ReplaceAllString(tmpString, "")
		}
	} else if obligatory {
		fmt.Println("El parametro id es obligatorio")
		tmpString = "no"
	} else if !obligatory {
		tmpString = "no"
	}

	return tmpString
}

/*The function contain the type of formating install in the partition*/
func (tmp *Lexer) TypeMkfsParameter(obligatory bool) string {
	var stringTmp string
	if matched, _ := regexp.Match(">type=full", []byte(tmp.CommandString)); matched {
		stringTmp = "full"
	} else if obligatory {
		fmt.Println("El parametro type es obligatorio")
		stringTmp = "no"
	} else if !obligatory {
		stringTmp = "full"
	}
	return stringTmp
}
