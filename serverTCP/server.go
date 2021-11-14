package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
Реализовать TCP-сервер, который возводит переданное ему число в
квадрат и возвращает результат. Количество обрабатываемых запросов
в один момент времени должно быть настраиваемым. Должен быть
предусмотрен graceful shutdown (перед завершением программы
необходимо обработать все открытые соединения).
Написать клиент для тестирования сервера
*/

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
	CONN_AMOUNT = 5
)

func main() {
	server := NewServer(CONN_HOST + ":" + CONN_PORT)


	fmt.Print("Enter command: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "stop" {
			server.Stop()
			break
		}
	}


}

// quit for signal to quit
// wg for waiting for request to process
type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

// constructor for Server
func NewServer(addr string) *Server {
	s := &Server{
		quit: make(chan interface{}),
	}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s.listener = listen
	s.wg.Add(1)
	go s.serve()

	log.Printf("listening on: %s\n", addr)
	return s
}

// creates new requestHandler for requests
func (s *Server) serve() {
	defer s.wg.Done()

	// semaphore for regulating request per unit of time
	semaphore := make(chan struct{}, CONN_AMOUNT)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				semaphore <- struct{}{}
				s.requestHandler(conn, semaphore)
				s.wg.Done()
			}()
		}
	}
}


// function for graceful shutdown
func (s *Server) Stop() {
	log.Println("shutting down server...")
	close(s.quit)
	s.listener.Close()
	log.Println("processing remaining requests...")
	s.wg.Wait()
}


// processes request
func (s *Server) requestHandler(conn net.Conn, semaphore chan struct{}) {
	defer func() { <- semaphore } ()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Make a buffer to hold incoming data.
	buff := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buff)
	if err != nil {
		log.Println("Error reading:", err.Error())
	}

	buff = delEmptyFromBuff(buff)
	result := numSquarer(buff)

	time.Sleep(2 * time.Second)

	_, err = conn.Write([]byte("Message received, result: " + fmt.Sprintf("%f", result)))
	if err != nil {
		log.Println("Error writing:", err.Error())
	}
}


// clears buffer form empty values
func delEmptyFromBuff(buff []byte) []byte {
	var index int
	for i, v := range buff {
		if v == 0 {
			index = i
			break
		}
	}
	return buff[:index]
}

// squares float from buffer
func numSquarer(buff []byte) float64 {
	float, err := strconv.ParseFloat(string(buff), 64)

	if err != nil {
		log.Println("Error converting: ", err.Error())
		return 0.0
	}

	return float * float
}




