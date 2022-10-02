package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddStudentRequest struct {
	ID int64 `json:"id" binding:"required"`
}

type Result struct {
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
	Message  string      `json:"message"`
	Result   interface{} `json:"result"`
	Status   int         `json:"status"`
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
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Result:   result,
		Message:  "Estudiante creado con exito",
		Response: req,
		Status:   http.StatusOK,
	})

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
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (student:Student),(group:Group)
	WHERE student.id = $id_student and group.id = $id_group
	RETURN EXISTS((student)-[:Enrolls]->(group))`, u.StructToMap(req))
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
				Error:    "El usuario ya está inscrito en la asignatuura" + err.Error(),
				Response: req,
			})
			return
		}
	}

	result, err = server.store.Run(`MATCH (student:Student),(group:Group)
	WHERE student.id = $id_student and group.id = $id_group
	CREATE (student)-[:Enrolls]->(group)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Result:   result,
		Message:  "Estudiante inscrito a la asignatura con exito",
		Response: req,
		Status:   http.StatusOK,
	})
}

type DeleteGroupToStudentRequest struct {
	IDStudent int64 `json:"id_student" binding:"required"`
	IDGroup   int64 `json:"id_group" binding:"required"`
}

func (server *Server) deleteGroupToStudent(c *gin.Context) {
	var req AddGroupToStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (student:Student)-[rel:Enrolls]->(group:Group)
		WHERE student.id = $id_student and group.id = $id_group
		DELETE rel`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Result:   result,
		Message:  "Cancelación de estudiante a grupo exitosa",
		Status:   http.StatusOK,
		Response: req,
	})

}

type AddCareerToStudentRequest struct {
	IDStudent int64 `json:"id_student" binding:"required"`
	IDCareer  int64 `json:"id_career" binding:"required"`
}

func (server *Server) AddCareerToStudent(c *gin.Context) {
	var req AddCareerToStudentRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (student:Student),(career:Career)
	WHERE student.id = $id_student and career.id = $id_career
	RETURN EXISTS((student)-[:Study]->(career))`, u.StructToMap(req))
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
				Error:    "El estudiante ya cursa la asignatura" + err.Error(),
				Response: req,
			})
			return
		}
	}

	result, err = server.store.Run(`MATCH (student:Student),(career:Career)
	WHERE student.id = $id_student and career.id = $id_career
	CREATE (student)-[:Study]->(career)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Result:   result,
		Message:  "Estudiante inscrito a carrera exitoso",
		Response: req,
		Status:   http.StatusOK,
	})

}

type getGroupsEnrolledByStudentRequest struct {
	IDStudent int64  `uri:"id" json:"id" binding:"required"`
	Semester  string `uri:"semester" json:"semester" binding:"required"`
}

func (server *Server) getGroupsEnrolledByStudent(c *gin.Context) {
	var req getGroupsEnrolledByStudentRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}
	queryStr := `MATCH (student:Student)-[:Enrolls]->(group:Group)
	WHERE student.id = $id and group.semester = $semester
	RETURN student,collect(group) as groups`

	result, err := server.store.Run(queryStr, u.StructToMap(req))
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

	groups, _ := userRecord.Get("groups")
	c.JSON(http.StatusOK, Result{
		Result:   groups,
		Message:  "Grupos en los cuales está inscrito el estudiante.",
		Response: req,
		Status:   http.StatusOK,
	})

}

type getCareerByStudentRequest struct {
	IDStudent int64 `uri:"id" json:"id" binding:"required"`
}

func (server *Server) getCareersByStudent(c *gin.Context) {
	var req getCareerByStudentRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (student:Student)-[:Study]->(career:Career)
	WHERE student.id = $id
	RETURN student,collect(career) as careers`, u.StructToMap(req))
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

	groups, _ := userRecord.Get("careers")
	c.JSON(http.StatusOK, Result{
		Result:   groups,
		Message:  "Carrearas que curso un estudiante",
		Response: req,
		Status:   http.StatusOK,
	})

}
