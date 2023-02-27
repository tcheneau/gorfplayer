package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type serialHandler struct {
	port *serial.Port
}

type rfplayerOrder struct {
	Order     string
	Address   string
	Protocol  string
	Percent   string `json:",omitempty"`
	Burst     string `json:",omitempty"`
	Qualifier string `json:",omitempty"`
}

const (
	START      = "ZIA++"
	TERMINATOR = "\r"
)

func (h *serialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if r.URL.Path == "/v1/read" {
		buf := h.read()
		log.Println(buf)
		w.Write(buf)
		return
	}

	if r.URL.Path == "/v1/command" {

		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var order rfplayerOrder

		err := decoder.Decode(&order)
		if err != nil {
			log.Fatalln(err)
		}

		if order.Order == "" || order.Address == "" || order.Protocol == "" {
			http.Error(w, "Missing JSON field order, address or protocol", http.StatusNotAcceptable)
			return
		}

		command := strings.Join([]string{order.Order, order.Address, order.Protocol}, " ")

		if order.Percent != "" {
			command += " %" + order.Percent
		}

		if order.Burst != "" {
			command += " BURST " + order.Burst
		}

		if order.Qualifier != "" {
			command += " QUALIFIER " + order.Qualifier
		}

		h.sendCommand(command)

		buf := h.read()
		log.Println(buf)
		w.Write(buf)
		return
	}

	if r.URL.Path == "/v1/ping" {
		command := "PING"

		h.sendCommand(command)

		buf := h.read()
		log.Println(buf)
		w.Write(buf)
		return
	}

	if r.URL.Path == "/v1/status" {
		command := "STATUS SYSTEM JSON"

		h.sendCommand(command)

		buf := h.read()
		log.Println(buf)
		w.Write(buf)
		return
	}

	// catchall
	http.NotFound(w, r)
}

// send a command to the RFPlayer
func (h *serialHandler) sendCommand(command string) bool {
	log.Println(command)
	h.port.Write([]byte(START + command + TERMINATOR))

	return true

}

// read reply from RFPlayer
func (h *serialHandler) read() []byte {
	var n int
	var err error
	buf := make([]byte, 256)
	tbuf := make([]byte, 0)

	for {
		n, err = h.port.Read(buf)

		tbuf = append(tbuf, buf[:n]...)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if n == 0 {
			break
		}
	}

	if bytes.HasPrefix(tbuf, []byte{'Z', 'I', 'A', '-', '-'}) {
		return tbuf[5:]
	} else {
		//return []byte{}
		return tbuf
	}
}

// drain the serial port
func (h *serialHandler) drain() {
	buf := make([]byte, 128)

	var err error
	var n int = 1

	for n != 0 {
		n, err = h.port.Read(buf)

		log.Println(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}

}

func main() {
	// Command line flag to enable or disable TLS
	var tlsEnabled bool
	var portPath string
	flag.BoolVar(&tlsEnabled, "tls", false, "enable TLS connections")
	flag.StringVar(&portPath, "port", "/dev/ttyUSB0", "path to the serial port")
	flag.Parse()

	// Setup serial port
	// baudrate 115200, 8bit data, no parity, 	1 stop bit
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, ReadTimeout: time.Second * 1}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	sconn := &serialHandler{port: s}

	// this make the rfplayer talk in JSON
	sconn.sendCommand("FORMAT JSON")
	// there is no output to the previous command
	// we drain the serial port in case there is an output from a previous command
	sconn.drain()

	r := http.NewServeMux()
	// Add API endpoints
	r.Handle("/v1/command", sconn)
	r.Handle("/v1/status", sconn)
	r.Handle("/v1/read", sconn)
	r.Handle("/v1/ping", sconn)

	// Setup server
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Enable TLS if specified
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			PreferServerCipherSuites: true,
		},
	}

	// Start server
	if tlsEnabled {
		log.Println("Starting server with TLS enabled")
		log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
	} else {
		log.Printf("Starting server without TLS on port %d\n", 8000)
		log.Fatal(srv.ListenAndServe())
	}
}

func endpoint1Handler(w http.ResponseWriter, r *http.Request) {
	// Stub
}

func endpoint2Handler(w http.ResponseWriter, r *http.Request) {
	// Stub

	// Common code for all requests can go here...

	switch r.Method {
	case http.MethodGet:
		// Handle the GET request...

	case http.MethodPost:
		// Handle the POST request...

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
