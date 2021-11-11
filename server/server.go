package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
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

func maxClients(h http.Handler, n int) http.Handler {
	semaphore := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		h.ServeHTTP(w, r)
	})
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}

}


func serve(ctx context.Context) (err error) {

	handler := http.HandlerFunc(squareEndpoint)
	heavyHandler := http.HandlerFunc(heavyOperation)

	mux := http.NewServeMux()

	mux.Handle("/square", maxClients(handler, 10))
	mux.Handle("/heavy", maxClients(heavyHandler, 1))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server shutting down...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}


