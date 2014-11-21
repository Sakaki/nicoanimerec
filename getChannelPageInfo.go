package main

import (
	"github.com/PuerkitoBio/goquery"
)

func getRecoadableVideos(ch_id string) (videolst []string){
        pageurl := "http://ch.nicovideo.jp/" + ch_id
        doc, _ := goquery.NewDocument(pageurl)

        doc.Find("div.g-video-left").Each(func(_ int, videobox *goquery.Selection) {
       		exists := videobox.Find("span.inner.ppv.all_pay").Text()
		if exists == "" {
		   	videolink, _ := videobox.Find("a.g-video-link").Attr("href")
			for _, vlink := range videolst {
			        if vlink == videolink {
				        return
				}
			}
		        videolst = append(videolst, videolink)
		}
	})

	return
}
