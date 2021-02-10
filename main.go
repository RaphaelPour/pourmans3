package main // import "github.com/RaphaelPour/pourmans3"

import (
	"flag"
	"fmt"
	"hash/maphash"
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
	<h3 style='text-align:left'><em>Link shortener, designed to store heave-payloaded urls.</em></h3>
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

var (
	BuildDate    string
	BuildVersion string
	Port         = flag.Int("port", 80, "Listen port")
	Version      = flag.Bool("version", false, "Print build information")
	hash         maphash.Hash
	links        = map[string]string{
		"about": "https://evilcookie.de",
	}
	postTemplate *template.Template
	getTemplate  *template.Template
)

type ShortLink struct {
	Links map[string]string
	Host  string
}

func shorten(link string) string {
	hash.Reset()
	hash.WriteString(link)
	return fmt.Sprintf("%x", hash.Sum64())
}

func extend(short string) string {
	if long, ok := links[short]; ok {
		return long
	}
	fmt.Println("unknown short link ", short)
	return "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
}

func main() {

	flag.Parse()

	if *Version {
		fmt.Println("BuildVersion: ", BuildVersion)
		fmt.Println("BuildDate: ", BuildDate)
		return
	}

	var err error
	postTemplate, err = template.New("post").Parse(POST_TEMPLATE)
	if err != nil {
		fmt.Println("error parsing post template:", err)
		return
	}

	getTemplate, err = template.New("get").Parse(GET_TEMPLATE)
	if err != nil {
		fmt.Println("error parsing get template:", err)
		return
	}

	http.HandleFunc("/", RequestHandler)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", *Port), nil))
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(
		"%s [%6s] %s\n",
		time.Now().UTC().Format("2006-01-02T15:04:05"),
		r.Method,
		r.URL,
	)

	hosts := r.Header["Refferer"]
	var host string
	if len(hosts) == 0 {
		host = "localhost"
	} else {
		host = hosts[0]
	}

	if r.Method == http.MethodPost {
		if r.FormValue("url") == "" {
			fmt.Println("error on post: url parameter missing")
			return
		}
		short := shorten(r.FormValue("url"))
		links[short] = r.FormValue("url")

		fmt.Fprintf(w, "%s/?%s", host, short)
		return
	}

	if r.Method == http.MethodGet && r.URL.RawQuery != "" {

		if err := getTemplate.Execute(w, extend(r.URL.RawQuery)); err != nil {
			fmt.Println("error rendering get template:", err)
		}
		return
	}

	if r.Method == http.MethodDelete && r.URL.RawQuery != "" {
		delete(links, r.URL.RawQuery)
		return
	}

	if err := postTemplate.Execute(w,
		ShortLink{
			Links: links,
			Host:  host,
		}); err != nil {
		fmt.Println("error rendering post template:", err)
		return
	}
}
