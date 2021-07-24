package main

import (
	"dutrozkladapi/header"
	"dutrozkladapi/router"
	"dutrozkladapi/scrapper/students"
	"dutrozkladapi/util"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command flags
	bind := flag.String("bind", "", "Binds to an IP Address (default \"0.0.0.0\")")
	port := flag.String("port", "9090", "Listens on a specific port")
	update := flag.Bool("update", false, "Forces listing update on start")
	debug := flag.Bool("debug", false, "Enables debug mode, shows more verbosity")
	cors := flag.String("acao", "*", "Sets CORS ACAO header")
	flag.Parse()

	log.Println("DUT Rozklad API init.")

	// Scrap faculties data if it doesn't exist otherwise load the file
	if _, err := os.Stat("data/faculties.json"); os.IsNotExist(err) || *update {
		students.UpdateFaculties()
	} else {
		util.LoadFaculties()
	}

	// // Scrap teachers data if it doesn't exist otherwise load the file
	// if _, err := os.Stat("data/teachers.json"); os.IsNotExist(err) || *update {
	// 	scrapper.UpdateTeachers()
	// } else {
	// 	util.LoadTeachers()
	// }

	log.Println("[STATS] Kafedras count: " + fmt.Sprint(len(header.Kafedras)))
	log.Println("[STATS] Teachers count: " + fmt.Sprint(util.CountTotalTeachers()))
	log.Println("[STATS] Faculties count: " + fmt.Sprint(len(header.Faculties)))
	log.Println("[STATS] Groups count: " + fmt.Sprint(util.CountTotalGroups()))

	if !*debug {
		// Disable debug mode
		gin.SetMode(gin.ReleaseMode)
	}

	// Create new router
	d := gin.Default()

	// Enable panic recovery
	d.Use(gin.Recovery())

	// Allow access origin
	d.Use(util.AllowOrigin(*cors))

	// Create router groups
	api := d.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/faculties", router.GetFaculties)
			v1.GET("/courses/:faculty", router.GetCourses)
			v1.GET("/groups/:faculty/:course", router.GetGroups)
			// v1.GET("/kafedras", router.GetKafedras)
			// v1.GET("/teachers", router.GetTeachers)

			v1.GET("/timetable/:group/:startdate/:enddate", router.GetTimeTable)
			// v1.POST("/teachertimetable", router.GetTeacherTimeTable)

			v1.GET("/stats", util.Stats)
		}
	}

	// Run server instance
	host := *bind + ":" + *port
	log.Println("Listening on " + host)
	d.Run(host)
}
