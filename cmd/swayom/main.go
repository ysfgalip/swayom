package main

import (
	"flag"
	"os"
	"path"

	"github.com/ysfgalip/swayom/internal/manager"
)

func main() {
	homeDir := os.Getenv("HOME")
	defaultConfigDir := path.Join(homeDir, ".config/swayom/config")
	configDir := flag.String("config", defaultConfigDir, "specify a custom config file")
	flag.Parse()

	manager.Manage(*configDir)
}
