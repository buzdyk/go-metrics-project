package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Address string
	Report  int
	Collect int
}

var config Config

func init() {
	addrFlag := flag.String("a", "", "Address of the server that collects metrics")
	reportFlag := flag.Int("r", 0, "Report interval in seconds")
	collectFlag := flag.Int("p", 0, "Poll interval in seconds")

	flag.Parse()

	var (
		envAddr    = os.Getenv("ADDRESS")
		envReport  = os.Getenv("REPORT")
		envCollect = os.Getenv("COLLECT")
	)

	switch {
	case envAddr != "":
		config.Address = envAddr
	case *addrFlag != "":
		config.Address = *addrFlag
	default:
		config.Address = "0.0.0.0:8080"
	}

	if !strings.HasPrefix(config.Address, "http://") {
		config.Address = "http://" + config.Address
	}

	switch {
	case envCollect != "":
		v, err := strconv.Atoi(envCollect)
		if err != nil {
			panic(err)
		}
		config.Collect = v
	case *addrFlag != "":
		config.Collect = *collectFlag
	default:
		config.Collect = 2
	}

	switch {
	case envReport != "":
		v, err := strconv.Atoi(envCollect)
		if err != nil {
			panic(err)
		}
		config.Report = v
	case *addrFlag != "":
		config.Report = *reportFlag
	default:
		config.Report = 10
	}
}
