package main

import (
	"flag"
	"log"

	zmq "github.com/pebbe/zmq4"
)

var (
	ringFlag = flag.Bool("ring", false, "Request help")
)

func main() {
	flag.Parse()

	requester, err := zmq.NewSocket(zmq.REQ)
	defer requester.Close()

	if err != nil {
		log.Fatalf("Error opening the socket: %v", err)
	}

	if err = requester.Connect("tcp://polaris.membrives.fr:9101"); err != nil {
		log.Fatalf("Error connecting socket: %v", err)
	}

	var data string
	if *ringFlag {
		log.Println("Requesting help")
		data = "r"
	} else {
		log.Println("Requesting time")
		data = "t"
	}

	if _, err = requester.SendBytes([]byte(data), 0); err != nil {
		log.Fatalf("Error while sending data: %v", err)
	}

	reply, err := requester.RecvBytes(0)
	if err != nil {
		log.Fatalf("Error while receiving data: %v", err)
	}

	log.Println(string(reply))
}
