package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddSubjectRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (server *Server) addSubject(c *gin.Context) {
	var req AddStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run("CREATE (est:Subject{id:$id})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Result:   result,
		Message:  "Petición realizada con exito: Matería añadida.",
		Response: req,
	})

}

type AddSubjectToCareerRequest struct {
	IDSubject int64 `json:"id_subject" binding:"required"`
	IDCareer  int64 `json:"id_career" binding:"required"`
}

func (server *Server) AddSubjectToCareer(c *gin.Context) {
	var req AddSubjectToCareerRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (subject:Subject),(career:Career)
	WHERE subject.id = $id_subject and career.id = $id_career
	RETURN EXISTS((subject)-[:Belongs]->(career))`, u.StructToMap(req))
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
				Error:    err.Error(),
				Response: req,
			})
			return
		}
	}

	result, err = server.store.Run(`MATCH (subject:Subject),(career:Career)
	WHERE subject.id = $id_subject and career.id = $id_career
	CREATE (subject)-[:Belongs]->(career)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Petición Realizada con exito: Añadido la materia a carrera",
		Result:   result,
		Response: req,
	})
}

type getGroupsBySubjectRequest struct {
	IDGroup  int64  `uri:"id" json:"id" binding:"required"`
	Semester string `uri:"semester" json:"semester"`
}

func (server *Server) getGroupsBySubject(c *gin.Context) {
	var req getGroupsBySubjectRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (subject:Subject)-[:Has]->(group:Group)
	WHERE subject.id = $id 
	RETURN subject,collect(group) as group`, u.StructToMap(req))
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
		Message:  "Petición realizada con exito, obtención de grupos por materias",
		Result:   groups,
		Response: req,
	})

}
