package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type addEnrollDateRequest struct {
	Semester  string    `json:"semester" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

func (server *Server) addEnrollDate(c *gin.Context) {
	var req addEnrollDateRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Error al bindear la request: "+err.Error()))
		return
	}

	result, err := server.store.Run("MATCH (sem:Semester{semester:$semester}) CREATE (enrdate:EnrollDate{start_time:datetime($start_time),end_time:datetime($end_time)})-[rel:Belongs]->(sem)", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)

}

type AddStudentToDate struct {
	IDStudent int       `json:"id_student" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

func (server *Server) addStudentToDate(c *gin.Context) {
	var req AddStudentToDate

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Error al bindear la request: "+err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (stu:Student),(enrdate:EnrollDate)
	WHERE stu.id = $id_student and (enrdate.start_time=datetime($start_time) and enrdate.end_time=datetime($end_time))
	CREATE (stu)-[daterel:EnrollsOn]->(enrdate)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)

}

type DeleteStudentToDate struct {
	IDStudent int       `json:"id_student" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

func (server *Server) deleteStudentToDate(c *gin.Context) {
	var req DeleteStudentToDate

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Error al bindear la request: "+err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (stu:Student) -[daterel:EnrollsOn]-> (enrdate:EnrollDate)
	WHERE stu.id = $id_student and (enrdate.start_time=datetime($start_time) and enrdate.end_time=datetime($end_time))
	DELETE daterel`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)

}
