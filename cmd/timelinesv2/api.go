package main

import "net/http"

type api struct {
	version uint
	routes  []handler
}

type handler struct {
	name, path, method string
	handler            func(w http.ResponseWriter, r *http.Request)
}
