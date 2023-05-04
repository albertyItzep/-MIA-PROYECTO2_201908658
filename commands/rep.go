package commands

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/mia/proyecto2/structs"
)

type Report struct {
	NameReport, PathRep, PathRepFile string
	idPartition, PathDisk            string
	SizeDisk, StartParition          int
}

func (rep *Report) Execute() string {
	rep.NameReport = rep.ReturnValueWithoutMarks(rep.NameReport)
	rep.PathRep = rep.ReturnValueWithoutMarks(rep.PathRep)
	rep.PathRepFile = rep.ReturnValueWithoutMarks(rep.PathRepFile)
	rep.idPartition = rep.ReturnValueWithoutMarks(rep.PathRepFile)
	rep.PathDisk = rep.ReturnValueWithoutMarks(rep.PathDisk)

	if rep.NameReport == "disk" {
		return rep.RepDisk()
	} else if rep.NameReport == "sb" {
		return rep.RepSb()
	} /*else if rep.NameReport == "tree" {

	} else if rep.NameReport == "file" {

	}*/
	return ""
}

func (rep *Report) RepDisk() string {

	file, err := os.OpenFile(rep.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivoasd:", err, "\033[0m")
		return ""
	}
	defer file.Close()

	var mbr structs.MBR
	posIExtendida := 0
	particioEx := 0
	ebrList := structs.SpacesList{}
	mbrList := structs.SpacesList{}

	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al leer el archivo:", err, "\033[0m")
		return ""
	}
	if mbr.Mbr_partition1.Part_type == 'e' {
		posIExtendida = int(mbr.Mbr_partition1.Part_start)
		particioEx = 1
	}
	if mbr.Mbr_partition2.Part_type == 'e' {
		posIExtendida = int(mbr.Mbr_partition2.Part_start)
		particioEx = 2
	}
	if mbr.Mbr_partition3.Part_type == 'e' {
		posIExtendida = int(mbr.Mbr_partition3.Part_start)
		particioEx = 3
	}
	if mbr.Mbr_partition4.Part_type == 'e' {
		posIExtendida = int(mbr.Mbr_partition4.Part_start)
		particioEx = 4
	}

	mbrList.InsertForSize(0, int(unsafe.Sizeof(structs.MBR{})), int(unsafe.Sizeof(structs.MBR{})), 'o')
	if mbr.Mbr_partition1.Part_status == 'o' {
		mbrList.InsertNodeRep(int(mbr.Mbr_partition1.Part_start), int(mbr.Mbr_partition1.Part_start+mbr.Mbr_partition1.Part_size), int(mbr.Mbr_partition1.Part_size), 'o', mbr.Mbr_partition1.Part_type)
	}
	if mbr.Mbr_partition2.Part_status == 'o' {
		mbrList.InsertNodeRep(int(mbr.Mbr_partition2.Part_start), int(mbr.Mbr_partition2.Part_start+mbr.Mbr_partition2.Part_size), int(mbr.Mbr_partition2.Part_size), 'o', mbr.Mbr_partition2.Part_type)
	}
	if mbr.Mbr_partition3.Part_status == 'o' {
		mbrList.InsertNodeRep(int(mbr.Mbr_partition3.Part_start), int(mbr.Mbr_partition3.Part_start+mbr.Mbr_partition3.Part_size), int(mbr.Mbr_partition3.Part_size), 'o', mbr.Mbr_partition3.Part_type)
	}
	if mbr.Mbr_partition4.Part_status == 'o' {
		mbrList.InsertNodeRep(int(mbr.Mbr_partition4.Part_start), int(mbr.Mbr_partition4.Part_start+mbr.Mbr_partition4.Part_size), int(mbr.Mbr_partition4.Part_size), 'o', mbr.Mbr_partition4.Part_type)
	}
	mbrList.FillList(rep.SizeDisk)
	mbrList.ShowList()
	if posIExtendida > 0 {
		var ebr structs.EBR
		file.Seek(int64(posIExtendida), 0)
		err = binary.Read(file, binary.LittleEndian, &ebr)
		if err != nil {
			fmt.Println("\033[31m[Error] > Al abrir el archivo1:", err, "\033[0m")
			return ""
		}
		if ebr.Part_next > 0 {
			ebrList.InsertNodeRep(int(ebr.Part_start), int(ebr.Part_start+ebr.Part_size), int(ebr.Part_size), ebr.Part_status, 'l')
			for ebr.Part_next > 0 {
				file.Seek(int64(ebr.Part_next), 0)
				err = binary.Read(file, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al abrir el archivo4:", err, "\033[0m")
					return ""
				}
				if ebr.Part_status == 'o' {
					ebrList.InsertNodeRep(int(ebr.Part_start), int(ebr.Part_start+ebr.Part_size), int(ebr.Part_size), 'o', 'l')
				}
			}
		} else {
			if ebr.Part_status == 'o' {
				ebrList.InsertNodeRep(int(ebr.Part_start), int(ebr.Part_start+ebr.Part_size), int(ebr.Part_size), 'o', 'l')
			}
		}
		if particioEx == 1 {
			ebrList.FillList(int(mbr.Mbr_partition1.Part_size + mbr.Mbr_partition1.Part_start))
		} else if particioEx == 2 {
			ebrList.FillList(int(mbr.Mbr_partition2.Part_size + mbr.Mbr_partition2.Part_start))
		} else if particioEx == 3 {
			ebrList.FillList(int(mbr.Mbr_partition3.Part_size + mbr.Mbr_partition3.Part_start))
		} else if particioEx == 4 {
			ebrList.FillList(int(mbr.Mbr_partition4.Part_size + mbr.Mbr_partition4.Part_start))
		}
	}
	ebrList.ShowList()
	cadena := ""
	cadena = "digraph G { \n rankdir = LR;\n"
	cadena += " nodoG[shape=record label=\"{ MBR "

	for i := 0; i < mbrList.Size; i++ {
		typeIndex := mbrList.ReturnTypeIndex(i)
		if typeIndex == 'p' {
			cadena += "| Primaria \\n "
			tamano := mbrList.ReturnSizeIndex(i)
			var porcentaje float64
			porcentaje = float64(tamano) / float64(rep.SizeDisk)
			porcentaje = porcentaje * 100
			cadena += strconv.FormatFloat(porcentaje, 'f', 6, 64) + "% "
		} else if typeIndex == 'e' {
			cadena += "| { Extendida "
			if ebrList.Size > 0 {
				cadena += " | { "
				for j := 0; j < ebrList.Size; j++ {
					typeEbr := ebrList.ReturnTypeIndex(j)
					if j != 0 {
						cadena += " | "
					}
					if typeEbr == 'f' {
						cadena += " Espacio Libre "
						tamano := ebrList.ReturnSizeIndex(j)
						var porcentaje float64
						porcentaje = float64(tamano) / float64(rep.SizeDisk)
						porcentaje = porcentaje * 100
						cadena += strconv.FormatFloat(porcentaje, 'f', 6, 64) + "%"
					} else if typeEbr == 'l' {
						cadena += " EBR | Logica "
						tamano := ebrList.ReturnSizeIndex(j)
						var porcentaje float64
						porcentaje = float64(tamano) / float64(rep.SizeDisk)
						porcentaje = porcentaje * 100
						cadena += strconv.FormatFloat(porcentaje, 'f', 6, 64) + "% "
					}
				}
				cadena += " } "
			} else {
				tamano := mbrList.ReturnSizeIndex(i)
				var porcentaje float64
				porcentaje = float64(tamano) / float64(rep.SizeDisk)
				porcentaje = porcentaje * 100
				cadena += strconv.FormatFloat(porcentaje, 'f', 6, 64) + "% "
			}
			cadena += " } "
		} else if typeIndex == 'f' {
			cadena += "| Espacio Libre \\n "
			tamano := mbrList.ReturnSizeIndex(i)
			var porcentaje float64
			porcentaje = float64(tamano) / float64(rep.SizeDisk)
			porcentaje = porcentaje * 100
			cadena += strconv.FormatFloat(porcentaje, 'f', 6, 64) + "% "

		}
	}
	cadena += "}\"];\n"
	cadena += "label=\"Reporte disco\"; \n}"

	patDir := rep.ReturnPathWithoutFileName(rep.PathRep)
	err = os.MkdirAll(patDir, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = os.Chmod(patDir, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	patDir += "/rep.dot"
	file, err2 := os.Create(patDir)
	if err2 != nil && !os.IsExist(err) {
		log.Fatal(err2)
	}
	defer file.Close()

	err = os.Chmod(patDir, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	_, err = file.WriteString(cadena)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	//comand := "sudo dot -Tpng " + patDir + " -o " + rep.PathRep
	cmd := exec.Command("dot", "-Tpng", patDir, "-o", rep.PathRep)
	err = cmd.Run()
	if err != nil {
		fmt.Errorf("no se pudo generar la imagen: %v", err)
	}
	return cadena
}

func (rep *Report) RepSb() string {
	fmt.Println("hello")
	file, err := os.OpenFile(rep.PathDisk, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo:", err, "\033[0m")
		return ""
	}
	defer file.Close()
	//para esto debe estar formateada la particion
	superBlock := structs.SuperBlock{}
	file.Seek(int64(rep.StartParition), 0)
	err = binary.Read(file, binary.LittleEndian, &superBlock)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo:", err, "\033[0m")
		return ""
	}

	cadena := ""
	cadena += "digraph G { \n rankdir = TB;\n"
	cadena += "nodoG[shape=record label=\"{ Super Bloque "
	cadena += "| s_filesystem_type: " + strconv.Itoa(int(superBlock.S_filesystem_type)) + " "
	cadena += "| s_inodes_count: " + strconv.Itoa(int(superBlock.S_inodes_count)) + " "
	cadena += "| s_blocks_count: " + strconv.Itoa(int(superBlock.S_blocks_count)) + " "
	cadena += "| s_free_blocks_count: " + strconv.Itoa(int(superBlock.S_free_blocks_count)) + " "
	cadena += "| s_free_inodes_count: " + strconv.Itoa(int(superBlock.S_free_inodes_count)) + " "
	date := string(superBlock.S_mtime[:])
	cadena += "| s_mtime: " + strconv.Itoa(int(superBlock.S_filesystem_type)) + date + " "
	cadena += "| s_mnt_count: " + strconv.Itoa(int(superBlock.S_mnt_count)) + " "
	cadena += "| s_magic: " + strconv.Itoa(int(superBlock.S_magic)) + " "
	cadena += "| s_inode_size: " + strconv.Itoa(int(superBlock.S_inode_size)) + " "
	cadena += "| s_block_size: " + strconv.Itoa(int(superBlock.S_block_size)) + " "
	cadena += "| s_firts_ino: " + strconv.Itoa(int(superBlock.S_firts_ino)) + " "
	cadena += "| s_first_blo: " + strconv.Itoa(int(superBlock.S_first_blo)) + " "
	cadena += "| s_bm_inode_start: " + strconv.Itoa(int(superBlock.S_bm_inode_start)) + " "
	cadena += "| s_bm_block_start: " + strconv.Itoa(int(superBlock.S_bm_block_start)) + " "
	cadena += "| s_inode_start: " + strconv.Itoa(int(superBlock.S_inode_start)) + " "
	cadena += "| s_block_start: " + strconv.Itoa(int(superBlock.S_block_start)) + " "
	cadena += " }\"];\n"
	cadena += "label=\"Reporte Super Boque\";\n}"

	patDir := rep.ReturnPathWithoutFileName(rep.PathRep)
	err = os.MkdirAll(patDir, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = os.Chmod(patDir, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	patDir += "/rep.dot"
	file, err2 := os.Create(patDir)
	if err2 != nil && !os.IsExist(err) {
		log.Fatal(err2)
	}
	defer file.Close()

	err = os.Chmod(patDir, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	_, err = file.WriteString(cadena)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	//comand := "sudo dot -Tpng " + patDir + " -o " + rep.PathRep
	cmd := exec.Command("dot", "-Tpng", patDir, "-o", rep.PathRep)
	err = cmd.Run()
	if err != nil {
		fmt.Errorf("no se pudo generar la imagen: %v", err)
	}
	return cadena
}

func (tmp *Report) ReturnPathWithoutFileName(value string) string {
	var tmpString string
	regPahtMkdir := regexp.MustCompile("/[a-zA-Z0-9]+\\.[a-zA-Z]+")
	tmpString = regPahtMkdir.ReplaceAllString(value, "")
	return tmpString
}

/*The function return a value without marks*/
func (tmp *Report) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
