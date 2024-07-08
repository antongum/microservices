package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"movieexample.com/metadata/internal/controller/metadata"
	httphandler "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

// константа названия сервиса
const serviceName = "metadata"

func main() {
	// ввожу порт через flag
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)

	// здесь использую консул, если что то будет меняться только эта строка
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// генерирую id
	instanceID := discovery.GenerateInstanceID(serviceName)
	// регистрирую сервис
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	// в отдельной голутине каждую секунду говорю, что сервис здоров
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	// если сервис сломается, отсановится, удаляю его из registry
	defer registry.Deregister(ctx, instanceID, serviceName)

	// дальше работа сервака
	repo := memory.New()
	svc := metadata.New(repo)
	h := httphandler.New(svc)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadataByID))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
