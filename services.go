package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/gin-gonic/gin"
)

func IsServiceExist(conn *dbus.Conn, serviceName string, ctx context.Context) (bool, error) {
	serviceData, err := conn.GetUnitPropertiesContext(ctx, serviceName)
	if err != nil {
		return false, err
	}
	LoadError := serviceData["LoadError"]
	LoadError2 := LoadError.([]interface{})
	var LoadError3 string
	for _, each := range LoadError2 {
		LoadError3 += each.(string)
	}
	if strings.Contains(LoadError3, "service not found") {
		return false, nil
	}
	return true, nil
}

func IsServiceActive(conn *dbus.Conn, serviceName string, ctx context.Context) (bool, error) {
	serviceData, err := conn.GetUnitPropertiesContext(ctx, serviceName)
	if err != nil {
		return false, err
	}

	if serviceData["ActiveState"] != "active" {
		return false, nil
	}
	return true, nil
}

func StopService(conn *dbus.Conn, serviceName string, ctx context.Context) (bool, error) {
	ch := make(chan<- string)
	_, err := conn.StopUnitContext(ctx, serviceName, "replace", ch)
	if err != nil {
		return false, err
	}
	return true, nil
}

func StartService(conn *dbus.Conn, serviceName string, ctx context.Context) (bool, error) {
	ch := make(chan<- string)
	_, err := conn.StartUnitContext(ctx, serviceName, "replace", ch)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ServiceHandler(c *gin.Context) {

	type details struct {
		ServiceName string `json:"serviceName"`
		Action      string `json:"action"`
	}
	var sd details
	if err := c.BindJSON(&sd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	log.Printf("service=%s", sd.ServiceName)
	log.Printf("action=%s", sd.Action)

	ctx := context.TODO()
	service := sd.ServiceName + ".service"

	dCon, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
		// log.Fatalf("failed to connect to systemd, error = %s", err)
	}
	defer dCon.Close()

	ok, err := IsServiceExist(dCon, service, ctx)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})

		return

	}
	if ok {
		log.Printf("%s exists", service)
	} else {
		// log.Fatalf("%s does not exist", service)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s does not exist", service)})
		return

	}

	switch sd.Action {
	case "status":
		ok, err = IsServiceActive(dCon, service, ctx)
		if err != nil {
			// log.Fatalf("error while getting status - %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return

		}
		if ok {
			log.Printf("service %s is active", service)
			c.JSON(200, gin.H{"msg": fmt.Sprintf("service %s is active", service)})
			return

		} else {
			log.Printf("service %s is not active", service)
			c.JSON(200, gin.H{"msg": fmt.Sprintf("service %s is not active", service)})
			return

		}
	case "start":
		ok, err = StartService(dCon, service, ctx)
		if err != nil {
			log.Printf("error while getting status - %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return

		}
		if ok {
			log.Printf("%s is started", service)
			c.JSON(200, gin.H{"msg": fmt.Sprintf("%s is started", service)})
			return

		}
	case "stop":
		ok, err = StopService(dCon, service, ctx)
		if err != nil {
			log.Printf("error while getting status - %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return

		}
		if ok {
			log.Printf("%s is stopped", service)
			c.JSON(200, gin.H{"msg": fmt.Sprintf("%s is stopped", service)})
			return

		}
	}
}

func RunServiceAction(serviceName, action string) (int, string, error) {
	ctx := context.TODO()
	service := serviceName + ".service"
	dCon, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	defer dCon.Close()

	ok, err := IsServiceExist(dCon, service, ctx)
	if err != nil {
		// log.Println(err)
		return http.StatusInternalServerError, "", err
	}
	if !ok {
		// log.Fatalf("%s does not exist", service)
		return http.StatusNotFound, "", fmt.Errorf("%s does not exist", service)
	}

	switch action {
	case "status":
		ok, err = IsServiceActive(dCon, service, ctx)
		if err != nil {
			// log.Fatalf("error while getting status - %s", err)
			return http.StatusInternalServerError, "", err
		}
		if ok {
			// log.Printf("service %s is active", service)
			return http.StatusAccepted, fmt.Sprintf("service %s is active", service), nil

		} else {
			// log.Printf("service %s is not active", service)
			return http.StatusAccepted, fmt.Sprintf("service %s is not active", service), nil

		}
	case "start":
		ok, err = StartService(dCon, service, ctx)
		if err != nil {
			// log.Printf("error while getting status - %s", err)
			return http.StatusInternalServerError, "", err
		}

		if ok {
			// log.Printf("%s is started", service)
			return http.StatusAccepted, fmt.Sprintf("%s is started", service), nil

		}
	case "stop":
		ok, err = StopService(dCon, service, ctx)
		if err != nil {
			// log.Printf("error while getting status - %s", err)
			return http.StatusInternalServerError, "", err
		}
		if ok {
			return http.StatusAccepted, fmt.Sprintf("%s is stopped", service), nil
		}
	}
	return http.StatusBadRequest, "", fmt.Errorf("invalid action : %s", action)
}

func multiServiceHandler(c *gin.Context) {
	type details struct {
		ServiceName string `json:"serviceName"`
		Action      string `json:"action"`
		Code        int    `json:"code"`
		Msg         string `json:"msg"`
		Err         string `json:"err"`
	}
	var multiSd []details
	if err := c.BindJSON(&multiSd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for each := range multiSd {
		// service := multiSd[each]
		// service.Code, service.Msg, service.Err = RunServiceAction(service.ServiceName, service.Action)
		var err error
		multiSd[each].Code, multiSd[each].Msg, err = RunServiceAction(multiSd[each].ServiceName, multiSd[each].Action)

		if err != nil {
			multiSd[each].Err = err.Error()
		}

		log.Println(multiSd[each])
	}
	// multiSd[0].Err = fmt.Errorf("test")
	log.Println(multiSd)
	c.JSON(200, gin.H{"msg": multiSd})
}
