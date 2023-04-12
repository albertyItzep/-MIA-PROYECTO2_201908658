package structs

import (
	"fmt"
	"regexp"
	"strings"
)

type NodeOrderList struct {
	Name, Path  string
	Size, Start int
	Next        *NodeOrderList
	Previus     *NodeOrderList
}
type OrderList struct {
	SizeList int
	RootNode *NodeOrderList
	EndNode  *NodeOrderList
}

func (list *OrderList) InsertNode(path string, name string, sizeP int, diskExists bool) {
	path = list.ReturnValueWithoutMarks(path)
	name = list.ReturnValueWithoutMarks(name)

	if diskExists {
		startPartion := list.ReturnStartPosition(path, name)
		if startPartion > 0 {
			newNodo := NodeOrderList{Name: name, Path: path, Size: sizeP}
			if list.RootNode == nil {
				list.RootNode = &newNodo
				list.EndNode = &newNodo
				list.SizeList++
			} else {
				list.EndNode.Next = &newNodo
				newNodo.Previus = list.EndNode
				list.EndNode = &newNodo
				list.SizeList++
			}
		} else {
			fmt.Println("La particion que desea montar no existe en el disco")
		}
	} else {
		fmt.Println("El disco no existe fisicamente para agregar la particion")
	}
}

/*Print the values in the list*/
func (list *OrderList) ShowList() {
	tmpNode := list.RootNode
	for i := 0; i < list.SizeList; i++ {
		fmt.Println("Name: " + tmpNode.Name + " Path: " + tmpNode.Path)
		tmpNode = tmpNode.Next
	}
}

/*Return the start of the RootNode*/
func (list *OrderList) ReturnStartPosition(path, name string) int {
	// buscamos la particion en el disco para ver que por lo menos exista
	return 137
}

/*Return value without marks*/
func (tmp *OrderList) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	remplaceString := regexp.MustCompile("\"")
	tmpString = remplaceString.ReplaceAllString(value, "")
	tmpString = strings.TrimSpace(tmpString)
	return tmpString
}
