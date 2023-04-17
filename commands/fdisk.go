package commands

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unsafe"

	"github.com/mia/proyecto2/structs"
)

type Fdisk struct {
	Name, Path      string
	Fit, Type, Unit byte
	Size            uint32
	MbrFdisk        structs.MBR
	MemoryList      structs.SpacesList
}

func (tmp *Fdisk) Execute() {
	tmp.Name = tmp.ReturnValueWithoutMarks(tmp.Name)
	tmp.Path = tmp.ReturnValueWithoutMarks(tmp.Path)
	tmp.Size = uint32(tmp.ReturnSize(int(tmp.Size), tmp.Unit))
	//Leemos el mbr del disco
	file, err := os.Open(tmp.Path)
	if err != nil {
		fmt.Println("\033[31m[Error] > Al abrir el archivo:", err, "\033[0m")
		return
	}
	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &tmp.MbrFdisk)
	if err != nil {
		fmt.Println("\033[31m [Error] > Al leer el archivo:", err, "\033[0m")
		return
	}
	file.Close()
	//verificamos el tipo de particion a crear
	if tmp.Type == 'p' || tmp.Type == 'o' {
		fmt.Println(tmp.MbrFdisk)
		tmp.PrimariPartition()
	} else if tmp.Type == 'e' {
		fmt.Println("Es extendida")
	} else if tmp.Type == 'l' {
		fmt.Println("Es logica")
	}
}

/*The function generate the asignation of a primary partition in the disk*/
func (tmp *Fdisk) PrimariPartition() {
	partitionTmp := tmp.FreePartition()
	if partitionTmp == nil {
		fmt.Println("No existe particion Libre")
		return
	}
	//add of partition fit, type, name
	isertDisc := tmp.StatusMemory()
	fmt.Println("estado de memoria apto")
	tmp.MemoryList.ShowList()
	if isertDisc {
		partitionTmp.Part_fit = tmp.Fit
		partitionTmp.Part_status = 'o'
		partitionTmp.Part_type = 'p'
		partitionTmp.Part_size = tmp.Size

		for i := 0; i < 16; i++ {
			partitionTmp.Part_name[i] = tmp.Name[i]
			if i == len(tmp.Name)-1 {
				break
			}
		}
		if tmp.MemoryList.ExistSpace(int(partitionTmp.Part_size)) {
			if tmp.Fit == 'f' {
				partitionTmp.Part_start = uint32(tmp.MemoryList.FirstSpace(int(partitionTmp.Part_size)))
			} else if tmp.Fit == 'b' {
				partitionTmp.Part_start = uint32(tmp.MemoryList.MinSpace(int(partitionTmp.Part_size)))
			} else if tmp.Fit == 'w' {
				partitionTmp.Part_start = uint32(tmp.MemoryList.MajSpace(int(partitionTmp.Part_size)))
			}
			//escribimos en el disco el mbr
			file2, err := os.OpenFile(tmp.Path, os.O_RDWR, 0644)
			if err != nil {
				fmt.Println("\033[31m[Error] > Al abrir el archivo:", err, "\033[0m")
				return
			}
			defer file2.Close()

			file2.Seek(0, 0)
			fmt.Println(unsafe.Sizeof(tmp.MbrFdisk))
			err2 := binary.Write(file2, binary.LittleEndian, &tmp.MbrFdisk)
			if err2 != nil && !os.IsExist(err) {
				fmt.Println("aqui es")
				log.Fatal(err2)
			}
			tmp.StatusMemory()
		} else {
			fmt.Println("El disco se encuentra fragmentado por ello no se encuentra el espacio disponible")
		}

	} else {
		fmt.Println("Todas las particiones se encuentran ocupadas")
	}
}

/*The function verify the status of dispotition of memory*/
func (fdiskTmp *Fdisk) StatusMemory() bool {
	cantPartCreated := 0
	fmt.Println(cantPartCreated)
	fdiskTmp.MemoryList.ClearList()
	sizeMBR := unsafe.Sizeof(structs.MBR{})
	fdiskTmp.MemoryList.InsertNode(0, int(sizeMBR), 'o')
	if fdiskTmp.MbrFdisk.Mbr_partition1.Part_status == 'o' {
		cantPartCreated++
		fdiskTmp.MemoryList.InsertNode(int(fdiskTmp.MbrFdisk.Mbr_partition1.Part_start), (int(fdiskTmp.MbrFdisk.Mbr_partition1.Part_size) + int(fdiskTmp.MbrFdisk.Mbr_partition1.Part_start)), 'o')
	}
	if fdiskTmp.MbrFdisk.Mbr_partition2.Part_status == 'o' {
		cantPartCreated++
		fdiskTmp.MemoryList.InsertNode(int(fdiskTmp.MbrFdisk.Mbr_partition2.Part_start), (int(fdiskTmp.MbrFdisk.Mbr_partition2.Part_size) + int(fdiskTmp.MbrFdisk.Mbr_partition2.Part_start)), 'o')
	}
	if fdiskTmp.MbrFdisk.Mbr_partition3.Part_status == 'o' {
		cantPartCreated++
		fdiskTmp.MemoryList.InsertNode(int(fdiskTmp.MbrFdisk.Mbr_partition3.Part_start), (int(fdiskTmp.MbrFdisk.Mbr_partition3.Part_size) + int(fdiskTmp.MbrFdisk.Mbr_partition3.Part_start)), 'o')
	}
	if fdiskTmp.MbrFdisk.Mbr_partition4.Part_status == 'o' {
		cantPartCreated++
		fdiskTmp.MemoryList.InsertNode(int(fdiskTmp.MbrFdisk.Mbr_partition4.Part_start), (int(fdiskTmp.MbrFdisk.Mbr_partition4.Part_size) + int(fdiskTmp.MbrFdisk.Mbr_partition4.Part_start)), 'o')
	}
	fdiskTmp.MemoryList.FillList(int(fdiskTmp.MbrFdisk.Mbr_tamano))
	if cantPartCreated < 4 {
		fmt.Print("")
		return true
	}
	return false
}

/*The function return the free partition in the disk*/
func (tmp *Fdisk) FreePartition() *structs.Partition {
	if tmp.MbrFdisk.Mbr_partition1.Part_status == 'f' {
		return &tmp.MbrFdisk.Mbr_partition1
	} else if tmp.MbrFdisk.Mbr_partition2.Part_status == 'f' {
		return &tmp.MbrFdisk.Mbr_partition2
	} else if tmp.MbrFdisk.Mbr_partition3.Part_status == 'f' {
		return &tmp.MbrFdisk.Mbr_partition3
	} else if tmp.MbrFdisk.Mbr_partition4.Part_status == 'f' {
		return &tmp.MbrFdisk.Mbr_partition4
	}
	return nil
}

/*The function return a value without marks*/
func (tmp *Fdisk) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}

/*Return the size in bytes*/
func (tmp *Fdisk) ReturnSize(sizeI int, unit byte) int {
	var size int
	switch {
	case unit == 'b':
		size = sizeI
	case unit == 'k' || tmp.Unit == 'o':
		size = sizeI * 1024
	case unit == 'm':
		size = sizeI * 1024 * 1024
	}
	return size
}
