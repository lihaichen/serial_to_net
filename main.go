package main
import (
	"fmt"
	"github.com/tarm/serial"
	"os"
	"time"
	"sync"
	"net"
)
var w sync.WaitGroup
var conn net.Conn
var serialP *serial.Port

func serialRead(s *serial.Port) {
	buf := make([]byte,0)
	recv_len := 0
	var err error
	for {
		n := 0
		tmp := make([]byte,2048)
		n, err = s.Read(tmp)
		if err != nil {
			if recv_len > 0 {
				_, err = conn.Write(buf[:recv_len])
				if err != nil {
					fmt.Println("write net error", err)
				}
				fmt.Println("serial recv len", recv_len)
				for i:=0; i < recv_len; i++ {
					fmt.Printf("%02X ", buf[i])
				}
				fmt.Println("")
				buf = make([]byte,0)
			}
			recv_len = 0
		} else {
			recv_len += n
			buf = append(buf, tmp[:n]...)
		}
	}
}
func socketRead(conn net.Conn) {
	for{
		tmp := make([]byte,2048)
		recvLen , _ := conn.Read(tmp)
		fmt.Println("net recv len", recvLen)
		for i:=0; i < recvLen; i++ {
			fmt.Printf("%02X ", tmp[i])
		}
		fmt.Println("")
		serialP.Write(tmp[:recvLen])
	}
}

func main() {
	fmt.Println("开始串口转网口")
	args := os.Args
	if args == nil || len(args) != 4 {
		fmt.Println("参数错误 ip 端口 com");
		os.Exit(-1)
	}
	ip := os.Args[1]
	port := os.Args[2]
	com := os.Args[3]
	fmt.Printf("eth[%s:%s] <==>%s\n", ip,port,com)
	w.Add(2)
	serialConfig := &serial.Config{Name:com, Parity: serial.ParityEven,
		Baud: 9600, ReadTimeout:time.Microsecond * 200}
	var err error
	serialP, err = serial.OpenPort(serialConfig)
	if err != nil {
		fmt.Println("打开串口失败",err)
		os.Exit(-1)
	}
	conn, err = net.Dial("tcp",ip + ":" + port)
	if err != nil {
		fmt.Println("连接网口失败",err)
		os.Exit(-1)
	}
	go serialRead(serialP)
	go socketRead(conn)
	w.Wait()
}