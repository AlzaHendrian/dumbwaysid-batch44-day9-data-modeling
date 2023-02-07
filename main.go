package main

import (
	"context"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"personal-web/connection"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title":   "Fauzan | Personal Web",
	"IsLogin": true,
}

type Project struct {
	Id             int
	ProjectName    string
	StartDate      time.Time
	EndDate        time.Time
	DurationText   string
	Description    string
	Technologies   []string
	ImageDirectory string
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	route.HandleFunc("/", Home).Methods("GET")
	route.HandleFunc("/contact-me", ContactMe).Methods("GET")
	route.HandleFunc("/add-project", AddProject).Methods("GET")
	route.HandleFunc("/projects/{id}", ProjectDetails).Methods("GET")
	// route.HandleFunc("/add-new-project", AddNewProject).Methods("POST")
	// route.HandleFunc("/delete-project/{id}", DeleteProject).Methods("GET")
	// route.HandleFunc("/update-project-page/{id}", UpdateProjectPage).Methods("GET")
	// route.HandleFunc("/update-project/{id}", UpdateProject).Methods("POST")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT * FROM tb_projects ORDER BY end_date DESC")

	var Projects []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.ImageDirectory)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.DurationText = CalculateDuration(each.StartDate, each.EndDate)

		Projects = append(Projects, each)
	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": Projects,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func ContactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func ProjectDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/projects/project-details.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetails := Project{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_projects WHERE id=$1", id).Scan(
		&ProjectDetails.Id, &ProjectDetails.ProjectName, &ProjectDetails.StartDate, &ProjectDetails.EndDate, &ProjectDetails.Description, &ProjectDetails.Technologies, &ProjectDetails.ImageDirectory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetails.DurationText = CalculateDuration(ProjectDetails.StartDate, ProjectDetails.EndDate)

	resp := map[string]interface{}{
		"Data":           Data,
		"ProjectDetails": ProjectDetails,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

// func AddNewProject(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ProjectName := r.PostForm.Get("project-name")
// 	StartDate := r.PostForm.Get("start-date")
// 	EndDate := r.PostForm.Get("end-date")
// 	Description := r.PostForm.Get("description")
// 	Technology1String := r.PostForm.Get("technology-1")
// 	Technology2String := r.PostForm.Get("technology-2")
// 	Technology3String := r.PostForm.Get("technology-3")
// 	Technology4String := r.PostForm.Get("technology-4")

// 	DurationText := CalculateDuration(StartDate, EndDate)
// 	Technology1, Technology2, Technology3, Technology4 := ConvertTechnologyToBoolean(
// 		Technology1String,
// 		Technology2String,
// 		Technology3String,
// 		Technology4String)

// 	var newProject = Project{
// 		ProjectName:           ProjectName,
// 		StartDate:             StartDate,
// 		EndDate:               EndDate,
// 		DurationText:          DurationText,
// 		Description:           Description,
// 		Technology1:           Technology1,
// 		Technology2:           Technology2,
// 		Technology3:           Technology3,
// 		Technology4:           Technology4,
// 		ImagePreviewDirectory: "/public/images/mobile-app-1.jpg",
// 		ImageDirectory:        "/public/images/mobile-app-1-large.jpg",
// 	}

// 	Projects = append(Projects, newProject)

// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
// }

// func DeleteProject(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	id, _ := strconv.Atoi(mux.Vars(r)["id"])

// 	Projects = append(Projects[:id], Projects[id+1:]...)

// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
// }

// func UpdateProjectPage(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	id, _ := strconv.Atoi(mux.Vars(r)["id"])

// 	var tmpl, err = template.ParseFiles("views/update-project.html")
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("message : " + err.Error()))
// 		return
// 	}

// 	ProjectDetails := Project{}

// 	for i, project := range Projects {
// 		if i == id {
// 			ProjectDetails = Project{
// 				ProjectName:           project.ProjectName,
// 				StartDate:             project.StartDate,
// 				EndDate:               project.EndDate,
// 				DurationText:          project.DurationText,
// 				Description:           project.Description,
// 				Technology1:           project.Technology1,
// 				Technology2:           project.Technology2,
// 				Technology3:           project.Technology3,
// 				Technology4:           project.Technology4,
// 				ImagePreviewDirectory: project.ImagePreviewDirectory,
// 				ImageDirectory:        project.ImageDirectory,
// 			}
// 		}
// 	}

// 	respData := map[string]interface{}{
// 		"Data":           Data,
// 		"Id":             id,
// 		"ProjectDetails": ProjectDetails,
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	tmpl.Execute(w, respData)
// }

// func UpdateProject(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	id, _ := strconv.Atoi(mux.Vars(r)["id"])

// 	Projects = append(Projects[:id], Projects[id+1:]...)

// 	err := r.ParseForm()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ProjectName := r.PostForm.Get("project-name")
// 	StartDate := r.PostForm.Get("start-date")
// 	EndDate := r.PostForm.Get("end-date")
// 	Description := r.PostForm.Get("description")
// 	Technology1String := r.PostForm.Get("technology-1")
// 	Technology2String := r.PostForm.Get("technology-2")
// 	Technology3String := r.PostForm.Get("technology-3")
// 	Technology4String := r.PostForm.Get("technology-4")

// 	DurationText := CalculateDuration(StartDate, EndDate)
// 	Technology1, Technology2, Technology3, Technology4 := ConvertTechnologyToBoolean(
// 		Technology1String,
// 		Technology2String,
// 		Technology3String,
// 		Technology4String)

// 	var newProject = Project{
// 		ProjectName:           ProjectName,
// 		StartDate:             StartDate,
// 		EndDate:               EndDate,
// 		DurationText:          DurationText,
// 		Description:           Description,
// 		Technology1:           Technology1,
// 		Technology2:           Technology2,
// 		Technology3:           Technology3,
// 		Technology4:           Technology4,
// 		ImagePreviewDirectory: "/public/images/mobile-app-1.jpg",
// 		ImageDirectory:        "/public/images/mobile-app-1-large.jpg",
// 	}

// 	Projects = append(Projects, newProject)

// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
// }

func CalculateDuration(StartDate time.Time, EndDate time.Time) string {
	// StartDateFormated, _ := time.Parse("2006-01-02", StartDate)
	// EndtDateFormated, _ := time.Parse("2006-01-02", EndDate)
	// Duration := EndtDateFormated.Sub(StartDateFormated)
	Duration := EndDate.Sub(StartDate)
	DurationHours := Duration.Hours()
	DurationDays := math.Floor(DurationHours / 24)
	DurationWeeks := math.Floor(DurationDays / 7)
	DurationMonths := math.Floor(DurationDays / 30)
	var DurationText string
	if DurationMonths > 1 {
		DurationText = strconv.FormatFloat(DurationMonths, 'f', 0, 64) + " months"
	} else if DurationMonths > 0 {
		DurationText = strconv.FormatFloat(DurationMonths, 'f', 0, 64) + " month"
	} else {
		if DurationWeeks > 1 {
			DurationText = strconv.FormatFloat(DurationWeeks, 'f', 0, 64) + " weeks"
		} else if DurationWeeks > 0 {
			DurationText = strconv.FormatFloat(DurationWeeks, 'f', 0, 64) + " week"
		} else {
			if DurationDays > 1 {
				DurationText = strconv.FormatFloat(DurationDays, 'f', 0, 64) + " days"
			} else if DurationDays > 0 {
				DurationText = strconv.FormatFloat(DurationDays, 'f', 0, 64) + " day"
			} else {
				DurationText = "less than a day"
			}
		}
	}
	return DurationText
}

func ConvertTechnologyToBoolean(
	Technology1String string,
	Technology2String string,
	Technology3String string,
	Technology4String string) (
	bool,
	bool,
	bool,
	bool) {
	var Technology1 bool
	if Technology1String == "on" {
		Technology1 = true
	}
	var Technology2 bool
	if Technology2String == "on" {
		Technology2 = true
	}
	var Technology3 bool
	if Technology3String == "on" {
		Technology3 = true
	}
	var Technology4 bool
	if Technology4String == "on" {
		Technology4 = true
	}
	return Technology1, Technology2, Technology3, Technology4
}
