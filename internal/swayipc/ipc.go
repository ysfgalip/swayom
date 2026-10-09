package swayipc

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	IPC_HEADER_SIZE = 14
	IPC_MAGIC_SIZE  = 6
	IPC_MAGIC       = "i3-ipc"
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

func (ipc *IPC) Close() {
	ipc.socket.Close()
}

func (ipc *IPC) sendHeader(length uint32, msgtype uint32) error {
	header := make([]byte, IPC_HEADER_SIZE)
	copy(header, []byte(IPC_MAGIC))

	binary.LittleEndian.PutUint32(header[IPC_MAGIC_SIZE:], length)
	binary.LittleEndian.PutUint32(header[IPC_MAGIC_SIZE+4:], msgtype)

	_, err := ipc.socket.Write(header)
	if err != nil {
		return fmt.Errorf("Error sending Header: %w", err)
	}

	return nil
}

func (ipc *IPC) sendPayload(payload []byte) error {
	_, err := ipc.socket.Write(payload)
	if err != nil {
		return fmt.Errorf("Error sending Payload: %w", err)
	}
	return nil
}

func (ipc *IPC) readReply() (*Reply, error) {
	header := make([]byte, IPC_HEADER_SIZE)

	_, err := ipc.socket.Read(header)
	if err != nil {
		return nil, fmt.Errorf("Cannot read header from reply: %w", err)
	}

	if string(header[:6]) != IPC_MAGIC {
		return nil, fmt.Errorf("Missing magic string in reply")
	}
	payloadLen := binary.LittleEndian.Uint32(header[6:10])

	payload := make([]byte, payloadLen)

	_, err = ipc.socket.Read(payload)
	if err != nil {
		return nil, fmt.Errorf("Cannot read payload from reply: %w", err)
	}

	reply := &Reply{}
	reply.MessageType = binary.LittleEndian.Uint32(header[10:])
	reply.Payload = payload

	return reply, nil
}
