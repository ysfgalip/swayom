// Package swayipc handles the ipc with the sway instance
package swayipc

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	IpcHeaderSize = 14
	IpcMagicSize  = 6
	IpcMagic      = "i3-ipc"
)

type Reply struct {
	MessageType uint32
	Payload     []byte
}

type IPC struct {
	socket net.Conn
}

func NewIPC(socketFile string) (*IPC, error) {
	ipc := &IPC{}

	sock, err := net.Dial("unix", socketFile)
	if err != nil {
		return nil, err
	}
	ipc.socket = sock
	return ipc, nil
}

func (ipc *IPC) Close() error {
	if err := ipc.socket.Close(); err != nil {
		return err
	}

	return nil
}

func (ipc *IPC) sendHeader(length uint32, msgtype uint32) error {
	header := make([]byte, IpcHeaderSize)
	copy(header, []byte(IpcMagic))

	binary.LittleEndian.PutUint32(header[IpcMagicSize:], length)
	binary.LittleEndian.PutUint32(header[IpcMagicSize+4:], msgtype)

	_, err := ipc.socket.Write(header)
	if err != nil {
		return fmt.Errorf("error sending Header: %w", err)
	}

	return nil
}

func (ipc *IPC) sendPayload(payload []byte) error {
	_, err := ipc.socket.Write(payload)
	if err != nil {
		return fmt.Errorf("error sending Payload: %w", err)
	}
	return nil
}

func (ipc *IPC) readReply() (*Reply, error) {
	header := make([]byte, IpcHeaderSize)

	_, err := ipc.socket.Read(header)
	if err != nil {
		return nil, fmt.Errorf("cannot read header from reply: %w", err)
	}

	if string(header[:6]) != IpcMagic {
		return nil, fmt.Errorf("missing magic string in reply")
	}
	payloadLen := binary.LittleEndian.Uint32(header[6:10])

	payload := make([]byte, payloadLen)

	_, err = ipc.socket.Read(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot read payload from reply: %w", err)
	}

	reply := &Reply{}
	reply.MessageType = binary.LittleEndian.Uint32(header[10:])
	reply.Payload = payload

	return reply, nil
}
