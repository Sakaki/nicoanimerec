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

func exeRec() {
        targets := readConf()
	for _, ch_id := range targets {
	        mkdircmd := exec.Command("mkdir", "videos/"+ch_id)
		mkdircmd.Run()
		videos := getRecoadableVideos(ch_id)
		for _, video := range videos {
		         video = strings.Replace(video, "http://www.nicovideo.jp/watch/", "", -1)
			 fmt.Println(video)
	                 reccmd := exec.Command("python", "pyniconico/downloadflv.py", "-l", "videos/"+ch_id, video)
			 reccmd.Run()
		}
	}
}
