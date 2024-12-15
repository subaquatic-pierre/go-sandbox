package main

import "net/http"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, 200, ApiError{Error: "no route found"})
}

func routeInfo(w http.ResponseWriter, r *http.Request) {
	data := ResponseData{
		URL: URL{Hostname: r.URL.Hostname(), Port: r.URL.Port(), Path: r.URL.Path},
		// Body:   string(body[:]),
		Agent:  r.UserAgent(),
		Method: r.Method,
	}

	WriteJson(w, 200, &data)
}
