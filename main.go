package main // import "github.com/RaphaelPour/pourmans3"

import (
	"flag"
	"fmt"
)

var (
	BuildDate    string
	BuildVersion string
	Port         = flag.Int("port", 80, "Listen port")
	Version      = flag.Bool("version", false, "Print build information")
)

func main() {

	flag.Parse()

	if *Version {
		fmt.Println("BuildVersion: ", BuildVersion)
		fmt.Println("BuildDate: ", BuildDate)
		return
	}

	service, err := NewService(*Port)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = service.Start(); err != nil {
		fmt.Println(err)
		return
	}
}
