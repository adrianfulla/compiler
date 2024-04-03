package main

import (
	"fmt"
	"net/http"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/gin-gonic/gin"
)

func main() {
	serve()
	// shuntingYard("b((a)?)")
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

		c.Data(http.StatusOK, "application/json", makeArbol(request.Regex))
	})

	r.Run()
}

func shuntingYard(Regex string) string {
	// Prueba de la función de validación
	postfix, err := automatas.InfixToPosfix(Regex)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	} else {
		fmt.Println("Regex postfix:", postfix)
		return postfix
	}
}

func makeArbol(Regex string) []byte {
	postfix := shuntingYard(Regex)
	arbol := &automatas.ArbolExpresion{}
	arbol.ConstruirArbol(postfix)
	jsonData, err := arbol.ToJson()
	if err != nil {
		return nil
	} else {
		return jsonData
	}
}
