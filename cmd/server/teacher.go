package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddTeacherRequest struct {
	ID int64 `json:"id_teacher" binding:"required"`
}

func (server *Server) addTeacher(c *gin.Context) {
	var req AddTeacherRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Message:  "Debes colocar campos validos en la petición",
			Error:    err.Error(),
			Status:   http.StatusBadRequest,
			Response: req,
		})
		return
	}

	result, err := server.store.Run("CREATE (tech:Teacher{id:$id_teacher})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB"))
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Profesor creado con exito!",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})

}

type AddTeacherToGroupRequest struct {
	IDTeacher int64 `json:"id_teacher" binding:"required"`
	IDGroup   int64 `json:"id_group" binding:"required"`
}

func (server *Server) AddTeacherToGroup(c *gin.Context) {
	var req AddTeacherToGroupRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (teacher:Teacher),(group:Group)
	WHERE teacher.id = $id_teacher and group.id = $id_group
	RETURN EXISTS((teacher)-[:Taught]->(group))`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	if result.Next() {
		if result.Record().Values[0].(bool) {
			c.JSON(http.StatusBadRequest, Result{
				Error:    "La asignatura ya pertenece a la carrera" + err.Error(),
				Response: req,
			})
			return
		}
	}

	result, err = server.store.Run(`MATCH (teacher:Teacher),(group:Group)
	WHERE teacher.id = $id_teacher and group.id = $id_group
	CREATE (teacher)-[:Taught]->(group)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

type getGroupsTaughtbyTeacherRequest struct {
	IDTeacher int64  `uri:"id" json:"id" binding:"required"`
	Semester  string `uri:"semester" json:"semester" binding:"required"`
}

func (server *Server) getGroupsTaughtbyTeacher(c *gin.Context) {
	var req getGroupsTaughtbyTeacherRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (teacher:Teacher)-[:Taught]->(group:Group)
	WHERE teacher.id = $id and group.semester = $semester
	RETURN teacher,collect(group) as group`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	userRecord, err := result.Single()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	groups, _ := userRecord.Get("group")
	c.JSON(http.StatusOK, Result{
		Status:   http.StatusOK,
		Message:  "Peticion Exitosa: Grupos dictador por el profesor",
		Result:   groups,
		Response: req,
	})

}
