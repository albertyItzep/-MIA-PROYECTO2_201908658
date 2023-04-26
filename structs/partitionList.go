package structs

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type PartitionNode struct {
	Name, PathFile                string
	SizePartition, StartPartition int
	NextNode                      *PartitionNode
	PreviusNode                   *PartitionNode
}

type PartitionList struct {
	Size     int
	RootNode *PartitionNode
	EndNode  *PartitionNode
}

/*The function insert a node in the list*/
func (list *PartitionList) InsertNode(pathFile string, name string, SizePartition int, existDisk bool) {
	name = list.ReturnValueWithoutMarks(name)
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	if existDisk {
		startPartition := list.ReturnStartOfPartition(pathFile, name)
		if startPartition > 0 {
			tmp := PartitionNode{Name: name, PathFile: pathFile, SizePartition: SizePartition, StartPartition: startPartition}
			list.Size++
			if list.RootNode == nil {
				list.RootNode = &tmp
				list.EndNode = &tmp
			} else {
				list.EndNode.NextNode = &tmp
				tmp.PreviusNode = list.EndNode
				list.EndNode = &tmp
			}
		}
	}
}

/*The function return the start of partition*/
func (list *PartitionList) ReturnStartPartitionValue(pathFile string, name string) int {
	name = list.ReturnValueWithoutMarks(name)
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	tmp := list.RootNode
	for i := 0; i < list.Size; i++ {
		if tmp.Name == name && tmp.PathFile == pathFile {
			return tmp.StartPartition
		}
	}
	return -1
}

/*The function return the size of the partition*/
func (list *PartitionList) ReturnSizePartition(pathFile string, name string) int {
	name = list.ReturnValueWithoutMarks(name)
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	tmp := list.RootNode
	for i := 0; i < list.Size; i++ {
		if tmp.Name == name && tmp.PathFile == pathFile {
			return tmp.SizePartition
		}
	}
	return -1
}

/*The function return the start of the partition*/
func (list *PartitionList) ReturnStartOfPartition(pathFile string, name string) int {
	name = list.ReturnValueWithoutMarks(name)
	pathFile = list.ReturnValueWithoutMarks(pathFile)

	startPartition := 0

	file, err := os.OpenFile(pathFile, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
		log.Fatal(err)
	}
	defer file.Close()

	//read the MBR
	var mbr MBR
	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
		log.Fatal(err)
	}

	//validate the partitions
	isExtended := false

	startExtendedPartition := 0
	namePartition := list.ReturnNameString(mbr.Mbr_partition1.Part_name)
	if namePartition == name {
		startPartition = int(mbr.Mbr_partition1.Part_start)
	} else if mbr.Mbr_partition1.Part_type == 'e' && mbr.Mbr_partition1.Part_status == 'o' {
		isExtended = true
		startExtendedPartition = int(mbr.Mbr_partition1.Part_start)
	}

	namePartition = list.ReturnNameString(mbr.Mbr_partition2.Part_name)
	if namePartition == name {
		startPartition = int(mbr.Mbr_partition2.Part_start)
	} else if mbr.Mbr_partition2.Part_type == 'e' && mbr.Mbr_partition2.Part_status == 'o' {
		isExtended = true
		startExtendedPartition = int(mbr.Mbr_partition2.Part_start)
	}

	namePartition = list.ReturnNameString(mbr.Mbr_partition3.Part_name)
	if namePartition == name {
		startPartition = int(mbr.Mbr_partition3.Part_start)
	} else if mbr.Mbr_partition3.Part_type == 'e' && mbr.Mbr_partition3.Part_status == 'o' {
		isExtended = true
		startExtendedPartition = int(mbr.Mbr_partition3.Part_start)
	}

	namePartition = list.ReturnNameString(mbr.Mbr_partition4.Part_name)
	if namePartition == name {
		startPartition = int(mbr.Mbr_partition4.Part_start)
	} else if mbr.Mbr_partition4.Part_type == 'e' && mbr.Mbr_partition4.Part_status == 'o' {
		isExtended = true
		startExtendedPartition = int(mbr.Mbr_partition4.Part_start)
	}

	if isExtended {
		file.Seek(int64(startExtendedPartition), 0)
		var ebr EBR
		err = binary.Read(file, binary.LittleEndian, &ebr)
		if err != nil {
			fmt.Println("Error al abrir el archivo")
			log.Fatal(err)
		}
		namePartition = list.ReturnNameString(ebr.Part_name)
		if namePartition == name {
			startPartition = int(ebr.Part_start)
		} else {
			for ebr.Part_size != -1 {
				file.Seek(int64(ebr.Part_next), 0)
				err = binary.Read(file, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Println("Error al abrir el archivo")
					log.Fatal(err)
				}
				namePartition = list.ReturnNameString(ebr.Part_name)
				if namePartition == name {
					startPartition = int(ebr.Part_start)
					break
				}
			}
		}
	}
	return startPartition
}

/*The function return the string for an arrays of bytes*/
func (list *PartitionList) ReturnNameString(name [16]byte) string {
	tmp := ""
	var c byte
	for i := 0; i < 16; i++ {
		if name[i] == c {
			break
		}
		tmp += string(name[i])
	}
	return tmp
}

/*The function return the value withouth marks*/
func (tmp *PartitionList) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}

/*This function show every nodes in the list*/
func (list *PartitionList) ShowListPartition() {
	tmp := list.RootNode
	for i := 0; i < list.Size; i++ {
		fmt.Println(tmp.Name, tmp.PathFile)
		tmp = tmp.NextNode
	}
}
