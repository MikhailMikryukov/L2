package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var timeout int
	flag.IntVar(&timeout, "timeout", 0, "timeout in seconds")
	flag.Parse()

	host := os.Args[0]
	port := os.Args[1]

	address := host + ":" + port

	var conn net.Conn
	var err error
	if timeout > 0 {
		conn, err = net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	} else {
		conn, err = net.Dial("tcp", address)
	}

	if err != nil {
		fmt.Println("connection err:", err)
		return
	}

	fmt.Println("connected to", address)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Читаем из STDIN и отправляем в сокет
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				text := scanner.Text()
				_, err := conn.Write([]byte(text + "\n"))
				if err != nil {
					fmt.Println("write message error", err)
				}
				fmt.Println("sent to socket:", text)
			}
		}

		// При Ctrl+D заканчивается сканирование и закрываем соединение
		if err := scanner.Err(); err != nil {
			fmt.Println("reading STDIN err:", err)
		} else {
			fmt.Println("input finished")
		}

		fmt.Println("closing the connection")
		conn.Close()
	}()

	// Читаем из сокета и отправляем в STDOUT
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		connReader := bufio.NewReader(conn)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := connReader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						fmt.Println("connection closed ")
						return
					}
					if strings.Contains(err.Error(), "closed network") {
						return
					}

					fmt.Println("read message error", err)
					return
				}
				fmt.Println("read from socket: ", msg)
			}
		}
	}()

	wg.Wait()
}
