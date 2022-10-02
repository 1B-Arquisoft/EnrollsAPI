package server

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Server struct {
	store  neo4j.Session
	router *gin.Engine
}

func Setup() *Server {
	server := Server{store: neo4jConnection()}
	router := gin.Default()

	//*------------Creación de Nodos--------------*
	//? Endpoint | Agregar nodo materia -> {id_materia}
	router.PUT("/Subjects", server.addSubject)
	//? Endpoint | Agregar nodo grupo -> {id_grupo,id_semestre}
	router.PUT("/Groups", server.addGroup)
	//? Endpoint | Agregar nodo profesor -> {id_profesor}
	router.PUT("/Teachers", server.addTeacher)
	//? Endpoint | Agregar nodo estudiante -> {id_estudiantes}
	router.PUT("/Students", server.addStudent)
	//? Endpoint | Agregar nodo carrera -> {id_carrera}
	router.PUT("/Careers", server.addCareer)
	//? Endpoint | Agregar nodo semestre -> {Semestre}
	router.PUT("/Semesters", server.addSemester)
	//? Endpoint | Agregar nodo fecha de inscripcion -> {Semestre,{hora_inicio,hora_fin,fecha}}
	router.PUT("/EnrollDate", server.addEnrollDate)
	//*--------------------------------------------*

	//*-------Inscripcion de asignaturas----------*
	//? Endpoint | Inscribir materia -> {id_estudiante,id_grupo}
	router.POST("/Enroll", server.addGroupToStudent)
	//? Endpoint | Cancelar materia -> {id_estudiante,id_grupo}
	router.DELETE("/Enroll", server.deleteGroupToStudent)
	//*--------------------------------------------*

	//*-------Asignacion de relaciones------------*
	//? Endpoint | Asignar profe a grupo -> {id_profesor,id_grupo}
	router.POST("/Groups/Assing/Teachers", server.AddTeacherToGroup)
	//? Endpoint | Quitar profe a grupo -> {id_profesor,id_grupo}
	router.DELETE("/Groups/Assing/Teachers", server.DeleteGroupToTeacher)
	//? Endpoint | Añadir grupo a materia -> {id_grupo,id_materia}
	router.POST("/Groups/Assing/Subjects", server.AddGroupToSubject)
	//? Endpoint | Vincular grupo a semestre -> {id_grupo,id_semestre}
	router.POST("/Groups/Assing/Semesters", server.AddGroupToSemester)

	//*-------Cita de inscripción------------*
	//? Endpoint | Añadir estudiante a cita de inscripcion -> {id_student,{start_time,end_time}}
	router.POST("/Student/Assing/EnrollDate", server.addStudentToDate)
	//? Endpoint | Quitar estudiante a cita de inscripcion -> {id_student,{start_time,end_time}}
	router.DELETE("/Student/Assing/EnrollDate", server.deleteStudentToDate)
	//*--------------------------------------------*

	//*------Asignar relaciones a carrera--------*
	//? Endpoint | Añadir Materia a carrera -> {id_materia,id_carrera}
	router.POST("/Subjects/Assing/Careers", server.AddSubjectToCareer)
	//? Endpoint | Añadir Estudiante a carrera -> {id_materia,id_carrera}
	router.POST("/Students/Assing/Careers", server.AddCareerToStudent)
	//*--------------------------------------------*

	//*------Obtener Información Estudiante--------*
	//? Endpoint | Obtener grupos inscritos -> {id_estudiante,opcional:semestre}
	router.GET("/Students/:id/Groups/:semester", server.getGroupsEnrolledByStudent)
	//? Endpoint | Obtener carreras de un estudiante -> {id_estudiante}
	router.GET("/Students/:id/Careers/", server.getCareersByStudent)
	//*--------------------------------------------*

	//*--------Obtener Información Profesor--------*
	//? Endpoint | Obtener grupos inscritos -> {id_estudiante,opcional:semestre}
	router.GET("/Teacher/:id/Groups/:semester", server.getGroupsTaughtbyTeacher)
	//*--------------------------------------------*

	//*--------Obtener información Materias----------*
	//? Endpoint | Obtener grupos de la materia -> {id_materia,opcional:semestre}
	router.GET("/Subject/:id/Groups/:semester", server.getGroupsBySubject)
	//*--------------------------------------------*

	server.router = router

	return &server
}

func errorResponse(err string) gin.H {
	return gin.H{
		"error": err,
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
