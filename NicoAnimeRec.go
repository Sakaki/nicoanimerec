package main

import (
        "flag"
)

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