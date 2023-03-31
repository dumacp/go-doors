package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	gpioPath = "/sys/class/gpio" // ruta del archivo que representa la GPIO
)

var pin int

func init() {
	flag.IntVar(&pin, "pin", 163, "ping number")
}

func main() {
	flag.Parse()
	gpioFile, err := os.OpenFile(fmt.Sprintf("%s/gpio%d/value", gpioPath, pin), os.O_RDONLY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		panic(err)
	}
	defer gpioFile.Close()

	// crear una nueva instancia de epoll
	epollFd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	defer unix.Close(epollFd)

	// crear una estructura epoll_event para representar el archivo
	event := unix.EpollEvent{
		Fd:     int32(gpioFile.Fd()),
		Events: unix.EPOLLPRI,
	}

	// agregar el archivo a la lista de descriptores monitoreados por epoll
	if err := unix.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, int(gpioFile.Fd()), &event); err != nil {
		panic(err)
	}

	// bucle infinito para esperar cambios en el estado de la GPIO
	value := -1
	for {
		// esperar cambios en el archivo utilizando epoll_wait()
		events := make([]unix.EpollEvent, 1)
		_, err := unix.EpollWait(epollFd, events, -1)
		if err != nil {
			panic(err)
		}

		// leer el valor del archivo para determinar el nuevo estado de la GPIO
		buf := make([]byte, 1024)
		n, err := gpioFile.ReadAt(buf[:], 0)
		if err != nil {
			if !errors.Is(err, io.EOF) || n <= 0 {
				panic(err)
			}
		}

		// imprimir el valor actual de la GPIO
		if value != int(buf[0]) {
			value = int(buf[0])
			fmt.Printf("Valor de la GPIO 17: %s\n", string(buf[:]))
		}
	}
}
