package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MountNode struct {
	PathFile, IdPartition, NamePartition string
	StartParticion, SizeOfPartitio       int
	NextNode                             *MountNode
	PreviusNode                          *MountNode
}

type MountList struct {
	Size     int
	RootNode *MountNode
	EndNode  *MountNode
}

/*the function moun partition in RAM*/
func (list *MountList) InserMount(pathFile string, namePartition string, startpartition int, sizeOfPartition int, numberOfPartitionMounted int, indexDisk int) {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	namePartition = list.ReturnValueWithoutMarks(namePartition)
	if !list.ExistPartitionList(pathFile, namePartition) {
		id := "58"
		id += strconv.Itoa(indexDisk)
		id += list.ReturnLetterAsigned(numberOfPartitionMounted)
		tmp := MountNode{PathFile: pathFile, NamePartition: namePartition, IdPartition: id, StartParticion: startpartition, SizeOfPartitio: sizeOfPartition}
		if list.RootNode == nil {
			list.RootNode = &tmp
			list.EndNode = &tmp
			list.Size++
		} else {
			list.EndNode.NextNode = &tmp
			tmp.PreviusNode = list.EndNode
			list.EndNode = &tmp
			list.Size++
		}
	} else {
		fmt.Println("se encuentra en memoria la particion indicada")
	}
}

/*The metod return of the letter asigned of mounted*/
func (list *MountList) ReturnLetterAsigned(numberPartition int) string {
	letterArray := [27]string{"aa", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	return letterArray[numberPartition]
}

/*The funtion return if partition existed in the list*/
func (list *MountList) ExistPartitionList(pathFile string, namePartition string) bool {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	namePartition = list.ReturnValueWithoutMarks(namePartition)
	tmp := list.RootNode
	for i := 0; i < list.Size; i++ {
		if tmp.PathFile == pathFile && tmp.NamePartition == namePartition {
			return true
		}
		tmp = tmp.NextNode
	}
	return false
}

/*The function show partitions*/
func (list *MountList) ShowPartition() {
	tmp := list.RootNode
	for i := 0; i < list.Size; i++ {
		fmt.Println(tmp.NamePartition, tmp.IdPartition)
		tmp = tmp.NextNode
	}

}

/*The function return a string without marks*/
func (tmp *MountList) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
