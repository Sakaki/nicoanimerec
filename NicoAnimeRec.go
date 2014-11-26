package main

import (
        "flag"
	"bitbucket.org/kardianos/osext"
	"strings"
)

func absPath(path string) string {
        runpath, _ := osext.ExecutableFolder()
	if strings.HasPrefix(path, "./") {
	        path = path[2:]
	}

	return runpath+path
}

func main() {
        var (
	        recmode bool
		servermode bool
		port string
	)

	flag.BoolVar(&recmode, "r", false, "record mode")
	flag.BoolVar(&servermode, "s", false, "server mode")
	flag.StringVar(&port, "p", "8080", "server port")
	flag.Parse()

	if recmode {
	        getAllAnimes()
                exeRec()
	}else if servermode {
	        Server(port)
	}
}