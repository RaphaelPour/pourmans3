package main

import(
	"fmt"
	"net/http"
	"log"
	"hash/maphash"
	"encoding/base64"
)

var hash maphash.Hash
var links = map[string]string{
	"123":"234",
}

func shorten(link string) string {
	hash.Reset()
	hash.WriteString(link)
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x",hash.Sum64())))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		for short, long := range links {
			fmt.Fprintf(w, "%s -> %s\n", short, long)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
