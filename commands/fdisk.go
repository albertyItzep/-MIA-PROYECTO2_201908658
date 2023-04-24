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
	if tmp.Fit == 'o' {
		tmp.Fit = 'w'
	}
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
		tmp.ExtendPartition()
	} else if tmp.Type == 'l' {
		tmp.LogicPartition()
	}
	fmt.Println(tmp.MbrFdisk)
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
			if tmp.MbrFdisk.Dsk_fit == 'f' || tmp.MbrFdisk.Dsk_fit == 'o' {
				partitionTmp.Part_start = uint32(tmp.MemoryList.FirstSpace(int(partitionTmp.Part_size)))
			} else if tmp.MbrFdisk.Dsk_fit == 'b' {
				partitionTmp.Part_start = uint32(tmp.MemoryList.MinSpace(int(partitionTmp.Part_size)))
			} else if tmp.MbrFdisk.Dsk_fit == 'w' {
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

/*This function generate the asignatio of extend partition in the disk*/
func (fdisk *Fdisk) ExtendPartition() {
	if !fdisk.ExistExtendedPartition() {
		partitionTmp := fdisk.FreePartition()
		if partitionTmp == nil {
			fmt.Println("No existe particion Libre")
			return
		}
		//add of partition fit, type, name
		isertDisc := fdisk.StatusMemory()
		fmt.Println("estado de memoria apto")
		fdisk.MemoryList.ShowList()
		if isertDisc {
			partitionTmp.Part_status = 'o'
			partitionTmp.Part_type = 'e'
			partitionTmp.Part_size = fdisk.Size
			partitionTmp.Part_fit = fdisk.Fit
			for i := 0; i < 16; i++ {
				partitionTmp.Part_name[i] = fdisk.Name[i]
				if i == len(fdisk.Name)-1 {
					break
				}
			}
			if fdisk.MemoryList.ExistSpace(int(partitionTmp.Part_size)) {
				if fdisk.MbrFdisk.Dsk_fit == 'f' || fdisk.MbrFdisk.Dsk_fit == 'o' {
					fdisk.MbrFdisk.Dsk_fit = 'f'
					partitionTmp.Part_start = uint32(fdisk.MemoryList.FirstSpace(int(partitionTmp.Part_size)))
				} else if fdisk.MbrFdisk.Dsk_fit == 'b' {
					partitionTmp.Part_start = uint32(fdisk.MemoryList.MinSpace(int(partitionTmp.Part_size)))
				} else if fdisk.MbrFdisk.Dsk_fit == 'w' {
					partitionTmp.Part_start = uint32(fdisk.MemoryList.MajSpace(int(partitionTmp.Part_size)))
				}
				//escribimos en el disco el mbr
				file2, err := os.OpenFile(fdisk.Path, os.O_RDWR, 0644)
				if err != nil {
					fmt.Println("\033[31m[Error] > Al abrir el archivo:", err, "\033[0m")
					return
				}
				defer file2.Close()

				file2.Seek(0, 0)
				fmt.Println(unsafe.Sizeof(fdisk.MbrFdisk))
				err2 := binary.Write(file2, binary.LittleEndian, &fdisk.MbrFdisk)
				if err2 != nil && !os.IsExist(err) {
					fmt.Println("Error en la escritura de la particion")
					log.Fatal(err2)
				}
				file2.Seek(int64(partitionTmp.Part_start), 0)
				ebr0 := structs.EBR{Part_status: 'f', Part_start: -1, Part_size: -1, Part_fit: fdisk.Fit}
				err2 = binary.Write(file2, binary.LittleEndian, &ebr0)
				if err2 != nil {
					fmt.Println("Error en la escritura de la particion")
					log.Fatal(err2)
				}

				fdisk.StatusMemory()
			} else {
				fmt.Println("El disco se encuentra fragmentado por ello no se encuentra el espacio disponible")
			}

		} else {
			fmt.Println("Todas las particiones se encuentran ocupadas")
		}
	} else {
		fmt.Println("Existe una particion Extendida en el disco")
	}
}

/*This function asigned the logic partition in the disk*/
func (fdisk *Fdisk) LogicPartition() {
	if fdisk.ExistExtendedPartition() {
		extendedPartition := fdisk.ReturnExtendedPartition()
		if extendedPartition != nil {
			if extendedPartition.Part_size > fdisk.Size {
				file, err := os.OpenFile(fdisk.Path, os.O_RDWR, 0644)
				if err != nil {
					fmt.Println("Error al abrir el archivo")
					log.Fatal(err)
				}
				defer file.Close()
				//read the initial ebr
				var ebr structs.EBR
				file.Seek(int64(extendedPartition.Part_start), 0)
				err = binary.Read(file, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Println("Error al abrir el archivo")
					log.Fatal(err)
				}
				if ebr.Part_status == 'f' {
					ebr.Part_status = 'o'
					ebr.Part_fit = fdisk.Fit
					for i := 0; i < 16; i++ {
						ebr.Part_name[i] = fdisk.Name[i]
						if i == len(fdisk.Name)-1 {
							break
						}
					}
					ebr.Part_size = int32(fdisk.Size)
					ebr.Part_start = int32(extendedPartition.Part_start)
					file.Seek(int64(extendedPartition.Part_start), 0)
					err = binary.Write(file, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al escribir en el archivo")
						log.Fatal(err)
					}
					fmt.Println("Inserto una en el primer espacio vacio")
				} else {
					listTmp := structs.SpacesList{}
					for ebr.Part_next != 0 {
						listTmp.InsertNode(int(ebr.Part_start), int(ebr.Part_size+ebr.Part_start), 'o')
						file.Seek(int64(ebr.Part_next), 0)
						err = binary.Read(file, binary.LittleEndian, &ebr)
						if err != nil {
							fmt.Println("Error al leer el archivo")
							log.Fatal(err)
						}
					}
					listTmp.InsertNode(int(ebr.Part_start), int(ebr.Part_start+ebr.Part_size), 'o')
					file.Seek(int64(extendedPartition.Part_start), 0)
					listTmp.FillList(int(extendedPartition.Part_size))
					var startLogicP, freeSpace int
					freeSpace = int(extendedPartition.Part_size) - listTmp.SpaceFill
					if freeSpace > int(fdisk.Size) {
						if extendedPartition.Part_fit == 'f' {
							startLogicP = listTmp.FirstSpace(int(fdisk.Size))
						} else if extendedPartition.Part_fit == 'b' {
							startLogicP = listTmp.MinSpace(int(fdisk.Size))
						} else if extendedPartition.Part_fit == 'w' {
							startLogicP = listTmp.MajSpace(int(fdisk.Size))
						}
						nextPartitio := listTmp.NextSpace(startLogicP)
						tmp := structs.EBR{Part_status: 'o', Part_fit: fdisk.Fit, Part_start: int32(startLogicP), Part_size: int32(fdisk.Size), Part_next: int32(nextPartitio)}
						for i := 0; i < 16; i++ {
							tmp.Part_name[i] = fdisk.Name[i]
							if i == len(fdisk.Name)-1 {
								break
							}
						}
						file.Seek(int64(tmp.Part_start), 0)
						err = binary.Write(file, binary.LittleEndian, &tmp)
						if err != nil {
							fmt.Println("Error al escribir en el archivo")
							log.Fatal(err)
						}
						previusPart := listTmp.PreviusSpace(startLogicP)
						if previusPart != -1 {
							file.Seek(int64(previusPart), 0)
							var tmp2 structs.EBR
							err = binary.Read(file, binary.LittleEndian, &tmp2)
							if err != nil {
								fmt.Println("Error al leer el archivo")
								log.Fatal(err)
							}
							tmp2.Part_next = int32(startLogicP)
							file.Seek(int64(tmp2.Part_start), 0)
							err = binary.Write(file, binary.LittleEndian, &tmp2)
							if err != nil {
								fmt.Println("Error al escribir en el archivo")
								log.Fatal(err)
							}
							fmt.Println("Escribio una logica mas")
						}
					} else {
						fmt.Println("El espacio para la particion de momento no se encuentra disponible")
					}
				}
			} else {
				fmt.Println("Espacio insuficiente")
			}

		} else {
			fmt.Println("No se encontro la particion Extendida")
		}
	} else {
		fmt.Println("No existe particion Extendida")
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

/*This function return the if exist a extend partition*/
func (fdisk *Fdisk) ExistExtendedPartition() bool {
	if fdisk.MbrFdisk.Mbr_partition1.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition1.Part_type == 'e' {
		return true
	}
	if fdisk.MbrFdisk.Mbr_partition2.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition2.Part_type == 'e' {
		return true
	}
	if fdisk.MbrFdisk.Mbr_partition3.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition3.Part_type == 'e' {
		return true
	}
	if fdisk.MbrFdisk.Mbr_partition4.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition4.Part_type == 'e' {
		return true
	}
	return false
}

/*This function return the extend partition*/
func (fdisk *Fdisk) ReturnExtendedPartition() *structs.Partition {
	if fdisk.MbrFdisk.Mbr_partition1.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition1.Part_type == 'e' {
		return &fdisk.MbrFdisk.Mbr_partition1
	} else if fdisk.MbrFdisk.Mbr_partition2.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition2.Part_type == 'e' {
		return &fdisk.MbrFdisk.Mbr_partition2
	} else if fdisk.MbrFdisk.Mbr_partition3.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition3.Part_type == 'e' {
		return &fdisk.MbrFdisk.Mbr_partition3
	} else if fdisk.MbrFdisk.Mbr_partition4.Part_status == 'o' && fdisk.MbrFdisk.Mbr_partition4.Part_type == 'e' {
		return &fdisk.MbrFdisk.Mbr_partition4
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

//fdisk >Size=300 >path=/home/user/Disco1.dsk >name=Particion1
