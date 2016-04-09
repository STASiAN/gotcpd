package main

import (
	"os"
	"fmt"
	"log"
	"net"
	"time"
	"syscall"
	"os/signal"
//	"encoding/json"
	"github.com/takama/daemon"
//	"bytes"
	"bufio"
	"bytes"
	"encoding/json"
)

const (
	_DN    		= "goeppd"
	_DD 		= "goeppd"
	_LT		= "\x0D\x0A"
	_KVT 		= ":"               // header value separator
	_READ_BUF     	= 512               // buffer size for socket reader
)

var (
	EPPDPORT, EPPDHOST, CSLOGIN, CSPASS, TLDLOGIN, TLDPASS, TLDHOST string
	LOGPATH = fmt.Sprintf("/var/log/%s.log", _DN)
	stdlog, errlog *log.Logger
)

type Config struct {
	EPPD EPPD
	Zaloopa Zaloopa
	Master Master
}

type EPPD struct {
	EPPDPort string
	EPPDHost string
}

type Zaloopa struct{
	CsLogin string
	CsPass string
}

type Master struct {
	TLDLogin string
	TLDPass string
	TLDHost string
}

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {
	usage := "Usage: myservice install | remove | start | stop | status"
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	ln, err := net.Listen("tcp", EPPDHOST + ":" + EPPDPORT)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for connections on " + EPPDHOST + ":" + EPPDPORT)
//	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
//			continue
		}
//		fmt.Println(conn)

		spawn(conn) //BUFF parse xml xuemel!!!!!!!!!!!!!!!!!!!!!!!!!!!
		fmt.Println("0xDEADBEEFFFF")
		conn.Close()
	}
	return usage, nil
}

func spawn(c net.Conn) {
	if c != nil {
		r := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
		fmt.Println(r)
		conn, _ := net.Dial("tcp", TLDHOST+":"+EPPDPORT)
		buf0 := bytes.NewBufferString("")
		buf0.Write([]byte(CSLOGIN))
		buf0.Write([]byte(_LT))
		buf0.Write([]byte(CSPASS))
		buf0.Write([]byte(_LT))
		buf0.Write([]byte("pidor"))
		buf0.Write([]byte(_LT))
		fmt.Fprintf(conn, buf0.String())
		conn.Close()
	}
}

func dialback(f int) {
	fmt.Println(f)
}

func init() {
	file, e1 := os.Open("/tmp/config.json")
	if e1 != nil {
		fmt.Println("Error: ", e1)
	}
	decoder := json.NewDecoder(file)
	conf := Config{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	EPPDHOST = conf.EPPD.EPPDHost
	EPPDPORT = conf.EPPD.EPPDPort
	CSLOGIN = conf.Zaloopa.CsLogin
	CSPASS = conf.Zaloopa.CsPass
	TLDHOST = conf.Master.TLDHost
	TLDLOGIN = conf.Master.TLDLogin
	TLDPASS = conf.Master.TLDPass

	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	srv, err := daemon.New(_DN, _DD)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}

func LoggerMap(s map[string]string) {
  	tf := timeFormat()
	f, _ := os.OpenFile(LOGPATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
 	log.SetOutput(f)
  	log.Print(tf)
  	log.Print(s)
  	fmt.Println(s)
}

func LoggerString(s string) {
	tf := timeFormat()
	f, _ := os.OpenFile(LOGPATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
	log.Print(tf)
	log.Print(s)
	fmt.Println(s)
}

func LoggerMapMap(m map[string][]map[string]string) {
	f, _ := os.OpenFile(LOGPATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
	log.Print(m)
	fmt.Println(m)
}

func LoggerErr(e error) {
	f, _ := os.OpenFile(LOGPATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
	log.Print(e)
	fmt.Println(e)
}

func timeFormat() (string) {
	t := time.Now()
  	tf := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
  	return tf
}
