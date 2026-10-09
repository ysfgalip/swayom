package swayipc

import (
	"encoding/json"
	"fmt"
)

const (
	SubscribeNr    = uint32(2)
	EventOutput    = uint32(0x80000001)
	EventOutputStr = "output"
	GetOutputsNr   = uint32(3)
)

type Output struct {
	Name   string `json:"name"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Serial string `json:"serial"`
	Active bool   `json:"active"`
}

func (ipc *IPC) GetOutputs() ([]*Output, error) {
	err := ipc.sendHeader(0, GetOutputsNr)
	if err != nil {
		return nil, err
	}

	reply, err := ipc.readReply()
	if err != nil {
		return nil, err
	}

	if reply.MessageType != GetOutputsNr {
		return nil, fmt.Errorf("wrong reply message")
	}

	displays := []*Output{}
	if err := json.Unmarshal(reply.Payload, &displays); err != nil {
		return nil, fmt.Errorf("error unmarshalling the available displays: %w", err)
	}

	return displays, nil
}

func (ipc *IPC) SubscribeOutput() error {
	payload, err := json.Marshal([]string{EventOutputStr})
	if err != nil {
		return fmt.Errorf("error marshalling event: %w", err)
	}

	err = ipc.sendHeader(uint32(len(payload)), SubscribeNr)
	if err != nil {
		return err
	}

	err = ipc.sendPayload(payload)
	if err != nil {
		return err
	}

	reply, err := ipc.readReply()
	if err != nil {
		return err
	}

	if reply.MessageType != SubscribeNr {
		return fmt.Errorf("wrong reply message")
	}

	var replyPayload map[string]any
	err = json.Unmarshal(reply.Payload, &replyPayload)
	if err != nil {
		return fmt.Errorf("cannot unmarshal the subscription reply: %w", err)
	}

	if !replyPayload["success"].(bool) {
		return fmt.Errorf("cannot subscribe to the event")
	}
	return nil
}

func (ipc *IPC) ListenEventOutput() (bool, error) {
	reply, err := ipc.readReply()
	if err != nil {
		return false, fmt.Errorf("error reading the event message")
	}

	if reply.MessageType != EventOutput {
		return false, nil
	}

	return true, nil
}
