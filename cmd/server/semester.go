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
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run("CREATE (sem:Semester{semester:$semester})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Semestre creado con exito!",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})

}
