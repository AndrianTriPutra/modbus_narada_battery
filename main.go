package main

import (
	"atp/modbus/narada"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	flag.Usage = func() {
		log.Printf("Usage: go run . port packID")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	port := flag.Args()[0]
	pack := flag.Args()[1]
	packID, err := strconv.Atoi(pack)
	if err != nil {
		log.Fatalf("failed parse pack -> %s", err.Error())
	}

	setting := narada.Setting{
		Port:     port,
		Baudrate: 9600,
		Timeout:  3 * time.Second,
	}

	ctx := context.Background()
	repo := narada.NewRepository(setting)

	data, err := repo.Modbus(ctx, uint8(packID))
	if err != nil {
		log.Fatalf("Narada:%s", err.Error())
	}

	js, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		log.Fatalf("json:%s", err.Error())
	}
	msg := string(js)
	log.Printf("payload:\n%s", msg)
}
