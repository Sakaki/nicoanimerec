package main

import (
        "fmt"
	"os"
	"path/filepath"
	"os/exec"
	"strings"
)

func getExecLoc() string{
        filename := os.Args[0]
	fmt.Println(os.Args)
    	filedir := filepath.Dir(filename)
    	abspath, _ := filepath.Abs(filedir)

	return abspath
}

func exeRec(mail string, passwd string) {
        targets := readConf()
	for _, ch_id := range targets {
	        loc := absPath("./videos/")+ch_id
	        mkdircmd := exec.Command("mkdir", loc)
		mkdircmd.Run()
		videos := getRecoadableVideos(ch_id)
		for _, video := range videos {
		         video = strings.Replace(video, "http://www.nicovideo.jp/watch/", "", -1)
			 fmt.Println(video)
			 Download(video, loc, mail, passwd)
		}
	}
}
