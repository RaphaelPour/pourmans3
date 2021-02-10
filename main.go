package main // import "github.com/RaphaelPour/pourmans3"

import (
	"flag"
	"fmt"
	"hash/maphash"
	"html/template"
	"net/http"
)

const (
	PORT          = 80
	POST_TEMPLATE = `
	<html>
	<body>
	<form method='POST'>
	URL:<input type='text' name='url'>
	<input type='submit' value='Shorten'>
	</form>
	<hr>
	<ul>
	{{ range $short, $long := .Links }}
		<li><a href='?{{ $short }}'>{{ $short }}</a> -> <a href='{{ $long }}'>{{ $long }}</a></li>
	{{ end }}
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
	Version      = flag.Bool("version", false, "Print build information")
	hash         maphash.Hash
	links        = map[string]string{
		"about": "https://evilcookie.de",
	}
)

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

	postTemplate, err := template.New("test").Parse(POST_TEMPLATE)
	if err != nil {
		fmt.Println("error parsing html template:", err)
		return
	}

	getTemplate, err := template.New("redirrect").Parse(GET_TEMPLATE)
	if err != nil {
		fmt.Println("error parsing html template:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header["Refferer"])
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

			fmt.Fprintf(
				w,
				"%s/?%s",
				host,
				short,
			)
			return
		}

		if r.Method == http.MethodGet && r.URL.RawQuery != "" {

			if err := getTemplate.Execute(w, extend(r.URL.RawQuery)); err != nil {
				fmt.Println("error rendering get template:", err)
			}
			return
		}

		if err := postTemplate.Execute(w,
			struct {
				Links map[string]string
				Host  string
			}{
				Links: links,
				Host:  host,
			}); err != nil {
			fmt.Println("error rendering post template:", err)
			return
		}
	})

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
