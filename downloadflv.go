package main

import (
    "strings"
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "fmt"
    "strconv"
    "io"
    "os"
    "encoding/xml"
    "encoding/json"
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

func urlEncode(source string) (decoded string){
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
    fmt.Println(decoded)
    return
}

func downloadFromUrl(url string, fname string, client *http.Client) {
     fmt.Println("Downloading", url, "to", fname)

     output, err := os.Create(fname)
     if err != nil {
     	fmt.Println("Error while creating", fname, "-", err)
     	return
     }
     defer output.Close()
     req, _ := http.NewRequest("GET", url, nil)
     response, err := client.Do(req)
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

type Thumb struct {
      Title string `xml:"thumb>title"`
}

func readVName(vid string) (vname string) {
    var vnames map[string]string 

    bdata, _ := os.Open(absPath("./data/vnames.json"))
    dec := json.NewDecoder(bdata)
    dec.Decode(&vnames)

    vname = vnames[vid]
    if vname == "" {
	vname = "none"
    }

    return
}

func saveVName(vid string, vname string) {
    vnames := make(map[string]string)

    bdata, _ := os.Open(absPath("./data/vnames.json"))
    dec := json.NewDecoder(bdata)
    dec.Decode(&vnames)

    vnames[vid] = vname

    jsonstr, _ := json.Marshal(vnames)
    fmt.Println(vnames)
    ioutil.WriteFile(absPath("./data/vnames.json"), jsonstr, 0644)
}

func Download(videoid string, location string, mail string, passwd string) {
    client := login(mail, passwd)

    vreq, _ := http.NewRequest("GET", "http://flapi.nicovideo.jp/api/getflv/"+videoid, nil)
    vres, _ := client.Do(vreq)
    apibody, _ := ioutil.ReadAll(vres.Body)

    vres.Body.Close()

    xmlres, _ := http.Get("http://ext.nicovideo.jp/api/getthumbinfo/"+videoid)
    rawxml, _ := ioutil.ReadAll(xmlres.Body)

    dreq, _ := http.NewRequest("GET", "http://www.nicovideo.jp/watch/"+videoid, nil)
    dres, _ := client.Do(dreq)
    dres.Body.Close()

    nicoXml := Thumb{""}
    xml.Unmarshal(rawxml, &nicoXml)

    apiinfo := strings.Split(string(apibody[:]), "&")[2][4:]
    vfile := location+"/"+videoid
    downloadFromUrl(urlEncode(apiinfo), vfile, client)
    saveVName(videoid, nicoXml.Title)
}
