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
	router.PUT("/Subjects")
	//? Endpoint | Agregar nodo grupo -> {id_grupo}
	router.PUT("/Groups", server.addGroup)
	//? Endpoint | Agregar nodo profesor -> {id_profesor}
	router.PUT("/Teachers", server.addTeacher)
	//? Endpoint | Agregar nodo estudiante -> {id_estudiantes}
	router.PUT("/Students", server.addStudent)
	//? Endpoint | Agregar nodo carrera -> {id_carrera}
	router.PUT("/Careers")
	//? Endpoint | Agregar nodo semestre -> {Semestre}
	router.PUT("/Semesters")
	//*--------------------------------------------*

	//*-------Inscripcion de asignaturas----------*
	//? Endpoint | Inscribir materia -> {id_estudiante,id_grupo}
	router.POST("/Enroll")
	//? Endpoint | Cancelar materia -> {id_estudiante,id_grupo}
	router.DELETE("/Enroll")
	//*--------------------------------------------*

	//*-------Asignacion de relaciones------------*
	//? Endpoint | Asignar profe a grupo -> {id_profesor,id_grupo}
	router.POST("/Groups/Assing/Teachers")
	//? Endpoint | Añadir grupo a materia -> {id_grupo,id_materia}
	router.POST("/Groups/Assing/Subjects")
	//? Endpoint | Vincular grupo a semestre -> {id_grupo,id_semestre}
	router.POST("/Groups/Assing/Semesters")
	//*--------------------------------------------*

	//*------Asignar relaciones a carrera--------*
	//? Endpoint | Añadir Materia a carrera -> {id_materia,id_carrera}
	router.POST("/Subjects/Assing/Carrers")
	//? Endpoint | Añadir Estudiante a carrera -> {id_materia,id_carrera}
	router.POST("/Students/Assing/Carrers")
	//*--------------------------------------------*

	//*------Obtener Información Estudiante--------*
	//? Endpoint | Obtener grupos inscritos -> {id_estudiante,opcional:semestre}
	router.GET("/Students/Groups")
	//? Endpoint | Obtener carreras de un estudiante -> {id_estudiante}
	router.GET("/Students/Careers")
	//*--------------------------------------------*

	//*------Obtener Información Profesor--------*
	//? Endpoint | Obtener grupos inscritos -> {id_estudiante,opcional:semestre}
	router.GET("/Teacher/Groups")
	//*--------------------------------------------*

	//*--------Obtener información Grupos----------*
	//? Endpoint | Obtener grupos de la materia -> {id_materia,opcional:semestre}
	router.GET("/Groups/BelongTo")
	//? Endpoint | Obtener grupos en los que dicta un profe -> {id_profe,opcional: semestre}
	router.GET("/Groups/TeachBy")
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
