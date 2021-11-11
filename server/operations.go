package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func squareEndpoint(writer http.ResponseWriter, request *http.Request)  {
	resp, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Fatal(err)
	}

	num := data["number"].(float64)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(struct {
		Square float64 `json:"square"`
	}{
		Square: num * num,
	})
}


func heavyOperation(writer http.ResponseWriter, request *http.Request) {
	for i := 0; i < 1e6; i++ {
		_ = strconv.Itoa(i)
	}
	_, _ = writer.Write([]byte("Done"))
}
