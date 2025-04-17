package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	testRoutes := routes()
	chiRoutes := testRoutes.(chi.Routes)
	
	// Auth routes
	routeExists(t, chiRoutes, "/users/login")
	routeExists(t, chiRoutes, "/users/logout")
	routeExists(t, chiRoutes, "/validate-token")

	// Admin user routes
	routeExists(t, chiRoutes, "/admin/users")
	routeExists(t, chiRoutes, "/admin/users/save")
	routeExists(t, chiRoutes, "/admin/users/get/{id}")
	routeExists(t, chiRoutes, "/admin/users/delete")
	routeExists(t, chiRoutes, "/admin/log-user-out/{id}")

	// Admin tag routes
	routeExists(t, chiRoutes, "/admin/dashboard/tags")
	routeExists(t, chiRoutes, "/admin/dashboard/tags/{id}")

	// Admin budget routes
	routeExists(t, chiRoutes, "/admin/dashboard/budgets")
	routeExists(t, chiRoutes, "/admin/dashboard/budgets/{id}")

	// Admin transaction recurrence routes
	routeExists(t, chiRoutes, "/admin/dashboard/transactions-recurrences")

	// Admin transaction routes
	routeExists(t, chiRoutes, "/admin/dashboard/transactions")
	routeExists(t, chiRoutes, "/admin/dashboard/transactions/{id}")
	routeExists(t, chiRoutes, "/admin/dashboard/transactions/set-all-active")
	routeExists(t, chiRoutes, "/admin/dashboard/transactions/set-status")

	// Admin logs routes
	routeExists(t, chiRoutes, "/admin/dashboard/logs")

	// Static files route
	routeExists(t, chiRoutes, "/static/*")
}

func routeExists(t *testing.T, routes chi.Routes, route string) {

	found := false

	chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if foundRoute == route {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("route %s not found", route)
	}
}
