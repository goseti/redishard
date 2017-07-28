package client

import (
    "net"
    "log"
)

type client struct {
    command chan<- []byte
    conn net.Conn
    buffer []byte
}

func NewClient(conn net.Conn, command chan []byte) *client {
    return &client{
        command: command,
        conn: conn,
        buffer: make([]byte, 0),
    }
}

func (c *client)Read() {
    buffer := make([]byte, 2048)
    for {
        n, err := c.conn.Read(buffer)
        log.Printf("bytes read: %d\n%s\n", n, string(buffer[:n]))
        if err != nil {
            log.Println(c.conn.RemoteAddr().String(), "connection error:", err)
            break;
        }
        c.parseData(buffer[:n])
    }
}

func (c *client)parseData(data []byte) {
    start := 0
    for i, b := range data {
        if b == '\n' && i > 3 && data[i - 1] == '\r' && data[i - 2] == '\n' && data[i - 3] == '\r' {
            command := append(c.buffer, data[start:i + 1]...)
            log.Println("got new command:\n", string(command))
            c.command <- command
            start = i
            c.buffer = make([]byte, 0)
        }
    }
    c.buffer = append(c.buffer, data[start:]...)
    log.Println("buffer:", string(c.buffer))
    log.Println("buffer bytes:", c.buffer)
}