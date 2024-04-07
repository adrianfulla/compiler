package main

import (
	"fmt"
	"net/http"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/gin-gonic/gin"
)

func main() {
	serve()
	// makeDirAfd("(a|b)*abb")

}

func serve() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/automata", func(c *gin.Context) {
		var request struct {
			Regex string `json:"regex"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := makeDirAfd(request.Regex)
		fmt.Print(string(response))
		c.Data(http.StatusOK, "application/json", response)
	})

	r.Run()
}

// func shuntingYard(Regex string) string {
// 	// Prueba de la función de validación
// 	postfix, err := automatas.InfixToPosfix(Regex)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return ""
// 	} else {
// 		fmt.Println("Regex postfix:", postfix)
// 		return postfix
// 	}
// }

func makeDirAfd(Regex string) []byte {
	// Prueba de la función de validación
	postfix, err := automatas.InfixToPosfix(Regex)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	} else {
		fmt.Println("Regex postfix:", postfix)
		afd := automatas.NewDirectAfd(postfix)
		jsonAfd, err := afd.MarshalJson()
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return nil
		}
		jsonArbol, err := afd.Arbol.ToJson()
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return nil
		}
		return append(jsonAfd, jsonArbol...)
	}
}
