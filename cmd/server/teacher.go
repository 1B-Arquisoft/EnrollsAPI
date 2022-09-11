package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddTeacherRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (server *Server) addTeacher(c *gin.Context) {
	var req AddTeacherRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run("CREATE (Teacher{id:$id})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB"))
		return
	}

	c.JSON(http.StatusOK, result)

}

type AddTeacherToGroupRequest struct {
	SemesterPeriod int64 `json:"semester_period" binding:"required"`
	IDGroup        int64 `json:"id_subject" binding:"required"`
}

func (server *Server) AddTeacherToGroup(c *gin.Context) {
	var req AddTeacherToGroupRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run(`MATCH (semester:Semester),(group:Group)
	WHERE semester.year = $semester_period and group.id = $id_group
	RETURN EXISTS((group)-[:Taught]->(semester))`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if result.Next() {
		if result.Record().Values[0].(bool) {
			c.JSON(http.StatusBadRequest, errorResponse("La asignatura ya pertenece a la carrera"))
			return
		}
	}

	result, err = server.store.Run(`MATCH (semester:Semester),(group:Group)
	WHERE semester.year = $semester_period and group.id = $id_group
	CREATE (group)-[:Taught]->(semester)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
