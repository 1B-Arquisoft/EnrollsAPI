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
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run("CREATE (est:Subject{id:$id})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("Error al ingresar el nodo en la DB"))
		return
	}

	c.JSON(http.StatusOK, result)

}

type AddSubjectToCareerRequest struct {
	IDSubject int64 `json:"id_subject" binding:"required"`
	IDCareer  int64 `json:"id_career" binding:"required"`
}

func (server *Server) AddSubjectToCareer(c *gin.Context) {
	var req AddSubjectToCareerRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run(`MATCH (subject:Subject),(career:Career)
	WHERE subject.id = $id_subject and career.id = $id_career
	RETURN EXISTS((subject)-[:Belongs]->(career))`, u.StructToMap(req))
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

	result, err = server.store.Run(`MATCH (subject:Subject),(career:Career)
	WHERE subject.id = $id_subject and career.id = $id_career
	CREATE (subject)-[:Belongs]->(career)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
