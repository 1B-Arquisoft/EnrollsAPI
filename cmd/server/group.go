package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddGroupRequest struct {
	ID       int64  `json:"id_group" binding:"required"`
	Semester string `json:"semester" binding:"required"`
}

func (server *Server) addGroup(c *gin.Context) {
	var req AddGroupRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run("CREATE (est:Group{id:$id_group,semester:$semester})", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Grupo creado con exito!",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})

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
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
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
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
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
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (semester:Semester),(group:Group)
	WHERE semester.semester = $semester and group.id = $id_group
	RETURN EXISTS((group)-[:Ocurred]->(semester))`, u.StructToMap(req))
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
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Grupo añadido al periodo academico descrito",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})
}

type GetStudentsEnrolledInGroupRequest struct {
	IDGroup int64 `uri:"id" json:"id_group" binding:"required"`
}

func (server *Server) getGroupInfo(c *gin.Context) {
	var req GetStudentsEnrolledInGroupRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (group:Group),(student:Student),(teacher:Teacher)
	WHERE group.id = $id_group and (student)-[:Enrolls]->(group) and (teacher)-[:Taught]->(group)
	RETURN group,collect(student) as students, teacher as teacher`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	records, err := result.Single()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	students, _ := records.Get("students")
	teacher, _ := records.Get("teacher")

	c.JSON(http.StatusOK, Result{
		Message:  "Grupo añadido al periodo academico descrito",
		Status:   http.StatusOK,
		Response: req,
		Result: GetGroupInfoResult{
			IDGroup:  req.IDGroup,
			Students: students,
			Teacher:  teacher,
		},
	})
}

type GetGroupInfoResult struct {
	IDGroup  int64       `json:"id_group"`
	Students interface{} `json:"id_students"`
	Teacher  interface{} `json:"id_teacher"`
}

type getTeacherInGroupRequest struct {
	IDGroup int64 `uri:"id" binding:"required"`
}

func (server *Server) getTeacherInGroup(c *gin.Context) {
	var req getTeacherInGroupRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (group:Group) (teacher:Teacher)
	WHERE group.id = $id_group
	RETURN (teacher)-[:Taught]->(group)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Grupo añadido al periodo academico descrito",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})
}
