package myapp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/hyson007/assignment5/message"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

const createTableString = `
CREATE TABLE IF NOT EXISTS courses (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255),
	duration INT,
	description VARCHAR(255)
)`

func New(reset bool) (*App, error) {
	a := App{}
	var err error
	a.Router = mux.NewRouter()
	a.DB, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/my-database")
	if err != nil {
		return nil, err
	}
	if reset {
		_, err = a.DB.Exec("DROP TABLE IF EXISTS courses")
		if err != nil {
			return nil, err
		}
	}
	_, err = a.DB.Exec(createTableString)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

type Course struct {
	ID          int
	Title       string
	Description string
	Duration    int
}

func (c Course) String() string {
	return fmt.Sprintf("course ID: %d, course title: %s, course duration: %d, course description: %s", c.ID, c.Title, c.Duration, c.Description)
}

// func GetAllCoursesHandler will return all courses in database
func (a *App) GetAllCoursesHandler(w http.ResponseWriter, r *http.Request) {
	courses, err := a.GetAllCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}

// func GetCourseByTitleHandler will return course based on the given title
func (a *App) GetCourseByTitleHandler(w http.ResponseWriter, r *http.Request) {
	var course Course
	var err error
	vars := mux.Vars(r)
	course, err = a.GetCourseByTitle(vars["courseid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if course.ID == 0 {
		json.NewEncoder(w).Encode(message.New("Course not found"))
		return
	}
	json.NewEncoder(w).Encode(course)
}

// func AddCourseHandler will first check the content type and http method
// as well as request body, if the necessary fields are present, then it will
// perform database insertion
func (a *App) AddCourseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {

		if r.Method == http.MethodPost {
			var course Course
			var err error

			requestBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(requestBody, &course)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if course.Title == "" || course.Duration == 0 || course.Description == "" {
				http.Error(w, "Course title, duration and description are required", http.StatusBadRequest)
				return
			}
			err = a.AddCourse(course.Title, course.Duration, course.Description)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(message.New("Course added successfully"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Content Type not allowed", http.StatusBadRequest)
	}
}

// func UpdateCourseHandler will first check the content type and http method
// as well as request body, if the necessary fields are present, then it will
// perform database update
func (a *App) UpdateCourseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {

		if r.Method == http.MethodPut {
			var course Course
			var err error

			requestBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(requestBody, &course)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if course.Title == "" || course.Duration == 0 || course.Description == "" {
				http.Error(w, "Course title, duration and description are required", http.StatusBadRequest)
				return
			}
			err = a.UpdateCourse(course.Title, course.Duration, course.Description)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(message.New("Course updated successfully"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Content Type not allowed", http.StatusBadRequest)
	}
}

// func DeleteCourseHandler will based on provided course title to perform
// database deletion
func (a *App) DeleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {

		if r.Method == http.MethodDelete {
			var err error
			vars := mux.Vars(r)

			courseTitle, ok := vars["courseid"]
			if !ok {
				http.Error(w, "Course title is required", http.StatusBadRequest)
				return
			}

			err = a.DeleteCourseByTitle(courseTitle)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(message.New("Course deleted successfully"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Content Type not allowed", http.StatusBadRequest)
	}
}

// func GetCourseByTitle will return course with given title
// it expect to return at most one course
// if no record, it will return empty course, note, this is not considered error
func (a *App) GetCourseByTitle(title string) (Course, error) {
	var course Course
	queryString := fmt.Sprintf("SELECT * FROM courses WHERE title='%s'", title)
	row := a.DB.QueryRow(queryString)
	err := row.Scan(&course.ID, &course.Title, &course.Duration, &course.Description)
	if err == sql.ErrNoRows {
		return course, nil
	}
	if err != nil {
		return course, err
	}
	return course, nil
}

// func GetAllCourses will return all courses in database
func (a *App) GetAllCourses() ([]Course, error) {
	var courses []Course
	rows, err := a.DB.Query("SELECT * FROM courses")
	if err != nil {
		return courses, err
	}
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.Title, &course.Duration, &course.Description)
		if err != nil {
			return courses, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

// func AddCourse will first check if course with same title exist, if not will perform
// database insertion
func (a *App) AddCourse(title string, duration int, description string) error {
	course, err := a.GetCourseByTitle(title)
	if err != nil {
		return err
	}
	if course.Title != "" {
		return fmt.Errorf("Course with title '%s' already exists", title)
	}

	query := fmt.Sprintf("INSERT INTO courses (title, duration, description) VALUES ('%s', %d, '%s')", title, duration, description)
	_, err = a.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// func UpdateCourse will first check if course with same title exist,
// if no, it will perform insert
// if yes, it will perform update
func (a *App) UpdateCourse(title string, duration int, description string) error {
	course, err := a.GetCourseByTitle(title)
	if err != nil {
		return err
	}
	if course.Title == "" {
		// log.Printf("Course with title '%s' does not exist, inserting\n", title)
		return a.AddCourse(title, duration, description)
	}

	query := fmt.Sprintf("UPDATE courses SET duration=%d, description='%s', title='%s' where ID=%d", duration, description, title, course.ID)
	_, err = a.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// func DeleteCourseByTitle will first check if course with same title exist,
// if no, it will return error
// if yes, it will perform delete
func (a *App) DeleteCourseByTitle(title string) error {
	course, err := a.GetCourseByTitle(title)
	if err != nil {
		return err
	}
	if course.Title == "" {
		return fmt.Errorf("Course with title '%s' does not exist", title)
	}

	return a.deleteCourse(title)
}

// func DeleteCourse will delete the course with given title
func (a *App) deleteCourse(title string) error {
	query := fmt.Sprintf("DELETE FROM courses WHERE title='%s'", title)
	_, err := a.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
