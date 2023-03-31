package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

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

	fdSet := &unix.FdSet{}
	fdSet.Bits[int(gpioFile.Fd()/64)] |= 1 << uint(gpioFile.Fd()%64)

	// bucle infinito para esperar cambios en el estado de la GPIO
	value := -1
	for range time.Tick(600 * time.Millisecond) {
		// configurar el conjunto de descriptores de archivo para esperar en
		timeout := unix.NsecToTimeval(10000000) // establecer un tiempo de espera de 10 ms
		ready, err := unix.Select(int(gpioFile.Fd()+1), fdSet, nil, nil, &timeout)
		if err != nil {
			panic(err)
		}

		// verificar si la llamada a select() ha expirado o si hay descriptores de archivo listos
		if ready > 0 {
			// leer el valor del archivo para determinar el nuevo estado de la GPIO
			// var buf [1]byte
			// _, err = gpioFile.Read(buf[:])
			// if err != nil {
			//     panic(err)
			// }

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

}
