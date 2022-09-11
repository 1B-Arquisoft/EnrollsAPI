package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddCareerRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (server *Server) addCareer(c *gin.Context) {
	var req AddStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run("CREATE (est:Career{id:$id})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB"))
		return
	}

	c.JSON(http.StatusOK, result)

}
