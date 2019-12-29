package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	messages "github.com/yuvalili138/golang_exc/pkg/messages"
)

func main() {
	urls := [...]messages.UrlMessage{
		messages.UrlMessage{
			Domain: "ynet.co.il",
			Path:   "/page=2",
		},
		messages.UrlMessage{
			Domain: "ynet.co.il",
			Path:   "/page=1",
		},
		messages.UrlMessage{
			Domain: "shesh.co.il",
			Path:   "/chat",
		},
		messages.UrlMessage{
			Domain: "ynet.co.il",
			Path:   "/page=2",
		},
		messages.UrlMessage{
			Domain: "ynet.co.il",
			Path:   "/page=1",
		},
		messages.UrlMessage{
			Domain: "shesh.co.il",
			Path:   "/chat",
		},
		messages.UrlMessage{
			Domain: "harta.co.il",
			Path:   "/chat",
		},
		messages.UrlMessage{
			Domain: "harta.co.il",
			Path:   "/lolz",
		},
		messages.UrlMessage{
			Domain: "ynet.co.il",
			Path:   "/page=2",
		},
	}

	j := 0
	for i := 0; i < 200; i++ {
		// Create message
		clientMessage := urls[j % len(urls)]
		j++

		msg, err  := clientMessage.MarshalJSON()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Send message
		url := "http://127.0.0.1:8080/yuvalili"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(msg))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}

	//// Create message
	//clientMessage := messages.UrlMessage {
	//	Domain: os.Args[1],
	//	Path: os.Args[2],
	//}
	//
	//msg, err  := clientMessage.MarshalJSON()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//// Send message
	//url := "http://127.0.0.1:8080/yuvalili"
	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(msg))
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}
