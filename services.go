package main

import (
	"context"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
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
