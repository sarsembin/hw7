package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	CLIENTS = 50
)

func main() {

	var wg sync.WaitGroup
	for i := 0; i < CLIENTS; i++ {
		wg.Add(1)
		go func (num int) {
			defer wg.Done()
			// Подключаемся к сокету
			conn, _ := net.Dial("tcp", "127.0.0.1:8080")
			defer func() {
				err := conn.Close()
				if err != nil {
					log.Println(err.Error())
				}
			}()
			// Отправляем в socket
			fmt.Fprintf(conn, fmt.Sprintf("%f", float64(num)))
			// Прослушиваем ответ
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(message)
		}(i)
	}
	wg.Wait()
}
