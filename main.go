package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mia/proyecto2/commands"
)

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	PartitionID string `json:"idparticion"`
}

type CmdComand struct {
	Cmd string
}

func main() {
	router := mux.NewRouter()
	lexer := commands.Lexer{}
	router.Use(corsMiddleware)
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var login User
		err := json.NewDecoder(r.Body).Decode(&login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		commandLogin := "login >user=" + login.Username + " >pwd=" + login.Password + " >id=" + login.PartitionID
		fmt.Println(commandLogin)
		commandLogin = strings.TrimSpace(commandLogin)
		commandLogin = strings.ToLower(commandLogin)
		contentR := lexer.GeneralComand(commandLogin)
		res := struct{ Message string }{Message: contentR}
		json.NewEncoder(w).Encode(res)
	}).Methods("POST")

	router.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		var comM CmdComand
		var respuestas []string
		err := json.NewDecoder(r.Body).Decode(&comM)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		comandos := strings.Split(comM.Cmd, "\n")
		for i := 0; i < len(comandos); i++ {
			tmp := comandos[i]
			if tmp != "" {
				tmp = strings.TrimSpace(tmp)
				if tmp[0] != '#' {
					tmp = strings.TrimSpace(tmp)
					tmp = strings.ToLower(tmp)
					tmp2 := lexer.GeneralComand(tmp)
					respuestas = append(respuestas, tmp2+"\n")
				} else if tmp[0] == '#' {
					tmp2 := strings.Replace(tmp, "#", "---- Comentario: ", -1)
					respuestas = append(respuestas, tmp2)
				}
			}
		}
		res := struct{ Message []string }{Message: respuestas}
		json.NewEncoder(w).Encode(res)
	}).Methods("POST")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		respuestas := lexer.GeneralComand("logout")
		res := struct{ Message string }{Message: respuestas}
		json.NewEncoder(w).Encode(res)
	}).Methods("GET")

	router.HandleFunc("/individualComand", func(w http.ResponseWriter, r *http.Request) {
		var comM CmdComand
		err := json.NewDecoder(r.Body).Decode(&comM)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		comM.Cmd = strings.TrimSpace(comM.Cmd)
		comM.Cmd = strings.ToLower(comM.Cmd)
		mess := lexer.GeneralComand(comM.Cmd)
		typeC := "command"
		if matched, _ := regexp.Match("(rep)(.*)", []byte(comM.Cmd)); matched {
			typeC = "rep"
		}
		res := struct{ Message, typeC string }{Message: mess, typeC: typeC}
		json.NewEncoder(w).Encode(res)
	}).Methods("POST")

	fmt.Println("Servidor activo")
	log.Fatal(http.ListenAndServe(":8000", router))
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir solicitudes desde cualquier origen con cualquier método y encabezado
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
