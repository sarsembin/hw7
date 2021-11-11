package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)


func main() {
	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go postSquare("http://localhost:8080/square", i / 100, &wg)
		//go getHeavy("http://localhost:8080/heavy", i, &wg)
	}
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Elapsed time: %v\n", elapsed)
}



func postSquare(url string, num int, wg *sync.WaitGroup) float64 {
	defer wg.Done()

	postBody, err := json.Marshal(struct {
		Number int `json:"number"`
	}{
		Number: num,
	})
	if err != nil {
		log.Fatal(err)
	}

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(jsonData["square"])
	return jsonData["square"].(float64)
}



func getHeavy(url string, number int, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Operation number: %v, response: %v\n", number, string(body))
}
