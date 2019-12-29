package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	db "github.com/yuvalili138/golang_exc/pkg/db"
	messages "github.com/yuvalili138/golang_exc/pkg/messages"
)

func handleRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		_, _ = io.WriteString(w, "Supporting only POST requests")
		return
	}

	// Parse client message
	body, _ := ioutil.ReadAll(req.Body)
	clientMessage := &messages.UrlMessage{}
	if err := clientMessage.UnmarshalJSON(body); err != nil {
		_, _ = io.WriteString(w, fmt.Sprintf("Error parsing message: %s", err.Error()))
		return
	}

	// Get result
	exists := db.Get(clientMessage.Domain, clientMessage.Path)
	result := ""
	if exists {
		result = fmt.Sprintf("https://%s%s", clientMessage.Domain, clientMessage.Path)
	}

	// Prepare the response
	response := messages.LocationMessage{Location: result}
	msg, err := response.MarshalJSON()
	if err != nil {
		fmt.Println("Oi Vei!")
		return
	}

	_, _ = w.Write(msg)
}

func main() {
	http.HandleFunc("/yuvalili", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}