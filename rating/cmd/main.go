package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"movieexample.com/rating/internal/controller/rating"
	httphandler "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
)

var serviceName = "rating"

func main() {
	// получаю порт через флаги
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting the rating service on port %d", port)
	// делаю новый экземпляр registry, здесь consul
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// создаю ID
	instanceID := discovery.GenerateInstanceID(serviceName)
	// регистрирую сервис в registry
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	// каждую секунду говорю что сервис жив
	go func() {
		for {
			if err := registry.
				ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	// при закрытии сервиса удаляю сервис из registry
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
