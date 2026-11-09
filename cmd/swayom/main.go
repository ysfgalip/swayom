package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/ysfgalip/swayom/internal/manager"
)

func main() {
	homeDir := os.Getenv("HOME")
	defaultConfigDir := path.Join(homeDir, ".config/swayom/config")
	configDir := flag.String("config", defaultConfigDir, "specify a custom config file")
	flag.Parse()

	err := manager.Manage(*configDir)
	fmt.Fprintln(os.Stderr, err)
}
