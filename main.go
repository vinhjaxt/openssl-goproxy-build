package main

import (
	"log"
	"net"
	"time"

	"github.com/libp2p/go-openssl"
)

var dialer = (&net.Dialer{
	DualStack: true,
	Timeout:   7 * time.Second,
}).Dial

func TlsDial(hostname, addr string) (net.Conn, error) {
	conn, err := dialer("tcp", addr)
	if err != nil {
		return nil, err
	}
	ctx, err := openssl.NewCtxWithVersion(openssl.SSLv3)
	if err != nil {
		conn.Close()
		return nil, err
	}
	tlsConn, err := openssl.Client(conn, ctx)
	if err != nil {
		conn.Close()
		return nil, err
	}
	err = tlsConn.SetTlsExtHostName(hostname)
	if err != nil {
		tlsConn.Close()
		return nil, err
	}
	err = tlsConn.Handshake()
	if err != nil {
		tlsConn.Close()
		return nil, err
	}
	err = tlsConn.VerifyHostname(hostname)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return tlsConn, err
}

func main() {
	conn, err := TlsDial("megadomain.vnn.vn", "megadomain.vnn.vn:443")
	if err != nil {
		log.Panicln(err)
	}
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: megadomain.vnn.vn\r\nUser-Agent: -\r\nAccept: */*\r\nConnection: close\r\n\r\n"))
	if err != nil {
		log.Panicln(err)
	}

	buf := make([]byte, 4096)
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		data := buf[:n]
		log.Println(string(data))
	}
	conn.Close()
}
