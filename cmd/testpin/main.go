package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
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
	gpioFile, err := os.OpenFile(fmt.Sprintf("%s/gpio%d/value", gpioPath, pin), os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer gpioFile.Close()

	value := 0
	for range time.Tick(3 * time.Second) {
		buf := make([]byte, 1024)
		n, err := gpioFile.ReadAt(buf[:], 0)
		if err != nil {
			log.Printf("err: %s, n: %d, v: %v", err, n, buf[:n])
			continue
		}
		if value != int(buf[0]) {
			fmt.Printf("Valor de la GPIO %d: %s\n", pin, string(buf[:]))
		}
	}

}
