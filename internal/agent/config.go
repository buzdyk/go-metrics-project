package agent

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Address string
	Report  int
	Collect int
}

func NewConfig() Config {
	return Config{"0.0.0.0:8080", 10, 2}
}

func NewConfigFromCLI() *Config {
	flag.Usage = usage

	addrFlag := flag.String("a", "", "Address of the server that collects metrics")
	reportFlag := flag.Int("r", 0, "Report interval in seconds")
	collectFlag := flag.Int("p", 0, "Poll interval in seconds")

	flag.Parse()

	config := NewConfig()

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
	}

	return &config
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("version 0.1")
}
