package main

import (
	"encoding/json"
	"fmt"
	"github.com/drone/routes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func makeHeader(title string) (header string) {
	header = loadHtml("header")
	header = strings.Replace(header, "%Title%", title, -1)

	return header
}

func loadHtml(filename string) string {
	bdata, _ := ioutil.ReadFile(absPath("./data/html/" + filename))
	return string(bdata)
}

func redirect_home(w http.ResponseWriter, r *http.Request) {
	home(w, r)
}

func home(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filter := params.Get(":filter")

	var (
		animelst []AnimeInfo
		trflag   bool
	)

	//要素の読み込み
	header := makeHeader("メインページ")
	footer := loadHtml("footer")
	titleview := loadHtml("main_titleview")
	lstview := loadHtml("main_lst")
	traytmp := loadHtml("main_tray")

	//アニメ一覧の読み込み
	reader, _ := os.Open(absPath("./data/animelst.json"))
	dec := json.NewDecoder(reader)
	dec.Decode(&animelst)

	//録画設定の読み込み
	rec_chs := readConf()

	fmt.Fprintf(w, header)
	fmt.Fprintf(w, titleview)
	fmt.Fprintf(w, lstview)
	for _, data := range animelst {
		tray := strings.Replace(traytmp, "%Weekday%", time.Weekday(data.Weekday).String(), -1)
		tray = strings.Replace(tray, "%Thumb%", data.Thumb, -1)
		tray = strings.Replace(tray, "%BtnId%", data.Channel, -1)
		btntxt := "録画"
		for _, item := range rec_chs {
			if item == data.Channel {
				btntxt = "録画中"
			}
		}
		if btntxt != "録画中" && filter == "recording" {
			continue
		}
		trflag = !trflag
		traycolor := ""
		if trflag {
			traycolor = "class=\"pure-table-odd\""
		}
		tray = strings.Replace(tray, "%Color%", traycolor, -1)
		tray = strings.Replace(tray, "%ToRec%", btntxt, -1)
		tray = strings.Replace(tray, "%Title%", data.Title, -1)
		tray = strings.Replace(tray, "%Channel%", "<a href=\"./channel/"+data.Channel+"\">"+data.Channel+"</a>", -1)
		fmt.Fprintf(w, tray)
	}
	fmt.Fprintf(w, string(footer))
}

func channel(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")

	fmt.Fprintf(w, "<p><a href=\"../\">HOME</a></p><br>")
	fmt.Fprintf(w, makeHeader("チャンネル"))
	footer := loadHtml("footer")

	files, _ := ioutil.ReadDir(absPath("./videos/" + id))
	for _, video := range files {
		vname := video.Name()
		fmt.Fprintf(w, "<p><a href=\"../watch/"+id+"/"+vname+"\">"+readVName(vname)+"</a></p>")
	}
	fmt.Fprintf(w, footer)
}

func watch(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	video := params.Get(":video")

	moviepage := loadHtml("watch")
	fmt.Fprintf(w, strings.Replace(moviepage, "%Video%", id+"/"+video, -1))
}

func readConf() (reader []string) {
	bdata, _ := os.Open(absPath("./data/config.json"))
	dec := json.NewDecoder(bdata)
	dec.Decode(&reader)

	return
}

func writeConf(writer []string) {
	jsonstr, _ := json.Marshal(writer)
	ioutil.WriteFile(absPath("./data/config.json"), jsonstr, 0644)
}

func addrec(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")

	rec_chs := readConf()

	for _, item := range rec_chs {
		if item == id {
			return
		}
	}
	rec_chs = append(rec_chs, id)

	writeConf(rec_chs)
	fmt.Fprintf(w, "complete!")
}

func delrec(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")

	rec_chs := readConf()

	for num, item := range rec_chs {
		if item == id {
			rec_chs = append(rec_chs[:num], rec_chs[num+1:]...)
		}
	}
	writeConf(rec_chs)
	fmt.Fprintf(w, "complete!")
}

func Server(port string) {
	mux := routes.New()

	pwd, _ := os.Getwd()
	mux.Static("/videos", pwd)

	mux.Get("/", redirect_home)
	mux.Get("/:filter", home)
	mux.Get("/channel/:id", channel)
	mux.Get("/add/:id", addrec)
	mux.Get("/del/:id", delrec)
	mux.Get("/watch/:id/:video", watch)

	http.Handle("/", mux)
	http.ListenAndServe("0.0.0.0:"+port, nil)
}
