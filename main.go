package main

import (
    "os"
    "fmt"
    "log"
    "net"
    "flag"
    "time"
    "net/http"
    "os/signal"
    
    "github.com/gorilla/websocket"
)

func gracefulExit(l net.Listener, c net.Conn, wc *websocket.Conn, msg string){
    wc.Close()
    c.Close()
    l.Close()
    log.Println(msg)
    os.Exit(0)
}

func WSocket2UDSocket(l net.Listener, c net.Conn, wc *websocket.Conn) {
    for {
        _, message, err := wc.ReadMessage()
        if err != nil {
            gracefulExit(l, c, wc, fmt.Sprintf("Websockets Read: %s", err))
        }
        
        _, err = c.Write(message)
        if err != nil {            
            gracefulExit(l, c, wc, fmt.Sprintf("Unable to write unix socket: %s", err))
        }
    }
}

func UDSocket2WSocket(l net.Listener, c net.Conn, wc *websocket.Conn) {
    message := make([]byte, 4)
    for {
        nr, err := c.Read(message[:])
        if err != nil {
            gracefulExit(l, c, wc, fmt.Sprintf("Unix socket reading error: %s", err))
        }
        
        data := message[0:nr]
        if err := wc.WriteMessage(websocket.BinaryMessage, data); err != nil {
            gracefulExit(l, c, wc, fmt.Sprintf("Websockets writing error: %s", err))
        }
        time.Sleep(3*time.Millisecond)
    }
}

func main() {
    ws_url := flag.String("u", "", "Assign the websockets endpoint. Starting with ws:// or wss://")
    randN := flag.Int("r", 9487, "Assign a random number to distinguish unix domain sockets")
    flag.Parse()
    
    l, err := net.Listen("unix", fmt.Sprintf("/tmp/telnetBYwebsocket.%d.sock", *randN))
    if err != nil {
        log.Fatal("listen error:", err)
    }
    
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, os.Interrupt)
    go func(){
        for sig := range ch {
            l.Close()
            log.Println("Aborted by user", sig)
            os.Exit(0)
        }
    }()
    
    var dialer *websocket.Dialer
    //c, _, err := dialer.Dial("wss://ws.ptt.cc/bbs", http.Header{"Origin": {"app://pcman"}})
    c, _, err := dialer.Dial(*ws_url, http.Header{"Origin": {"app://pcman"}})
    if err != nil {
        l.Close()
    	log.Fatal("websocket connect error:", err)
    }

    for {
        fd, err := l.Accept()
        if err != nil {
            log.Fatal("accept error:", err)
        }
        
        go WSocket2UDSocket(l, fd, c)
        go UDSocket2WSocket(l, fd, c)
    }
    
}
