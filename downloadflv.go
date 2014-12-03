package main

import (
    "strings"
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"
    "net/url"
//    "regexp"
    "fmt"
    "strconv"
    "io"
    "os"
)

func login(mail string, passwd string) *http.Client {
    cookieJar, _ := cookiejar.New(nil)

    client := &http.Client {
        Jar: cookieJar,
    }

    apiUrl := "https://secure.nicovideo.jp/secure/login/"
    data := url.Values{}
    data.Add("mail", mail)
    data.Add("password", passwd)
    data.Add("next_url", "")
    data.Add("site", "niconico")

    res, _ := client.PostForm(apiUrl, data)
    res.Body.Close()

    return client
}

func UrlEncode(source string) (decoded string){
    for idx := 0; idx < len(source); idx++ {
    	if source[idx] == '%' {
	   numchar, _ := strconv.ParseUint(source[idx+1:idx+3], 16,0)
	   char := fmt.Sprintf("%c", numchar)
	   decoded += char
	   idx += 2
	} else {
	   decoded += string(source[idx])
	}
    }

    return
}

func downloadFromUrl(url string, fname string) {
     fmt.Println("Downloading", url, "to", fname)

     output, err := os.Create(fname)
     if err != nil {
     	fmt.Println("Error while creating", fname, "-", err)
     	return
     }
     defer output.Close()
     response, err := http.Get(url)
     if err != nil {
     	fmt.Println("Error while downloading", url, "-", err)
     	return
     }
     defer response.Body.Close()
     n, err := io.Copy(output, response.Body)
     if err != nil {
     	fmt.Println("Error while downloading", url, "-", err)
	return
     }
     fmt.Println(n, "bytes downloaded.")
}

func main() {
    videoid := "sm24905365"
    client := login("sakakicks333@yahoo.co.jp", "begizagon")

    vreq, _ := http.NewRequest("GET", "http://flapi.nicovideo.jp/api/getflv/"+videoid, nil)
    vres, _ := client.Do(vreq)
    apibody, _ := ioutil.ReadAll(vres.Body)

    vres.Body.Close()

    xml := http.Get("http://ext.nicovideo.jp/api/getthumbinfo/"+videoid)
    

    apiinfo := strings.Split(string(apibody[:]), "&")[2][4:]
    downloadFromUrl(UrlEncode(apiinfo), "testvideo.flv")
}
