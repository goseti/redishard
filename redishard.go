package main

import (
    "net"
    "log"
    "fmt"
    "runtime"
    "github.com/goseti/redishard/client"
    . "github.com/goseti/types"
)

type RedisServer struct {
    Host string
    Port int
}

func (rs RedisServer)Addr() string {
    return fmt.Sprintf("%s:%d", rs.Host, rs.Port)
}

func (rs RedisServer)Send(data []byte) []byte {
    log.Println("use redis, Addr:", rs.Addr())

    redis, err := net.Dial("tcp", rs.Addr())
    if nil != err {
        log.Println(err)
        return []byte{}
    }
    defer redis.Close()

    redis.Write(data)

    buffer := make([]byte, 2048)
    n, _ := redis.Read(buffer)
    log.Printf("redis send back: %d\n%s", n, string(buffer[:n]))
    return buffer[:n]
}

var (
    redis_pool Slice = Slice {
        RedisServer{"10.211.55.7", 6380},
        RedisServer{"10.211.55.7", 6381},
        RedisServer{"10.211.55.7", 6382},
        RedisServer{"10.211.55.7", 6383},
    }
)

func main() {
    netListen, _ := net.Listen("tcp", "localhost:9736")
    defer netListen.Close()

    log.Println("Waiting for connection")

    for {
        conn, err := netListen.Accept()
        if err != nil {
            continue
        }

        log.Println(conn.RemoteAddr().String(), "tcp connect success")
        go handleConnection(conn)
    }
}

//处理连接
func handleConnection(conn net.Conn) {
    defer log.Println("handle connection done!")
    defer conn.Close()

    rs := redis_pool.Random().(RedisServer)

    reader := make(chan []byte)

    c := client.NewClient(conn, reader)
    go c.Read()
    for {
        select {
        case command := <- reader:
            log.Println("read command:", (command))
            response := rs.Send(command)
            log.Println("redis response:", response)
            conn.Write(response)
        default:
            runtime.Gosched()
        }
    }
/*
    for {
        n, err := conn.Read(buffer)
        log.Printf("bytes read: %d\n%s\n", n, string(buffer[:n]))
        if err != nil {
            log.Println(conn.RemoteAddr().String(), "connection error:", err)
            break;
        } else {
            // response := rs.Send(buffer[:n])
            // conn.Write(response)
        }
    }
//*/
}
