package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-openssl"
)

var dialer = (&net.Dialer{
	DualStack: true,
	Timeout:   7 * time.Second,
}).Dial

var listen = flag.String(`listen`, `:8043`, `Listen address. Eg: :8443; unix:/tmp/proxy.sock`)
var addr = flag.String(`addr`, `megadomain.vnn.vn:443`, `Remote address`)
var sni = flag.String(`sni`, `megadomain.vnn.vn`, `TLS Hello SNI String`)
var sslVer = flag.Int(`ver`, 0x02, `SSL Version. SSLv3 = 2; TLSv1 = 3; TLSv1.1 = 4; TLSv1.2 = 5; AnyVer = 6`)
var noVerify = flag.Bool(`k`, false, `Don't verify SNI`)

func TlsDial() (net.Conn, error) {
	conn, err := dialer("tcp", *addr)
	if err != nil {
		return nil, err
	}
	ctx, err := openssl.NewCtxWithVersion(openssl.SSLVersion(*sslVer))
	if err != nil {
		conn.Close()
		return nil, err
	}
	tlsConn, err := openssl.Client(conn, ctx)
	if err != nil {
		conn.Close()
		return nil, err
	}
	err = tlsConn.SetTlsExtHostName(*sni)
	if err != nil {
		tlsConn.Close()
		return nil, err
	}
	err = tlsConn.Handshake()
	if err != nil {
		tlsConn.Close()
		return nil, err
	}
	if *noVerify == false {
		err = tlsConn.VerifyHostname(*sni)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	return tlsConn, err
}

func handleConn(l net.Conn) {
	defer l.Close()
	r, err := TlsDial()
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Close()
	go io.Copy(l, r)
	io.Copy(r, l)
}

func main() {
	flag.Parse()
	var err error
	var ln net.Listener
	if strings.HasPrefix(*listen, `unix:`) {
		unixFile := (*listen)[5:]
		os.Remove(unixFile)
		ln, err = net.Listen(`unix`, unixFile)
		os.Chmod(unixFile, os.ModePerm)
		log.Println(`Listening:`, unixFile)
	} else {
		ln, err = net.Listen(`tcp`, *listen)
		log.Println(`Listening:`, ln.Addr().String())
	}
	if err != nil {
		log.Panicln(err)
	}
	for {
		l, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(l)
	}
}
