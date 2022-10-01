package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddStudentRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (server *Server) addStudent(c *gin.Context) {
	var req AddStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("Debe colocar un id valido en la petición"))
		return
	}

	result, err := server.store.Run("CREATE (est:Student{id:$id})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)

}

/*
* MATCH (student:Student),(group:Group) WHERE student.id = 1 and group.id = 1
* return student,group
?
*/
type AddGroupToStudentRequest struct {
	IDStudent int64 `json:"id_student" binding:"required"`
	IDGroup   int64 `json:"id_group" binding:"required"`
}

func (server *Server) addGroupToStudent(c *gin.Context) {
	var req AddGroupToStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (student:Student),(group:Group)
	WHERE student.id = $id_student and group.id = $id_group
	RETURN EXISTS((student)-[:Enrolls]->(group))`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if result.Next() {
		if result.Record().Values[0].(bool) {
			c.JSON(http.StatusBadRequest, errorResponse("El usuario ya esta inscrito en la asignatura"))
			return
		}
	}

	result, err = server.store.Run(`MATCH (student:Student),(group:Group)
	WHERE student.id = $id_student and group.id = $id_group
	CREATE (student)-[:Enrolls]->(group)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

type DeleteGroupToStudentRequest struct {
	IDStudent int64 `json:"id_student" binding:"required"`
	IDGroup   int64 `json:"id_group" binding:"required"`
}

func (server *Server) deleteGroupToStudent(c *gin.Context) {
	var req AddGroupToStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (student:Student)-[rel:Enrolls]->(group:Group)
		WHERE student.id = $id_student and group.id = $id_group
		DELETE r`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)

}

type AddCareerToStudentRequest struct {
	IDStudent int64 `json:"id_student" binding:"required"`
	IDCareer  int64 `json:"id_career" binding:"required"`
}

func (server *Server) AddCareerToStudent(c *gin.Context) {
	var req AddCareerToStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (student:Student),(career:Career)
	WHERE student.id = $id_student and career.id = $id_career
	RETURN EXISTS((student)-[:Study]->(career))`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if result.Next() {
		if result.Record().Values[0].(bool) {
			c.JSON(http.StatusBadRequest, errorResponse("El usuario ya esta inscrito en la asignatura"))
			return
		}
	}

	result, err = server.store.Run(`MATCH (student:Student),(career:Career)
	WHERE student.id = $id_student and career.id = $id_career
	CREATE (student)-[:Study]->(career)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)

}

type getGroupsEnrolledByStudentRequest struct {
	IDStudent int64  `uri:"id" json:"id" binding:"required"`
	Semester  string `uri:"semester" json:"semester"`
}

func (server *Server) getGroupsEnrolledByStudent(c *gin.Context) {
	var req getGroupsEnrolledByStudentRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (student:Student)-[:Enrolls]->(group:Group)
	WHERE student.id = $id
	RETURN student,collect(group) as groups`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	userRecord, err := result.Single()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	groups, _ := userRecord.Get("groups")
	c.JSON(http.StatusOK, groups)

}

type getCareerByStudentRequest struct {
	IDStudent int64 `uri:"id" json:"id" binding:"required"`
}

func (server *Server) getCareersByStudent(c *gin.Context) {
	var req getCareerByStudentRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	result, err := server.store.Run(`MATCH (student:Student)-[:Study]->(career:Career)
	WHERE student.id = $id
	RETURN student,collect(career) as careers`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	userRecord, err := result.Single()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	groups, _ := userRecord.Get("careers")
	c.JSON(http.StatusOK, groups)

}
