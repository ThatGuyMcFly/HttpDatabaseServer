package main

import (
	"log"
	"net/http"
	"strconv"

	"gopkg.in/ini.v1"
)

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

	err = addRoute(DELETE, "/employees", deleteEmployee)
	if err != nil {
		log.Println("register delete employee route failed:", err)
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

func main() {

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
