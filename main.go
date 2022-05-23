package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "127.0.0.1"
	DB_PORT     = 5432
	DB_USER     = "tecmint"
	DB_PASSWORD = "root"
	DB_NAME     = "todo"
)

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Board struct {
	BoardID   int    `json:"boardid"`
	BoardName string `json:"boardname"`
}

type Task struct {
	TaskID   int    `json:"taskid"`
	TaskName string `json:"taskname"`
	BoardID  string `json:"boardid"`
	Status   string `json:"status"`
}

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponseBoard struct {
	Type    string  `json:"type"`
	Data    []Board `json:"data"`
	Message string  `json:"message"`
}

type JsonResponseTask struct {
	Type    string `json:"type"`
	Data    []Task `json:"data"`
	Message string `json:"message"`
}

func main() {

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all boards
	router.HandleFunc("/boards", GetBoards).Methods("GET")

	// Get all tasks
	router.HandleFunc("/tasks/{boardid}", GetTasks).Methods("GET")

	// Create a board
	router.HandleFunc("/boards", CreateBoard).Methods("POST")

	// Create a task
	router.HandleFunc("/tasks", CreateTask).Methods("POST")

	// Delete a task by the taskID
	router.HandleFunc("/tasks/{taskid}", DeleteTask).Methods("DELETE")

	// Update a task by the taskID
	router.HandleFunc("/tasks/{taskid}", UpdateTask).Methods("PATCH")

	// serve the app
	//fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get all movies
func GetBoards(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting boards...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM boards")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var boards []Board

	// Foreach movie
	for rows.Next() {
		var id int
		var boardName string

		err = rows.Scan(&id, &boardName)

		// check errors
		checkErr(err)

		boards = append(boards, Board{BoardID: id, BoardName: boardName})
	}

	var response = JsonResponseBoard{Type: "success", Data: boards}

	json.NewEncoder(w).Encode(response)
}

// Get all movies
func GetTasks(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting tasks...")

	params := mux.Vars(r)

	boardID := params["boardid"]

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM tasks WHERE tasks.board=$1", boardID)

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var tasks []Task

	// Foreach movie
	for rows.Next() {
		var id int
		var taskName, boardID, status string

		err = rows.Scan(&id, &taskName, &boardID, &status)

		// check errors
		checkErr(err)

		tasks = append(tasks, Task{TaskID: id, TaskName: taskName, BoardID: boardID, Status: status})
	}

	var response = JsonResponseTask{Type: "success", Data: tasks}

	json.NewEncoder(w).Encode(response)
}

// Create a board

// response and request handlers
func CreateBoard(w http.ResponseWriter, r *http.Request) {
	boardName := r.FormValue("boardname")

	var response = JsonResponseBoard{}

	if boardName == "" {
		response = JsonResponseBoard{Type: "error", Message: "You are missing boardName parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting board into DB")

		fmt.Println("Inserting new board with name: " + boardName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO boards(Name) VALUES($1) returning id;", boardName).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponseBoard{Type: "success", Message: "The board has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.FormValue("taskname")
	boardID := r.FormValue("boardid")

	var response = JsonResponseTask{}

	if boardID == "" {
		response = JsonResponseTask{Type: "error", Message: "You are missing boardid parameter."}
	} else if taskName == "" {
		response = JsonResponseTask{Type: "error", Message: "You are missing taskname parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting task into DB")

		fmt.Println("Inserting new task with name: " + taskName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO tasks(name, board) VALUES($1, $2) returning id;", taskName, boardID).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponseTask{Type: "success", Message: "The board has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// response and request handlers
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	taskID := params["taskid"]

	var response = JsonResponseTask{}

	if taskID == "" {
		response = JsonResponseTask{Type: "error", Message: "You are missing taskID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting task from DB")

		_, err := db.Exec("DELETE FROM tasks where id = $1", taskID)

		// check errors
		checkErr(err)

		response = JsonResponseTask{Type: "success", Message: "The task has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskStatus := r.FormValue("taskstatus") //"new"/"in progress"/"done"
	params := mux.Vars(r)

	taskID := params["taskid"]

	var response = JsonResponseTask{}

	if taskID == "" {
		response = JsonResponseTask{Type: "error", Message: "You are missing taskID parameter."}
	} else if taskStatus != "new" && taskStatus != "in progress" && taskStatus != "done" {
		response = JsonResponseTask{Type: "error", Message: "Task status must be new or in progress or done."}
	} else {
		db := setupDB()

		printMessage("Updating task")

		_, err := db.Exec("UPDATE tasks SET status = $1 where id = $2", taskStatus, taskID)

		// check errors
		checkErr(err)

		response = JsonResponseTask{Type: "success", Message: "The task has been updated successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all movies

// response and request handlers
func DeleteMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all movies...")

	_, err := db.Exec("DELETE FROM movies")

	// check errors
	checkErr(err)

	printMessage("All movies have been deleted successfully!")

	var response = JsonResponseBoard{Type: "success", Message: "All movies have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
