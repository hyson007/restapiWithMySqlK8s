package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hyson007/assignment5/myapp"
)

func main() {
	var reset bool = true
	time.Sleep(time.Second * 10)
	app, err := myapp.New(reset)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer app.DB.Close()

	err = app.AddCourse("Go", 1000, "Go is awesome")
	if err != nil {
		log.Println("insert error: ", err)
	}

	err = app.AddCourse("python", 2000, "python is awesome")
	if err != nil {
		log.Println("insert error: ", err)
	}

	err = app.AddCourse("java", 3000, "java is awesome")
	if err != nil {
		log.Println("insert error: ", err)
	}

	app.Router.HandleFunc("/api/v1/courses", app.GetAllCoursesHandler).Methods("GET")
	app.Router.HandleFunc("/api/v1/course/{courseid}", app.GetCourseByTitleHandler).Methods("GET")
	app.Router.HandleFunc("/api/v1/course", app.AddCourseHandler).Methods("POST")
	app.Router.HandleFunc("/api/v1/course", app.UpdateCourseHandler).Methods("PUT")
	app.Router.HandleFunc("/api/v1/course/{courseid}", app.DeleteCourseHandler).Methods("DELETE")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", app.Router))
}
