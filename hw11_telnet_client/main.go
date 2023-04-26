package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// nc -l localhost 4242
// go build -o go-telnet
// ./go-telnet localhost 4242

func main() {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGQUIT)

	var address string
	var port string
	timeout := flag.String("timeout", "10s", "timeout")
	flag.Parse()

	timeoutDuration, err := time.ParseDuration(*timeout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("address and port is empty")
		os.Exit(0)
	}
	if len(args) == 1 {
		fmt.Println("port is empty")
		os.Exit(0)
	}

	reader := bufio.NewReader(os.Stdin)
	out := &bytes.Buffer{}
	in := io.NopCloser(reader)
	address = args[0]
	port = args[1]
	tclient := NewTelnetClient(address+":"+port, timeoutDuration, in, out)
	tclient.Connect()

	wg := sync.WaitGroup{} // лучше использовать передачу по указателю
	wg.Add(2)
	go tclient.Send()
	go tclient.Receive()

	select {
	case <-interruptCh:
		tclient.Close()
		fmt.Println("...EOF")
		os.Exit(0)
	}

	wg.Wait()

}
