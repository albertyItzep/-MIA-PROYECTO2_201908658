package main

import (
	"github.com/mia/proyecto2/commands"
)

func main() {
	var lexer commands.Lexer
	//lexer.GeneralComand("mkdisk >Size=3000 >unit=K >path=/home/user/Disco1.dsk >fit=ff")

	lexer.GeneralComand("fdisk >size=300 >path=/home/user/disco1.dsk >name=Particion1")
}

/*
fdisk >Size=300 >path=/home/user/Disco1.dsk >name=Particion1
var lexer commands.Lexer
	//lexer.GeneralComand("mkdisk >Size=3000 >unit=K >path=/home/user/Disco1.dsk >fit=ff")

	for {
		fmt.Println("Ingrese el comando deseado")
		reader := bufio.NewReader(os.Stdin)
		commandString, _ := reader.ReadString('\n')
		commandString = strings.TrimSpace(commandString)
		commandString = strings.ToLower(commandString)

		fmt.Println(commandString)
		if commandString == "exit" {
			fmt.Println("Nos vemos Luego")
			return
		} else if commandString != "" {
			lexer.GeneralComand(commandString)
		} else if commandString == "" {
			fmt.Println("Por Favor ingrese un comando la proxima vez")
			return
		}
	}

list := structs.SpacesList{}
	list.InsertNode(1, 4, 'f')
	list.InsertNode(9, 13, 'f')
	list.InsertNode(16, 21, 'f')
	list.InsertNode(26, 30, 'f')
	list.FillList(40)
	list.ShowList()
	min := list.MinSpace(3)
	maj := list.MajSpace(5)
	fmt.Println(min, ',', maj)
	fmt.Println(list.FirstSpace(4))
*/
