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
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run("MATCH (sem:Semester{semester:$semester}) CREATE (enrdate:EnrollDate{start_time:datetime($start_time),end_time:datetime($end_time)})-[rel:Belongs]->(sem)", u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Cita de inscripción creada con exito",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})

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
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (stu:Student),(enrdate:EnrollDate)
	WHERE stu.id = $id_student and (enrdate.start_time=datetime($start_time) and enrdate.end_time=datetime($end_time))
	CREATE (stu)-[daterel:EnrollsOn]->(enrdate)`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Message:  "Estudiante añadido a cita de inscripción con exito!",
		Status:   http.StatusOK,
		Response: req,
		Result:   result,
	})

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
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (stu:Student) -[daterel:EnrollsOn]-> (enrdate:EnrollDate)
	WHERE stu.id = $id_student and (enrdate.start_time=datetime($start_time) and enrdate.end_time=datetime($end_time))
	DELETE daterel`, u.StructToMap(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Status:   http.StatusOK,
		Message:  "Desasignada cita de inscripción a estudiante",
		Result:   result,
		Response: req,
	})

}

type GetDatesByStudent struct {
	IDStudent int `uri:"id" json:"id" binding:"required"`
}

type GetDatesByStudentResult struct {
	IDStudent int         `json:"id_student"`
	Dates     interface{} `json:"dates"`
}

func (server *Server) getDatesByStudent(c *gin.Context) {
	var req GetDatesByStudent

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Error:    err.Error(),
			Response: req,
		})
		return
	}

	result, err := server.store.Run(`MATCH (stu:Student) -[daterel:EnrollsOn]-> (enrdate:EnrollDate)
	WHERE stu.id = $id
	RETURN stu,collect(enrdate) as dates`, u.StructToMap(req))
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

	dates, _ := records.Get("dates")

	c.JSON(http.StatusOK, Result{
		Status:  http.StatusOK,
		Message: "Lista de fechas de inscripción del estudiante",
		Result: GetDatesByStudentResult{
			IDStudent: req.IDStudent,
			Dates:     dates,
		},
		Response: req,
	})

}
