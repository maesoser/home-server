package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"sync"
)

const (
	writeWait = 10 * time.Second
	maxMessageSize = 2048
	pongWait = 100 * time.Second
	pingPeriod = pongWait / 4
	closeGracePeriod = pongWait * 4
)

type SockClient struct {
	Address  string
	CertPath string
	Verbose bool
	r io.Reader
	w io.Writer
	ws *websocket.Conn
	wg sync.WaitGroup
}

/*
Ping sends a ping every X time. If it receives a signal usign channel it stops.
*/
func (client *SockClient) ping(done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			vlog(client.Verbose, "\rSending ping to %s\n", client.Address)
			if err := client.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Printf("ping: %s\n", err)
			}
		case <-done:
			return
		}
	}
}

func (client *SockClient) pong(string) error {
	vlog(client.Verbose, "\rPong received\n")
	client.ws.SetReadDeadline(time.Now().Add(pongWait))
	return nil 
}

func (client *SockClient) toStdout(){
	defer client.wg.Done()
	for {
		mt, r, err := client.ws.NextReader()
		if err != nil {
			log.Printf("read: %v\n", err)
		}
		if mt == websocket.BinaryMessage{
			writer(r, os.Stdout, maxMessageSize)
		}
		// We read from the endpoint and write to the ws
	}
}

func (client *SockClient) fromStdin(){
	defer client.wg.Done()
	for {
		w, err := client.ws.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Printf("write: %v\n", err)
		}
		writer(os.Stdin, w, maxMessageSize)
	}
}

func vlog(enabled bool, format string, v ...interface{}){
	if (enabled){
		log.Printf(format, v...)
	}
}

func (client *SockClient) Close(){
	client.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	client.ws.Close()
}

func (client *SockClient) Run() {
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	var caCertPool *x509.CertPool
	if client.CertPath != ""{
		vlog(client.Verbose, "Using certificate %s\n", client.CertPath)
		caCert, err := ioutil.ReadFile(client.CertPath)
		if err != nil {
			log.Fatalf("Error reading file %s : %s\n", client.CertPath, err)
		}
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("Error creating cert pool : %s\n", err)
		}
		caCertPool.AppendCertsFromPEM(caCert)
		vlog(client.Verbose, "Certificate succesfully loaded\n")
	}

	vlog(client.Verbose, "Connecting to %s\n", client.Address)
	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}
	var err error
	client.ws, _, err = dialer.Dial("ws://" + client.Address, nil)
	if err != nil {
		log.Fatalf("Error dialing to ws://%s (%s)\n", client.Address, err)
	}
	defer client.ws.Close()
	vlog(client.Verbose, "Succesfully connected to %s\n", client.Address)

	stdoutDone := make(chan struct{})
	client.wg.Add(3)

	client.ws.SetReadLimit(maxMessageSize)
	client.ws.SetReadDeadline(time.Now().Add(pongWait))
	client.ws.SetPongHandler(client.pong)

	go client.ping(stdoutDone)
	go client.toStdout()
	go client.fromStdin()
	vlog(client.Verbose, "Ready")
	client.wg.Wait()
	client.Close()
}

func writer(r io.Reader, to io.Writer, bufSize int) {
	b := make([]byte, bufSize)
	n, err := r.Read(b)
	if err != nil {
		return
	}
	if n > 0 {
		to.Write(b[0:n])
	}
}

// Server code
type SockServer struct {
	Verbose bool
	Target  string
	wg sync.WaitGroup
	ws *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     (func(r *http.Request) bool { return true }),
	EnableCompression: true,
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

func (server *SockServer) ping(string) error {
	vlog(server.Verbose, "Ping received.\n")
	//if err := server.ws.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
	//	log.Printf("pong reponse: %s\n", err)
	//}
	return nil 
}

func (server *SockServer) close(i int, t string) error {
	vlog(server.Verbose, "Receiving CLOSE %d %s\n",i,t)
	return nil 
}

func (server *SockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var endConn *net.TCPConn
	var err error
	var addr *net.TCPAddr

	log.Printf("Received request from %s (%s)\n", r.Host, r.RemoteAddr)
	server.ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %s\n", err)
		return
	}
	defer server.ws.Close()
	vlog(server.Verbose, "Connection from %s succesfully upgraded\n", r.Host)

	if addr, err = net.ResolveTCPAddr("tcp4", server.Target); err != nil {
		fmt.Printf("Error resolving proxy address: %s\n%s", server.Target, err)
	}

	if endConn, err = net.DialTCP("tcp", nil, addr); err != nil {
		fmt.Printf("Error forwarding connection to: %s\n%s", server.Target,err)
	}
	vlog(server.Verbose, "Forwarding connection to %s\n", server.Target)
	bufReader := bufio.NewReader(endConn)

	server.ws.SetReadLimit(maxMessageSize)
	server.ws.SetPingHandler(server.ping)
	//server.ws.SetCloseHandler(server.close)
	
	server.wg.Add(2)
	go func() {
		defer server.wg.Done()
		for {
			mt, r, err := server.ws.NextReader()
			if err != nil {
				log.Printf("read err: %v\n", err)
				return
			}
			if mt != 2{
				vlog(server.Verbose, "RCVD MSG %d\n", mt)
			}
			writer(r, endConn, maxMessageSize)
		}
	}()
	func() {
		defer server.wg.Done()
		for {
			w, err := server.ws.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Printf("write err: %v\n", err)
				return
			}
			writer(bufReader, w, maxMessageSize)
		}
	}()
	server.wg.Wait()
	server.ws.Close()
}

func GetEnvStr(name, value string) string {
    if os.Getenv(name) != "" {
        return os.Getenv(name)
    }
    return value
}

/*
Client usage:
	wstun -client <addr>:<port> -cert <cert>
Server usage:
	wstun -server <port> -dst <addr>:<port>
*/

func main() {
	// sudo setcap 'cap_net_bind_service=+ep' pasarela
	clientAddr := flag.String("client", GetEnvStr("CLIENT_ADDR", ""), "WebSocket target address")
	clientCert := flag.String("cert", GetEnvStr("CLIENT_CERT","client-ca.pem"), "Client certificate for Access")
	serverPort := flag.String("server", GetEnvStr("SERVER_PORT","8080"), "Server port")
	targetAddr := flag.String("dest", GetEnvStr("SERVER_TARGET","127.0.0.1:22"), "Service target address")
	verboseOpt := flag.Bool("verbose",false, "Verbose")
	flag.Parse()

	if *clientAddr != "" && *clientCert != "" {
		client := &SockClient{
			Address:  *clientAddr,
			Verbose: *verboseOpt,
		}
		client.Run()
	} else if *clientAddr != "" && *clientCert == "" {
		client := &SockClient{
			Address:  *clientAddr,
			CertPath: *clientCert,
			Verbose: *verboseOpt,

		}
		client.Run()
	} else if *serverPort != "" {
		server := &SockServer{
			Target:  *targetAddr,
			Verbose: *verboseOpt,
		}
		addr := fmt.Sprintf("0.0.0.0:%s", *serverPort)
		log.Printf("Listening at %s, forwarding to %s\n", addr, *targetAddr)
		log.Fatal(http.ListenAndServe(addr, server))
	} else {
		fmt.Println("Error, no arguments")
		os.Exit(1)
	}
}
