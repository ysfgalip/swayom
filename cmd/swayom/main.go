package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ysfgalip/swayom/internal/config"
	"github.com/ysfgalip/swayom/internal/swayipc"
)

func main() {
	homeDir := os.Getenv("HOME")
	configDir := path.Join(homeDir, ".config/swayom/config")
	sc := config.NewSwayConfig(configDir)
	conf, err := sc.Read()
	if err != nil {
		fmt.Println(err)
	}
	for _, e := range conf {
		fmt.Println(e)
	}

	ipc, _ := swayipc.NewIPC(os.Getenv("SWAYSOCK"))

	ipc.SubscribeOutput()

	for {
		ok, err := ipc.ListenEventOutput()
		if err != nil {
			fmt.Println(err)
		}

		if ok {
			outputs, err := ipc.GetOutputs()
			if err != nil {
				fmt.Println(err)
			}
			jsonOut, err := json.Marshal(&outputs)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(jsonOut))
		}
	}
}
