package zabbix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"regexp"
)

const (
	FlagCommunications uint8 = 1
	FlagCompression
	FlagLargePacket
)

const protocolName string = "ZBXD"

type NetIO interface {
	Read(b []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

type Header struct {
	Protocol [len(protocolName)]byte
	Flag     uint8
	DataSize uint64
}

type Agent struct {
	Transport NetIO
}

func packageConstructor(data []byte) bytes.Buffer {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(protocolName))
	binary.Write(buffer, binary.LittleEndian, FlagCommunications)
	binary.Write(buffer, binary.LittleEndian, uint64(len(data)))
	buffer.Write(data)

	return *buffer
}

func (a *Agent) read(b []byte) (int, error) {
	return a.Transport.Read(b)
}

func (a *Agent) write(b []byte) (int, error) {
	return a.Transport.Write(b)
}

func (a *Agent) Get(key string) (string, error) {
	if a.Transport == nil {
		return "", errors.New("the transport for executing the request is not set")
	}

	buffer := packageConstructor([]byte(key))
	_, err := a.write(buffer.Bytes())
	if err != nil {
		return "", err
	}

	payload, err := a.GetPayload()
	if err != nil {
		return "", err
	}

	if regexp.MustCompile("^ZBX_NOTSUPPORTED").Match([]byte(payload)) {
		return "", errors.New(payload)
	}

	return payload, nil
}

func (a *Agent) GetPayload() (string, error) {
	header, err := a.readHeader()
	if err != nil {
		return "", err
	}

	payload, err := a.readData(header)
	if err != nil {
		return "", err
	}

	return payload, nil
}

func (a *Agent) readHeader() (Header, error) {
	header := Header{}

	_, err := a.read(header.Protocol[:])
	if err != nil {
		if err == io.EOF {
			return header, errors.New("check access restrictions in Zabbix agent configuration")
		}

		return header, err
	}

	if !bytes.Equal(header.Protocol[:], []byte(protocolName)) {
		return header, errors.New("the protocol in the response does not match Zabbix")
	}

	if binary.Read(a.Transport, binary.LittleEndian, &header.Flag) != nil {
		return header, err
	}

	if header.Flag != FlagCommunications {
		return header, errors.New("unsupported flag in agent response")
	}

	if binary.Read(a.Transport, binary.LittleEndian, &header.DataSize) != nil {
		return header, err
	}

	if header.DataSize == 0 {
		return header, errors.New("there is no payload in the response")
	}

	return header, nil
}

func (a *Agent) readData(h Header) (string, error) {
	buffer := make([]byte, 1024)
	data := bytes.Buffer{}

	for data.Len() < int(h.DataSize) {
		size, err := a.read(buffer)

		if err != nil && err != io.EOF {
			return "", err
		}

		if size == 0 {
			break
		}

		data.Write(buffer[0:size])
	}

	return data.String(), nil
}
