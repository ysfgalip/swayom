package swayipc

import (
	"encoding/json"
	"fmt"
)

const (
	RunCommandNr = uint32(0)
)

type CommandReply struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (ipc *IPC) RunCommand(command string) error {
	payloadLen := len(command)

	err := ipc.sendHeader(uint32(payloadLen), RunCommandNr)
	if err != nil {
		return err
	}

	err = ipc.sendPayload([]byte(command))
	if err != nil {
		return err
	}

	reply, err := ipc.readReply()
	if err != nil {
		return fmt.Errorf("error reading command reply: %w", err)
	}

	var replyPayload CommandReply
	err = json.Unmarshal(reply.Payload, &replyPayload)
	if err != nil {
		return fmt.Errorf("error unmarshalling command reply: %w", err)
	}

	if !replyPayload.Success {
		return fmt.Errorf("command unsuccessful: %s", replyPayload.Error)
	}

	return nil
}
