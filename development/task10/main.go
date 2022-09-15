package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

/*
Реализовать простейший telnet-клиент.

Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Требования:
Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP. После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться. При подключении к несуществующему сервер, программа должна завершаться через timeout
*/
func main() {
	timeoutFlag := flag.Int("timeout", 10, "timeout")
	flag.Parse()
	exitSignal := make(chan os.Signal)
	if os.Args[0] != "go-telnet" {
		os.Exit(1)
	}
	timeout := time.Duration(*timeoutFlag) * time.Second
	host := os.Args[len(os.Args)-2]                                   //хост = предпоследняя запись в пользовательском вводе
	port := os.Args[len(os.Args)-1]                                   //порт = последняя запись в пользовательском вводе
	signal.Notify(exitSignal, os.Interrupt, os.Kill)                  //ctrl+d это зарезервированный сигнал kill
	connectionParametrs := net.JoinHostPort(host, port)               //сокет с таймаутом
	conn, err := net.DialTimeout("tcp", connectionParametrs, timeout) //подключаемся
	if err != nil {
		log.Fatal(err) //обрабатываем ошибку
	}
	go func() { //прослушиваем на получение сигнала kill (ctrl+d)
		<-exitSignal //ждем сигнал в канале
		fmt.Println("Выходим из программы...")
		conn.Close() //закрываем сокет
		os.Exit(0)   //закрываем программу с кодом 0 который говорит что закрываемся в штатном режиме
	}()
	go func() { //прослушиваем на получение сообщений с сервера

		_, err := io.Copy(conn, os.Stdout) //копируем из консоли в сокет
		if err != nil {
			log.Fatal("Ошибка с сервера")

		} else {
			fmt.Println("Сообщение отправлено")
		}
	}()
	go func() { //прослушиваем на отсылку сообщений с stdin
		_, err := io.Copy(os.Stdin, conn) //копируем из сокета в консоль
		if err != nil {
			log.Fatal("Ошибка во время отправки сообщения")
		}
	}()
}
