// Package manager implements the management of the configuration and the ipc implementation
package manager

import (
	"errors"
	"fmt"
	"os"

	"github.com/ysfgalip/swayom/internal/config"
	"github.com/ysfgalip/swayom/internal/swayipc"
)

func Manage(configDir string) error {
	sc := config.NewSwayConfig(configDir)
	swaySocket := os.Getenv("SWAYSOCK")
	var errs []error

	ipcSubscribe, err := swayipc.NewIPC(swaySocket)
	errs = append(errs, err)
	ipcConfig, err := swayipc.NewIPC(swaySocket)
	errs = append(errs, err)
	ipcCmd, err := swayipc.NewIPC(swaySocket)
	errs = append(errs, err)

	if err = errors.Join(errs...); err != nil {
		return err
	}

	if err := ipcSubscribe.SubscribeOutput(); err != nil {
		return err
	}

	ok := true
	for {
		if ok {
			outputs, err := ipcConfig.GetOutputs()
			if err != nil {
				return err
			}
			outputSlice := getOutputIDs(outputs)
			conf, err := sc.Read(outputSlice)
			if err != nil {
				return err
			}
			for _, confCmd := range conf {
				err = ipcCmd.RunCommand(confCmd)
				if err != nil {
					return err
				}
			}
			if err = ipcSubscribe.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			ipcSubscribe, _ = swayipc.NewIPC(swaySocket)
			if err = ipcSubscribe.SubscribeOutput(); err != nil {
				return err
			}

			ok, err = ipcSubscribe.ListenEventOutput()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
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
