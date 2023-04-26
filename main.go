package main

import (
	//"github.com/mia/proyecto2/structs"

	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mia/proyecto2/commands"
)

func main() {
	lexer := commands.Lexer{}
	//lexer.GeneralComand("fdisk >size=300 >type=p >unit=k >fit=bf >path=/home/user/disco1.dsk >name=\"particion1\"")
	//lexer.GeneralComand("fdisk >size=300 >type=e >unit=k >fit=bf >path=/home/user/disco1.dsk >name=\"particion2\"")
	//lexer.GeneralComand("fdisk >size=100 >type=l >unit=k >fit=bf >path=/home/user/disco1.dsk >name=\"particion3\"")
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
}

/*
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
*/
