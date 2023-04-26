package structs

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type DiskNode struct {
	PathDisk                        string
	DiskSize, numPart, numPartMount int
	nextNode                        *DiskNode
	previusNode                     *DiskNode
}

type DiskList struct {
	Size     int
	rootNode *DiskNode
	endNode  *DiskNode
}

/*The function insert a disk in the list*/
func (list *DiskList) InsertNode(pathFile string, diskSize int) {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	if !list.ExistDiscList(pathFile) && list.ExistFileFisic(pathFile) {
		tmp := DiskNode{PathDisk: pathFile, DiskSize: diskSize}
		list.Size++
		if list.rootNode == nil {
			list.rootNode = &tmp
			list.endNode = &tmp
		} else {
			list.endNode.nextNode = &tmp
			tmp.previusNode = list.endNode
			list.endNode = &tmp
		}
	}
}

/*The function add a partition in his disk but not in the RAM*/
func (list *DiskList) InsertPartitionDisk(pathFile string) {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	tmp := list.rootNode
	for i := 0; i < list.Size; i++ {
		if tmp.PathDisk == pathFile {
			tmp.numPart++
			break
		}
		tmp = tmp.nextNode
	}
}

/*The function add a partition in his disk and in RAM*/
func (list *DiskList) InsertPartitionDiskMounted(pathFile string) {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	tmp := list.rootNode
	for i := 0; i < list.Size; i++ {
		if tmp.PathDisk == pathFile {
			tmp.numPartMount++
			break
		}
		tmp = tmp.nextNode
	}
}

/*the function return the cant of partitions in the disk*/
func (list *DiskList) ReturnPartitionsDiskMounted(pathFile string) int {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	tmp := list.rootNode
	for i := 0; i < list.Size; i++ {
		if tmp.PathDisk == pathFile {
			return tmp.numPartMount
		}
		tmp = tmp.nextNode
	}
	return -1
}

/*The function return the size of the disk*/
func (list *DiskList) ReturSizeDisk(pathFile string) int {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	tmp := list.rootNode
	for i := 0; i < list.Size; i++ {
		if tmp.PathDisk == pathFile {
			return tmp.DiskSize
		}
		tmp = tmp.nextNode
	}
	return -1
}

/*The function return if disk exist in the list*/
func (list *DiskList) ExistDiscList(pahtFile string) bool {
	pahtFile = list.ReturnValueWithoutMarks(pahtFile)
	tmp := list.rootNode
	for i := 0; i < list.Size; i++ {
		if tmp.PathDisk == pahtFile {
			return true
		}
		tmp = tmp.nextNode
	}
	return false
}

/*The function return if file exist in fisic*/
func (list *DiskList) ExistFileFisic(pathFile string) bool {
	pathFile = list.ReturnValueWithoutMarks(pathFile)
	if _, err := os.Stat(pathFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (tmp *DiskList) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}

func (list *DiskList) ShowDisk() {
	tmp := list.rootNode
	for i := 0; i < list.Size; i++ {
		fmt.Println(i)
		fmt.Println(tmp.PathDisk + ", " + strconv.Itoa(tmp.DiskSize))
		tmp = tmp.nextNode
	}
}
