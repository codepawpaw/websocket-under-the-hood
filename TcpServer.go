package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}

}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		return i + 2, dropCR(data[0:i]), nil
	}

	if atEOF {
		return len(data), dropCR(data), nil
	}

	return 0, nil, nil
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

	for {
		reader := bufio.NewReader(c)
		scanner := bufio.NewScanner(reader)
		secKey := ""

		scanner.Split(ScanCRLF)

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "Sec-WebSocket-Key") {
				key := strings.Split(scanner.Text(), ": ")[1]

				h := sha1.New()
				h.Write([]byte(key))
				h.Write(keyGUID)
				secKey = base64.StdEncoding.EncodeToString(h.Sum(nil))
				break
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Invalid input: %s", err)
		}

		upgrade := "HTTP/1.1 101 Switching Protocols\r\n" + "Upgrade: websocket\r\n" + "Connection: Upgrade\r\n" + "Sec-WebSocket-Accept: " + secKey + "\r\n\r\n"
		println(upgrade)
		c.Write([]byte(upgrade))
	}
	c.Close()
}
