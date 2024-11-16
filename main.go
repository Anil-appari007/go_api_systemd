package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-systemd/v22/dbus"
)

func main() {
	fmt.Println("Go...")
	serviceName := os.Args[1]
	action := os.Args[2]
	log.Printf("service=%s", serviceName)
	log.Printf("action=%s", action)

	ctx := context.TODO()
	service := serviceName + ".service"

	dCon, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		log.Fatalf("failed to connect to systemd, error = %s", err)
	}
	defer dCon.Close()

	// status, err := dCon.GetUnitPropertiesContext(ctx, service)
	// if err != nil {
	// 	log.Fatalf("error in getting status, error = %s", err)
	// }
	// fmt.Println(status)
	// fmt.Println(status["ActiveState"])
	// fmt.Println(status["LoadError"])
	// fmt.Printf("\n%T\n", status)

	// LoadError := status["LoadError"]
	// fmt.Printf("\n\n\nLoadErro of type - %T\n", LoadError)
	// LES, ok := LoadError.([]interface{})
	// if !ok {
	// 	log.Fatal("error in LoadError interface")
	// }

	// fmt.Println("LES")
	// fmt.Printf("\nlen of LES - %v", len(LES))
	// fmt.Println(LES...)
	// var LESDATA string
	// for _, each := range LES {
	// 	fmt.Printf("range - %s", each)
	// 	fmt.Printf("type - %T", each)
	// 	LESDATA += each.(string)
	// }
	// fmt.Printf("\n\nlesdata - %s", LESDATA)
	// if strings.Contains(LESDATA, "service not found") {
	// 	log.Fatal(LESDATA)
	// } else {
	// 	log.Printf("service exists %s", service)
	// }

	//

	// ok, err := IsServiceExist(dCon, service, ctx)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// if ok {
	// 	log.Printf("service %s exists", service)
	// }

	// ok, err = IsServiceActive(dCon, service, ctx)
	// if err != nil {
	// 	log.Fatalf("error while getting status - %s", err)
	// }
	// if ok {
	// 	log.Printf("service %s is active", service)
	// }

	// if data == LES {

	// }
	// fmt.Printf("data = %s", data...)

	// if status["ActiveState"] == "active" {
	// 	chan1 := make(chan<- string)
	// 	stopStatus, err := dCon.StopUnitContext(ctx, service, "replace", chan1)
	// 	if err != nil {
	// 		log.Fatalf("error in stopping service, error = %s", err)
	// 	}
	// 	fmt.Println(stopStatus)
	// }

	// if status["ActiveState"] == "inactive" {
	// 	ch2 := make(chan<- string)
	// 	startStatus, err := dCon.StartUnitContext(ctx, service, "replace", ch2)
	// 	if err != nil {
	// 		log.Fatalf("error in starting service, error = %s", err)
	// 	}
	// 	fmt.Printf("Start stauts = %d\n", startStatus)
	// }

	ok, err := IsServiceExist(dCon, service, ctx)
	if err != nil {
		log.Fatalln(err)
	}
	if ok {
		log.Printf("%s exists", service)
	} else {
		log.Fatalf("%s does not exist", service)
	}

	switch action {
	case "status":
		ok, err = IsServiceActive(dCon, service, ctx)
		if err != nil {
			log.Fatalf("error while getting status - %s", err)
		}
		if ok {
			log.Printf("service %s is active", service)
		} else {
			log.Printf("service %s is not active", service)
		}
	case "start":
		ok, err = StartService(dCon, service, ctx)
		if err != nil {
			log.Fatalf("error while getting status - %s", err)
		}
		if ok {
			log.Printf("%s is started")
		}
	case "stop":
		ok, err = StopService(dCon, service, ctx)
		if err != nil {
			log.Fatalf("error while getting status - %s", err)
		}
		if ok {
			log.Printf("%s is stopped")
		}
	}

}
