package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddSemesterRequest struct {
	Year string `json:"year" binding:"required"`
}

func (server *Server) addSemester(c *gin.Context) {
	var req AddSemesterRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run("CREATE (est:Semester{year:$year})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB"))
		return
	}

	c.JSON(http.StatusOK, result)

}
