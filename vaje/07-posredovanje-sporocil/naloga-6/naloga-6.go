package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

type log struct {
	retries int
	success bool
}

var N int
var id int
var root int

var logs []log
var stop bool
var mu sync.Mutex
var msg chan byte

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func receive(addr *net.UDPAddr) {
	// Poslušamo
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	fmt.Println("Proces", id, "posluša na", addr)
	buffer := make([]byte, 1024)
	if id == root {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		for !stop {
			// Preberemo sporočilo
			mLen, err := conn.Read(buffer)
			if err == nil {
				fmt.Println("Proces", id, "prejel sporočilo:", buffer[:mLen])
				idr := byte(buffer[0])
				mu.Lock()
				logs[idr].success = true
				mu.Unlock()
			}
		}
		conn.Close()
		close(msg)
	} else {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		mLen, err := conn.Read(buffer)
		checkError(err)
		fmt.Println("Proces", id, "prejel sporočilo:", string(buffer[:mLen]))
		idr := byte(buffer[0])
		conn.Close()
		msg <- idr
	}
}

func send(addr *net.UDPAddr) {
	// Odpremo povezavo
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()
	// Pripravimo sporočilo
	sMsg := [1]byte{byte(id)}
	_, err = conn.Write(sMsg[:])
	checkError(err)
	fmt.Println("Proces", id, "poslal sporočilo", sMsg, "procesu na naslovu", addr)
}

func broadcast(remoteAddrs []*net.UDPAddr) {
	for {
		count := 0
		for i := 0; i < N; i++ {
			if i == root {
				continue
			}
			mu.Lock()
			if logs[i].success || logs[i].retries >= 5 {
				count++
				mu.Unlock()
				continue
			}
			mu.Unlock()
			if count == N-1 {
				fmt.Println("Vsi procesi so prejeli sporočilo ali dosegli maksimalno število ponovitev.")
				stop = true
				return
			}
			logs[i].retries++
			send(remoteAddrs[i])
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// Preberi argumente
	portPtr := flag.Int("p", 9000, "# start port")
	idPtr := flag.Int("id", 0, "# process id")
	NPtr := flag.Int("N", 2, "total number of processes")
	rootPtr := flag.Int("root", 0, "# root id")
	flag.Parse()

	rootPort := *portPtr
	id = *idPtr
	N = *NPtr
	root = *rootPtr
	basePort := rootPort + id

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", basePort))
	checkError(err)
	go receive(localAddr)
	if id == root {
		remoteAddrs := make([]*net.UDPAddr, N)
		logs = make([]log, N)
		for i := 1; i < N; i++ {
			remotePort := rootPort + i
			remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", remotePort))
			checkError(err)
			remoteAddrs[i] = remoteAddr
		}
		<-msg
	} else {
		<-msg
		remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", rootPort+root))
		checkError(err)
		send(remoteAddr)
	}
}
