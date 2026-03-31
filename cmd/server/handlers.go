package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ThatGuyMcFly/HttpDatabaseServer/internal/authentication"
	"github.com/ThatGuyMcFly/HttpDatabaseServer/internal/database"
)

func getHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func addEmployee(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	db := database.ConnectDatabase(database.EmployeeDatabase)

	newIds, err := database.AddEmployee(db, decoder)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	jsonData, err := json.Marshal(newIds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.Write(jsonData)
}

func extractQueryString(r *http.Request) string {
	var urlParts = strings.Split(r.RequestURI, "?")

	if len(urlParts) == 2 {
		return urlParts[1]
	}

	return ""
}

func writeDataAsJson(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(jsonData)
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDatabase(database.EmployeeDatabase)
	defer db.Database.Close()

	var queryString = extractQueryString(r)

	var employees = database.GetEmployees(db, queryString)

	writeDataAsJson(w, employees)
}

func getEmployeeById(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDatabase(database.EmployeeDatabase)
	defer db.Database.Close()
	employeeId := r.PathValue("employeeId")
	employees := database.GetEmployees(db, "employeeId="+employeeId)

	if len(employees) < 1 || employees == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
		return
	}

	writeDataAsJson(w, employees[0])
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Write Employee"))
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDatabase(database.EmployeeDatabase)
	defer db.Database.Close()

	var queryString = extractQueryString(r)

	err := database.DeleteEmployee(db, queryString)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Employee not found"))
	}

	w.Write([]byte("Employee deleted"))
}

func deleteEmployeeById(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDatabase(database.EmployeeDatabase)
	defer db.Database.Close()
	employeeId := r.PathValue("employeeId")
	err := database.DeleteEmployee(db, "employeeId="+employeeId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Employee not found"))
	}

	w.Write([]byte("Employee deleted"))
}

func addItem(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("addItem"))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDatabase(database.InventoryDatabase)
	defer db.Database.Close()

	var items = database.GetItems(db, nil)

	jsonData, err := json.Marshal(items)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(jsonData)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateItem"))
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteItem"))
}

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok {
		var employeeId, err = strconv.Atoi(username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid username or password"))
			return
		}

		authToken, err := authentication.AuthenticateUser(employeeId, password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		log.Println(authToken)
	}
}
