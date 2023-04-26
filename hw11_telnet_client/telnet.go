package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &SimpleTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.

type SimpleTelnetClient struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (stc *SimpleTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", stc.address, stc.timeout) // todo err
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//defer conn.Close()

	stc.conn = conn
	return nil
}

func (stc *SimpleTelnetClient) Close() error {
	return stc.conn.Close()
}

func (stc *SimpleTelnetClient) Send() error {
	sc := bufio.NewScanner(stc.in)
	for sc.Scan() {
		writeData := sc.Bytes()
		writeData = append(writeData, "\n"...)
		if _, err := stc.conn.Write(writeData); err != nil {
			fmt.Println("error write: ", err)
			return nil
		}
	}
	return nil
}

func (stc *SimpleTelnetClient) Receive() error {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGQUIT)
	r := bufio.NewReader(stc.conn)
	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error write: ", err)
			return err
		}
		_, err = stc.out.Write(data)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	}
	return nil
}
