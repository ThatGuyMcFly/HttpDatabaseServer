package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ThatGuyMcFly/HttpDatabaseServer/internal/auth"
	"github.com/ThatGuyMcFly/HttpDatabaseServer/internal/database"
	"gopkg.in/ini.v1"
)

const DefaultAdminId = 1000000

func registerEmployeeRoutes() {

	var err = addRoute(POST, "/employees", addEmployee)
	if err != nil {
		log.Println("register add employee route failed:", err)
	}

	err = addRoute(GET, "/employees", getEmployees)
	if err != nil {
		log.Println("register get employee route failed:", err)
	}

	err = addRoute(GET, "/employees/{employeeId}", getEmployeeById)
	if err != nil {
		log.Println("register get employee by id route failed:", err)
	}

	err = addRoute(PUT, "/employees/{employeeId}", updateEmployee)
	if err != nil {
		log.Println("register put employee route failed:", err)
	}

	err = addRoute(DELETE, "/employees/{employeeId}", deleteEmployeeById)
}

func registerItemRoutes() {
	var err = addRoute(POST, "/items", addItem)
	if err != nil {
		log.Println("register add item route failed:", err)
	}

	err = addRoute(GET, "/items", getItems)
	if err != nil {
		log.Println("register get item route failed:", err)
	}

	err = addRoute(PUT, "/items", updateItem)
	if err != nil {
		log.Println("register put item route failed:", err)
	}

	err = addRoute(DELETE, "/items", deleteItem)
	if err != nil {
		log.Println("register delete item route failed:", err)
	}
}

func registerAuthenticationRoutes() {
	var err = addRoute(GET, "/authenticate", authenticateUser)
	if err != nil {
		log.Println("register authenticate route failed:", err)
	}
}

func initializeAdmin() bool {
	db := database.ConnectDatabase(database.EmployeeDatabase)

	employees := database.GetEmployees(db, "employeeId="+strconv.Itoa(DefaultAdminId))

	if len(employees) == 0 {
		admin := map[string]any{
			"employeeId": DefaultAdminId,
			"firstName":  "Admin",
			"lastName":   "Admin",
			"role":       "admin",
		}
		_, err := database.AddEmployee(db, admin)
		if err != nil {
			log.Println(err)
			return false
		}

		err = auth.SetUserPassword(DefaultAdminId, "admin")

		if err != nil {
			log.Println(err)
			return false
		}
	}

	return true
}

func main() {

	if !initializeAdmin() {
		log.Println("initializeAdmin failed")
		return
	}

	cfg, err := ini.Load("cmd/server/config/config.ini")

	if err != nil {
		log.Fatal(err)
	}

	httpPort := cfg.Section("server").Key("port").MustInt(8080)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(httpPort),
		Handler: routes(),
	}

	registerEmployeeRoutes()
	registerItemRoutes()
	registerAuthenticationRoutes()

	log.Println("Starting server on port " + strconv.Itoa(httpPort))

	err = server.ListenAndServe()

	if err != nil {
		return
	}
}
