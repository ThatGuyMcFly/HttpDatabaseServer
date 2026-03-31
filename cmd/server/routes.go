package main

import (
	"log"
	"net/http"
	"regexp"
	"slices"
)

type RouteExistsError struct{}

func (e RouteExistsError) Error() string {
	return "Route conflicts with existing route"
}

type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	PUT
	PATCH
	DELETE
)

func (m HTTPMethod) toString() string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}[m]
}

var routesMap = map[HTTPMethod][]string{
	GET:    {},
	POST:   {},
	PUT:    {},
	PATCH:  {},
	DELETE: {},
}

var re = regexp.MustCompile(`\{.*?}`)

var mux = http.NewServeMux()

func routes() http.Handler {
	mux.HandleFunc("GET /", getHome)
	return mux
}

func routeExists(method HTTPMethod, route string) bool {

	newRoute := re.ReplaceAllString(route, "")

	if slices.Contains(routesMap[method], newRoute) {
		return true
	}

	return false
}

/*
 */
func addRoute(method HTTPMethod, path string, handler func(http.ResponseWriter, *http.Request)) error {
	if routeExists(method, path) {
		return RouteExistsError{}
	}

	storedRoute := re.ReplaceAllString(path, "")

	routesMap[method] = append(routesMap[method], storedRoute)

	var route = method.toString() + " " + path

	log.Println("Adding route: " + route)
	mux.HandleFunc(route, handler)

	return nil
}
