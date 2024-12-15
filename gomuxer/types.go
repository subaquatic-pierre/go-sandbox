package main

type ApiError struct {
	Error string `json:"error"`
}

type URL struct {
	Port     string `json:"port"`
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
}

type ResponseData struct {
	URL    URL
	Agent  string `json:"agent"`
	Body   string `json:"body"`
	Method string `json:"method"`
}
