package manager

import (
	"fmt"
	"os"
	"path"

	"github.com/ysfgalip/swayom/internal/config"
	"github.com/ysfgalip/swayom/internal/swayipc"
)

func Manage() {
	homeDir := os.Getenv("HOME")
	configDir := path.Join(homeDir, ".config/swayom/config")
	sc := config.NewSwayConfig(configDir)

	ipcSubscribe, _ := swayipc.NewIPC(os.Getenv("SWAYSOCK"))
	ipcConfig, _ := swayipc.NewIPC(os.Getenv("SWAYSOCK"))
	ipcCmd, _ := swayipc.NewIPC(os.Getenv("SWAYSOCK"))

	if err := ipcSubscribe.SubscribeOutput(); err != nil {
		fmt.Println(err)
	}

	for {
		ok, err := ipcSubscribe.ListenEventOutput()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if ok {
			fmt.Println("output change")
			outputs, err := ipcConfig.GetOutputs()
			if err != nil {
				fmt.Println(err)
			}
			outputSlice := getOutputIDs(outputs)
			conf, err := sc.Read(outputSlice)
			if err != nil {
				fmt.Println(err)
			}
			for _, confCmd := range conf {
				err = ipcCmd.RunCommand(confCmd)
				if err != nil {
					fmt.Println(err)
				}
			}
			ipcSubscribe.Close()
			ipcSubscribe, _ = swayipc.NewIPC(os.Getenv("SWAYSOCK"))
			if err := ipcSubscribe.SubscribeOutput(); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func getOutputIDs(outputs []*swayipc.Output) (ids []string) {
	for _, output := range outputs {
		if output.Active {
			ids = append(ids, output.Make+" "+output.Model+" "+output.Serial)
		}
	}

	return ids
}
