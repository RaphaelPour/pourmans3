package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	POST_TEMPLATE = `
	<html>
	<head><title>pourmans3</title></head>
	<body>
	<h1 style='text-align:left'>pourmans3</h1>
	<h3 style='text-align:left'><em>Link shortener, designed to store heave-payloaded urls</em></h3>
	<hr>
	<form method='POST'>
	<input type='text' name='url'>
	<input type='submit' value='Shorten'>
	</form>
	<ul>
	{{ range $short, $long := .Links }}
		<li><a href='?{{ $short }}'>{{ $short }}</a> -> <a href='{{ $long }}'>{{ $long }}</a></li>
	{{ end }}
	</ul>
	<hr>
	<ul style="font-family: monospace;">
		<li>POST /{original link} -> {shortened link}</li>
		<li>GET /{shortened link} -> {original link}</li>
		<li>DELETE /{shortened link}</li>
	</ul>
	</body>
	</html>
	`
	GET_TEMPLATE = `
	<html>
	<head>
		<meta http-equiv="Refresh" content="0; URL={{ . }}">
	</head>
	</html>
	`
)

type Service struct {
	port         int
	postTemplate *template.Template
	getTemplate  *template.Template
	storage      Storage
}

type ShortLink struct {
	Links map[string]string
	Host  string
}

func NewService(port int) (*Service, error) {
	postTemplate, err := template.New("post").Parse(POST_TEMPLATE)
	if err != nil {
		return nil, fmt.Errorf("error parsing post template: %s", err)
	}

	getTemplate, err := template.New("get").Parse(GET_TEMPLATE)
	if err != nil {
		return nil, fmt.Errorf("error parsing get template: %s", err)
	}

	return &Service{
		port:         port,
		postTemplate: postTemplate,
		getTemplate:  getTemplate,
		storage:      NewStorage(),
	}, nil
}

func (s *Service) Start() error {
	http.HandleFunc("/", s.RequestHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", *Port), nil)
}

func (s *Service) RequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(
		"%s [%6s] %s\n",
		time.Now().UTC().Format("2006-01-02T15:04:05"),
		r.Method,
		r.URL,
	)

	if r.Method == http.MethodPost {
		if r.FormValue("url") == "" {
			fmt.Println("error on post: url parameter missing")
			return
		}

		key := s.storage.Set(r.FormValue("url"))
		fmt.Fprintf(w, "%s/?%s", r.Host, key)
		return
	}

	if r.Method == http.MethodGet && r.URL.RawQuery != "" {

		value, err := s.storage.Get(r.URL.RawQuery)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := s.getTemplate.Execute(w, value); err != nil {
			fmt.Println("error rendering get template:", err)
		}
		return
	}

	if r.Method == http.MethodDelete && r.URL.RawQuery != "" {
		s.storage.Delete(r.URL.RawQuery[1:])
		return
	}

	if err := s.postTemplate.Execute(w,
		ShortLink{
			Links: s.storage.All(),
			Host:  r.Host,
		}); err != nil {
		fmt.Println("error rendering post template:", err)
		return
	}
}
