package main

import (
	"log"
	"net/http"
	"testing"
)

func TestSendPostRequest(t *testing.T) {
	request, err := SendPostRequest()
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
	log.Print(request)
}

func TestFindToken(t *testing.T) {
	client := &http.Client{}
	url := "https://xskydata.jobs.feishu.cn/api/v1/search/job/posts"
	request, _ := http.NewRequest("POST", url, nil)
	FindToken(request)
	if _, err = client.Do(request); err != nil {
		log.Println(err)
		t.Error(err)
	}
}
