package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddSemesterRequest struct {
	Semester string `json:"semester" binding:"required"`
}

func (server *Server) addSemester(c *gin.Context) {
	var req AddSemesterRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"+err.Error()))
		return
	}

	result, err := server.store.Run("CREATE (sem:Semester{semester:$semester})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)

}
