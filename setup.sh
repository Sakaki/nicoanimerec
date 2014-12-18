#!/bin/sh

go get "bitbucket.org/kardianos/osext"
go get "github.com/PuerkitoBio/goquery"
go get "github.com/drone/routes"

go build -o NicoAnimeRec ./*.go
