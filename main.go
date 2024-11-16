package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// fmt.Println("Go...")
	// serviceName := os.Args[1]
	// action := os.Args[2]
	// log.Printf("service=%s", serviceName)
	// log.Printf("action=%s", action)

	// ctx := context.TODO()
	// service := serviceName + ".service"

	// dCon, err := dbus.NewSystemdConnectionContext(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to connect to systemd, error = %s", err)
	// }
	// defer dCon.Close()

	// ok, err := IsServiceExist(dCon, service, ctx)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// if ok {
	// 	log.Printf("%s exists", service)
	// } else {
	// 	log.Fatalf("%s does not exist", service)
	// }

	// switch action {
	// case "status":
	// 	ok, err = IsServiceActive(dCon, service, ctx)
	// 	if err != nil {
	// 		log.Fatalf("error while getting status - %s", err)
	// 	}
	// 	if ok {
	// 		log.Printf("service %s is active", service)
	// 	} else {
	// 		log.Printf("service %s is not active", service)
	// 	}
	// case "start":
	// 	ok, err = StartService(dCon, service, ctx)
	// 	if err != nil {
	// 		log.Fatalf("error while getting status - %s", err)
	// 	}
	// 	if ok {
	// 		log.Printf("%s is started", service)
	// 	}
	// case "stop":
	// 	ok, err = StopService(dCon, service, ctx)
	// 	if err != nil {
	// 		log.Fatalf("error while getting status - %s", err)
	// 	}
	// 	if ok {
	// 		log.Printf("%s is stopped", service)
	// 	}
	// }

	router := gin.Default()
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "running"})
	})
	router.POST("/service", ServiceHandler)
	// router.Run("localhost:8080")
	go func() {
		router.Run("localhost:8080")

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server")

}
