package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"time"
)

type AnimeInfo struct {
	Title   string
	Channel string
	Date    string
	Thumb   string
	Weekday time.Weekday
}

func getAllAnimes() {
	var (
		title      string
		channel_id string
		start_time string
		thumb_url  string
	)

	weekdaylst := map[string]time.Weekday{
		"playerNav3": time.Monday,
		"playerNav4": time.Tuesday,
		"playerNav5": time.Wednesday,
		"playerNav6": time.Thursday,
		"playerNav7": time.Friday,
		"playerNav8": time.Saturday,
		"playerNav9": time.Sunday}

	var animelst []AnimeInfo
	var tmpchs []string

	AnimeListUrl := "http://ch.nicovideo.jp/portal/anime?cc_referrer=nicotop_sidemenu"
	doc, _ := goquery.NewDocument(AnimeListUrl)
	doc.Find("div.playerNav").Each(func(_ int, s *goquery.Selection) {
		navid, _ := s.Attr("id")
		if weekday, exists := weekdaylst[navid]; exists {
			s.Find("li.video.cfix").Each(func(_ int, t *goquery.Selection) {
				t.Find("input").Each(func(_ int, u *goquery.Selection) {
					key, _ := u.Attr("name")
					value, _ := u.Attr("value")

					switch key {
					case "title":
						title = value
					case "channel_id":
						channel_id = value
					case "thumbnail_url":
						thumb_url = value
					case "start_time":
						start_time = value
					}
				})
				if channel_id == "" {
					return
				}
				for _, item := range tmpchs {
					if item == channel_id {
						return
					}
				}
				tmpchs = append(tmpchs, channel_id)

				info := AnimeInfo{Title: title, Channel: channel_id, Date: start_time, Thumb: thumb_url, Weekday: weekday}
				animelst = append(animelst, info)
			})
		}
	})

	jsonstr, _ := json.Marshal(animelst)
	ioutil.WriteFile("data/animelst.json", jsonstr, 0644)
}
