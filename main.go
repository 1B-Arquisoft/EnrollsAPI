package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type things struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var thingsList = []things{
	{ID: "1", Name: "Gonorrea"},
	{ID: "2", Name: "Alejandra"},
	{ID: "3", Name: "Alguien más"},
}

func getThings(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, thingsList)

}

func addThings(context *gin.Context) {
	var thing things

	err := context.BindJSON(&thing)
	if err != nil {
		return
	}

	thingsList = append(thingsList, thing)

	context.IndentedJSON(http.StatusCreated, thing)
}

func searchID(id string) (*things, error) {

	for i, v := range thingsList {
		if v.ID == id {
			return &thingsList[i], nil
		}
	}

	return nil, errors.New("No hay id en lista")
}

func getThingID(context *gin.Context) {
	id := context.Param("id")
	thingSearched, err := searchID(id)
	if err != nil {
		return
	}

	context.IndentedJSON(http.StatusAccepted, thingSearched)

}

func updateThingID(context *gin.Context) {
	id := context.Param("id")

	name := context.Query("name")

	thing, err := searchID(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, nil)
		return
	}

	thing.Name = name

	context.IndentedJSON(http.StatusOK, thing)

}

func main() {
	router := gin.Default()

	router.GET("/things", getThings)
	router.GET("/things/:id", getThingID)
	router.PATCH("/things/:id", updateThingID)
	router.POST("/things", addThings)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	router.Run("localhost:9090")
}
