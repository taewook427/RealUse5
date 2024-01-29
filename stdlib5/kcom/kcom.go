package kcom

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// test582 : kcom (go)

// little endian encoding
func encode(num int, length int) *[]byte {
	temp := make([]byte, length)
	for i := 0; i < length; i++ {
		temp[i] = byte(num % 256)
		num = num / 256
	}
	return &temp
}

// little endian decoding
func decode(data *[]byte) int {
	temp := 0
	for i, r := range *data {
		if r != 0 {
			exp := 1
			for j := 0; j < i; j++ {
				exp = exp * 256
			}
			temp = temp + int(r)*exp
		}
	}
	return temp
}

// port int + key bytes -> address str
func Pack(port int, key []byte) string {
	if 48%len(key) != 0 {
		panic("invalid keylen")
	}
	temp := make([]string, len(key))
	num := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i, r := range key {
		temp[i] = num[r/16] + num[r%16]
	}
	out := fmt.Sprintf("%d.%s", port, strings.Join(temp, "."))
	return out
}

// address str -> port int + key 48B
func Unpack(address string) (int, []byte) {
	address = strings.ToLower(address)
	temp := strings.Split(address, ".")
	port, _ := strconv.Atoi(temp[0])
	temp = temp[1:]
	if 48%len(temp) != 0 {
		panic("invalid keylen")
	}
	key := make([]byte, len(temp))
	num := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i, r := range temp {
		j := 0
		for r[0] != num[j][0] {
			j++
		}
		k := 0
		for r[1] != num[k][0] {
			k++
		}
		key[i] = byte(16*j + k)
	}
	return port, bytes.Repeat(key, 48/len(key))
}

type server struct {
	Ipv6  bool
	Port  int
	Close int
	Msg   string
}

func Initsvr() server {
	var temp server
	temp.Ipv6 = true
	temp.Port = 13600
	temp.Close = 150
	temp.Msg = ""
	return temp
}

// transmit data
func (self *server) Send(data []byte) {
	// 서버 접속
	var ipaddr string
	if self.Ipv6 {
		ipaddr = "[::1]" // linux에서는 "localhost"
	} else {
		ipaddr = "127.0.0.1"
	}
	svrsc, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ipaddr, self.Port))
	if err != nil {
		panic(err)
	}
	if self.Close > 0 {
		go func() {
			time.Sleep(time.Duration(self.Close) * time.Second)
			panic("timeout")
		}()
	}
	defer svrsc.Close()

	// KCOM5 + rand 3B echo
	conn, err := svrsc.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if self.Close > 0 {
		tvar := time.Now().Add(time.Duration(self.Close) * time.Second)
		conn.SetDeadline(tvar)
	}
	mnum := make([]byte, 8)
	conn.Read(mnum)
	if string(mnum[0:5]) == "KCOM5" {
		conn.Write(mnum)
	} else {
		panic("invalid connection")
	}

	// data transmit
	localmsg := []byte(self.Msg)
	buf := *encode(len(localmsg), 4)
	conn.Write(append(buf, localmsg...))
	buf = *encode(len(data), 8)
	conn.Write(append(buf, data...))
}

type client struct {
	Ipv6  bool
	Port  int
	Close int
	Msg   string
}

func Initcli() client {
	var temp client
	temp.Ipv6 = true
	temp.Port = 13600
	temp.Close = 150
	temp.Msg = ""
	return temp
}

// recieve data
func (self *client) Recieve() []byte {
	// 서버 접속
	var ipaddr string
	if self.Ipv6 {
		ipaddr = "[::1]"
	} else {
		ipaddr = "127.0.0.1"
	}
	clisc, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ipaddr, self.Port))
	if err != nil {
		panic(err)
	}
	defer clisc.Close()
	if self.Close > 0 {
		clisc.SetDeadline(time.Now().Add(time.Duration(self.Close) * time.Second))
	}

	// KCOM5 + rand 3B echo
	mnum := make([]byte, 8)
	copy(mnum, []byte("KCOM5"))
	rand.Seed(time.Now().UnixNano())
	rand.Read(mnum[5:])
	clisc.Write(mnum)
	chk := make([]byte, 8)
	clisc.Read(chk)
	for i := 0; i < 8; i++ {
		if mnum[i] != chk[i] {
			panic("invalid connection")
		}
	}

	// data recieve
	temp := make([]byte, 4)
	clisc.Read(temp)
	mlen := decode(&temp)
	localmsg := make([]byte, mlen)
	clisc.Read(localmsg)
	temp = make([]byte, 8)
	clisc.Read(temp)
	dlen := decode(&temp)
	data := make([]byte, dlen)
	clisc.Read(data)

	// socket close
	self.Msg = string(localmsg)
	return data
}
