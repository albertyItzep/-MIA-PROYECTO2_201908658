package structs

import (
	"fmt"
)

type NodeSpaces struct {
	Inicio, Fin, Tamano int
	Status              byte
	Next                *NodeSpaces
	Previus             *NodeSpaces
}
type SpacesList struct {
	Size, SpaceFill int
	Root            *NodeSpaces
	End             *NodeSpaces
}

func (list *SpacesList) InsertNode(inicio int, fin int, status byte) {
	newNodo := NodeSpaces{Inicio: inicio, Fin: fin, Tamano: fin - inicio, Status: status}
	list.SpaceFill += newNodo.Tamano
	list.Size++
	if list.Root == nil {
		list.Root = &newNodo
		list.End = &newNodo
	} else {
		list.End.Next = &newNodo
		newNodo.Previus = list.End
		list.End = &newNodo
	}
}

func (list *SpacesList) InsertForSize(Inicio int, fin int, tamano int, status byte) {
	newNode := NodeSpaces{Inicio: Inicio, Fin: fin, Status: status, Tamano: tamano}
	list.SpaceFill += newNode.Tamano
	list.Size++
	if list.Root == nil {
		list.Root = &newNode
		list.End = &newNode
	} else if tamano < list.Root.Tamano {
		newNode.Next = list.Root
		list.Root.Previus = &newNode
		list.Root = &newNode
	} else if tamano > list.End.Tamano {
		list.End.Next = &newNode
		newNode.Previus = list.End
		list.End = &newNode
	} else {
		tmp2 := list.Root
		x := 0
		for x < list.Size && tmp2 != nil {
			if tamano < tmp2.Tamano {
				tmp2.Previus.Next = &newNode
				newNode.Previus = tmp2.Previus
				newNode.Next = tmp2
				tmp2.Previus = &newNode
				break
			}
			x++
			tmp2 = tmp2.Next
		}
	}
}

func (list *SpacesList) ReturnOcupedSpace() int {
	return list.SpaceFill
}

func (list *SpacesList) FillList(tamanoTotal int) {
	tmp := list.Root
	for tmp.Next != nil {
		distance := tmp.Next.Inicio - tmp.Fin
		if distance >= 3 {
			startP := tmp.Fin + 1
			endP := tmp.Next.Inicio - 1
			tmp2 := NodeSpaces{Inicio: startP, Fin: endP, Status: 'f', Tamano: endP - startP}
			tmp2.Previus = tmp
			tmp2.Next = tmp.Next
			tmp.Next.Previus = &tmp2
			tmp.Next = &tmp2
			list.Size++
		}
		tmp = tmp.Next
	}

	if tamanoTotal-list.End.Fin > 3 {
		startP := list.End.Fin + 1
		tmp2 := NodeSpaces{Inicio: startP, Fin: tamanoTotal, Status: 'f', Tamano: tamanoTotal - startP}
		list.End.Next = &tmp2
		tmp2.Previus = list.End
		list.End = &tmp2
		list.Size++
	}
}

func (list *SpacesList) ClearList() {
	list.Root = nil
	list.End = nil
	list.Size = 0
}

func (list *SpacesList) FirstSpace(sizeP int) int {
	tmp := list.Root
	for i := 0; i < list.Size; i++ {
		if tmp.Status == 'f' && tmp.Tamano >= sizeP {
			return tmp.Inicio
		}
		tmp = tmp.Next
	}
	return -1
}

func (list *SpacesList) NextSpace(startP int) int {
	tmp := list.Root
	for i := 0; i < list.Size; i++ {
		if tmp.Inicio == startP {
			if tmp.Next != nil {
				return tmp.Next.Inicio
			} else {
				return -1
			}
		}
		tmp = tmp.Next
	}
	return -1
}

func (list *SpacesList) PreviusSpace(startP int) int {
	tmp := list.Root
	for i := 0; i < list.Size; i++ {
		if tmp.Inicio == startP {
			if tmp.Previus != nil {
				return tmp.Previus.Inicio
			} else {
				return -1
			}
		}
		tmp = tmp.Next
	}
	return -1
}

func (list *SpacesList) MinSpace(Tamano int) int {
	tmpList := SpacesList{}
	tmp := list.Root
	for i := 0; i < list.Size; i++ {
		if tmp.Status == 'f' {
			tmpList.InsertForSize(tmp.Inicio, tmp.Fin, tmp.Tamano, tmp.Status)
		}
		tmp = tmp.Next
	}
	tmp = tmpList.Root
	for i := 0; i < list.Size; i++ {
		if tmp.Tamano >= Tamano {
			return tmp.Inicio
		}
		tmp = tmp.Next
	}
	return -1
}

func (list *SpacesList) MajSpace(Tamano int) int {
	tmpList := SpacesList{}
	tmp := list.Root
	for i := 0; i < list.Size; i++ {
		if tmp.Status == 'f' {
			tmpList.InsertForSize(tmp.Inicio, tmp.Fin, tmp.Tamano, tmp.Status)
		}
		tmp = tmp.Next
	}
	return tmpList.End.Inicio
}

func (list *SpacesList) ExistSpace(Tamano int) bool {
	tmp := list.Root
	for i := 0; i < list.Size; i++ {
		if tmp.Tamano >= Tamano {
			return true
		}
		tmp = tmp.Next
	}
	return false
}

func (list *SpacesList) ShowList() {
	tmp := list.Root
	for i := 0; i < list.Size; i++ {
		fmt.Println(tmp.Inicio, ',', tmp.Fin)
		tmp = tmp.Next
	}
}
