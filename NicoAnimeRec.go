package main

import (
	"bitbucket.org/kardianos/osext"
	"flag"
	"strings"
)

func absPath(path string) string {
	runpath, _ := osext.ExecutableFolder()
	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}

	return runpath + path
}

func main() {
	var (
		recmode    bool
		servermode bool
		port       string
		mail       string
		passwd     string
	)

	flag.BoolVar(&recmode, "r", false, "record mode")
	flag.BoolVar(&servermode, "s", false, "server mode")
	flag.StringVar(&port, "P", "8080", "server port")
	flag.StringVar(&mail, "m", "", "login mail addr")
	flag.StringVar(&passwd, "p", "", "login passwd")
	flag.Parse()

	if recmode {
		GetAllAnimes()
		DoRecord(mail, passwd)
	} else if servermode {
		Server(port)
	}
}
