package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddGroupRequest struct {
	ID       int64  `json:"id" binding:"required"`
	Semester string `json:"semester" binding:"required"`
}

func (server *Server) addGroup(c *gin.Context) {
	var req AddGroupRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run("CREATE (est:Group{id:$id,semester:$semester})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)

}

type AddGroupToSubjectRequest struct {
	IDGroup   int64 `json:"id_group" binding:"required"`
	IDSubject int64 `json:"id_subject" binding:"required"`
}

func (server *Server) AddGroupToSubject(c *gin.Context) {
	var req AddGroupToSubjectRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run(`MATCH (subject:Subject),(group:Group)
	WHERE subject.id = $id_subject and group.id = $id_group
	RETURN (EXISTS((group)-[:Belongs]->(subject)) and EXISTS((subject)-[:Has]->(group)))`, u.StructToMap(req))
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

	result, err = server.store.Run(`MATCH (subject:Subject),(group:Group)
	WHERE subject.id = $id_subject and group.id = $id_group
	CREATE (subject)-[:Has]->(group)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

type AddGroupToSemesterRequest struct {
	SemesterPeriod string `json:"semester" binding:"required"`
	IDGroup        int64  `json:"id_group" binding:"required"`
}

func (server *Server) AddGroupToSemester(c *gin.Context) {
	var req AddGroupToSemesterRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"+err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (semester:Semester),(group:Group)
	WHERE semester.semester = $semester and group.id = $id_group
	RETURN EXISTS((group)-[:Ocurred]->(semester))`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if result.Next() {
		if result.Record().Values[0].(bool) {
			c.JSON(http.StatusBadRequest, Result{
				Error:    "La asignatura ya pertenece al periodo recibido",
				Status:   http.StatusBadRequest,
				Response: req,
			})
			return
		}
	}

	result, err = server.store.Run(`MATCH (semester:Semester),(group:Group)
	WHERE semester.semester = $semester and group.id = $id_group
	CREATE (group)-[ocr:Ocurred]->(semester)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Grupo añadido al periodo academico descrito",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})
}
