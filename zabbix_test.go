package zabbix

import (
	"bytes"
	"net"
	"testing"
	"time"
)

type netStub struct {
	response *bytes.Reader
}

func (n *netStub) Read(data []byte) (int, error) {
	return n.response.Read(data)
}

func (n *netStub) Write(data []byte) (int, error) {
	return len(data), nil
}

func TestPackageParsing(t *testing.T) {
	agent := Agent{
		Transport: &netStub{
			response: bytes.NewReader([]byte{0x5A, 0x42, 0x58, 0x44, 0x01, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74}),
		},
	}

	response, err := agent.Get("AnyKey")
	if err != nil {
		t.Error(err)
	}

	if response != "test" {
		t.Errorf("the payload does not match the expected result: %s", response)
	}
}

func TestClientServer(t *testing.T) {
	listen, err := net.Listen("tcp4", ":0")
	if err != nil {
		t.Error(err)
	}
	defer listen.Close()

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				break
			}

			go func() {
				defer conn.Close()

				agent := Agent{
					Transport: conn,
				}

				_, err := agent.GetPayload()
				if err != nil {
					t.Error(err)
				}

				agent.write([]byte{0x5A, 0x42, 0x58, 0x44, 0x01, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74})
			}()
		}
	}()

	conn, err := net.DialTimeout("tcp4", listen.Addr().String(), time.Second*2)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	agent := Agent{
		Transport: conn,
	}

	response, err := agent.Get("AnyKey")
	if err != nil {
		t.Error(err)
	}

	if response != "test" {
		t.Errorf("the payload does not match the expected result: %s", response)
	}
}
