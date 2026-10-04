package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ysfgalip/swayom/internal/swayipc"
)

func main() {
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
