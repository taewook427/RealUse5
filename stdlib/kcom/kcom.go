// test626 : stdlib5.kcom

package kcom

import (
	"crypto/rand"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"net/http"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// go mod init example.com
// go get github.com/PuerkitoBio/goquery

// get domain txt from (http~ *.html)
func Gettxt(url string, domain string) (string, error) {
	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("get fail")
	}
	defer resp.Body.Close()

	// HTTP response check 200 OK
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("200 fail")
	}

	// HTML phrasing
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.New("parse fail")
	}

	// find p id=domain
	text := doc.Find(fmt.Sprintf("p#%s", domain)).Text()
	if text == "" {
		return "", errors.New("text find fail")
	}

	return text, nil
}

// download binary name -> path from (http~ */) + (*.num)
func Download(url string, name string, num int, path string, proc *float64) error {
	*proc = 0.0
	// URL slash update
	if url[len(url)-1] != '/' {
		url = url + "/"
	}

	// generate file
	file, err := kio.Open(path, "w")
	if err != nil {
		return errors.New("file creation fail")
	}
	defer file.Close()

	// write data on file
	for i := 0; i < num; i++ {
		*proc = float64(i) / float64(num)
		resp, err := http.Get(fmt.Sprintf("%s%s.%d", url, name, i))
		if err != nil {
			return errors.New("download fail")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.New("connection fail")
		}

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return errors.New("writing fail")
		}
	}

	*proc = 2.0
	return nil
}

// (port int, key 4B) -> address str
func Pack(port int, key []byte) string {
	temp := append(kobj.Encode(port, 2), key...)
	out := kio.Bprint(temp)
	return fmt.Sprintf("%s.%s.%s", out[0:4], out[4:8], out[8:12])
}

// address str -> (port int, key 4B)
func Unpack(address string) (int, []byte, error) {
	address = strings.Replace(strings.ToLower(address), ".", "", -1)
	out, err := kio.Bread(address)
	if err != nil {
		return 0, nil, err
	}
	return kobj.Decode(out[0:2]), out[2:6], nil
}

type node struct {
	Ipv6  bool
	Port  int
	Close int
}

// init svr/cli node
func Initcom() node {
	var out node
	out.Ipv6 = true
	out.Port = 13600
	out.Close = 150
	return out
}

// send simple data (data B, key 4B)
func (server *node) Send(data []byte, key []byte) (exit error) {
	defer func() { _ = recover() }()

	temp := make([]byte, len(data))
	for i, r := range data {
		temp[i] = byte((int(r) + int(key[i%4]) + 256) % 256)
	}
	msg := kobj.Encode(len(data), 8)
	msg = append(msg, kobj.Encode(int(crc32.ChecksumIEEE(data)), 4)...)
	msg = append(msg, temp...)

	var ipaddr string
	if server.Ipv6 {
		ipaddr = "[::1]" // use "localhost" in linux!!
	} else {
		ipaddr = "127.0.0.1"
	}
	svrsc, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ipaddr, server.Port))
	if err != nil {
		return err
	}
	if server.Close > 0 {
		go func() {
			time.Sleep(time.Duration(server.Close) * time.Second)
			exit = errors.New("timeout")
			svrsc.Close()
		}()
	}
	defer svrsc.Close()

	conn, err := svrsc.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	if server.Close > 0 {
		tvar := time.Now().Add(time.Duration(server.Close) * time.Second)
		conn.SetDeadline(tvar)
	}
	mnum := make([]byte, 8)
	conn.Read(mnum)
	if string(mnum[0:5]) == "kcom5" {
		conn.Write(mnum)
	} else {
		return errors.New("invalid connection")
	}

	conn.Write(msg)
	return nil
}

// recieve data (key 4B)
func (client *node) Recieve(key []byte) ([]byte, error) {
	var ipaddr string
	if client.Ipv6 {
		ipaddr = "[::1]" // use "localhost" in linux!!
	} else {
		ipaddr = "127.0.0.1"
	}
	clisc, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ipaddr, client.Port))
	if err != nil {
		return nil, err
	}
	defer clisc.Close()
	if client.Close > 0 {
		clisc.SetDeadline(time.Now().Add(time.Duration(client.Close) * time.Second))
	}

	mnum := make([]byte, 8)
	copy(mnum, []byte("kcom5"))
	rand.Read(mnum[5:])
	clisc.Write(mnum)
	chk := make([]byte, 8)
	clisc.Read(chk)
	for i := 0; i < 8; i++ {
		if mnum[i] != chk[i] {
			return nil, errors.New("invalid connection")
		}
	}

	temp := make([]byte, 8)
	clisc.Read(temp)
	mlen := kobj.Decode(temp)
	temp = make([]byte, 4)
	clisc.Read(temp)
	crcv := kobj.Decode(temp)
	temp = make([]byte, mlen)
	clisc.Read(temp)
	data := make([]byte, mlen)
	for i, r := range temp {
		data[i] = byte((int(r) - int(key[i%4]) + 256) % 256)
	}

	if crcv != int(crc32.ChecksumIEEE(data)) {
		return nil, errors.New("invalid key")
	}
	return data, nil
}
