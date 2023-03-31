package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-doors/internal/app"
	"github.com/dumacp/go-doors/internal/gps"
	"github.com/dumacp/go-doors/internal/listenner"
	"github.com/dumacp/go-doors/internal/pubsub"
	"github.com/dumacp/go-logs/pkg/logs"
)

const (
	GPIOPuerta1 int = 162
	GPIOPuerta2 int = 163
)

const (
	version = "1.2.0_test"
)

var mqtt, debug, logStd bool
var showVersion bool
var gpioFrontDoor, gpioRearDoor int

func init() {
	flag.BoolVar(&mqtt, "mqtt", false, "send messages to local broker. (DECRECATED)")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.BoolVar(&logStd, "logStd", false, "log All in stderr")
	flag.IntVar(&gpioFrontDoor, "frontaldoor", GPIOPuerta1, "gpionumber for frontal door")
	flag.IntVar(&gpioRearDoor, "backdoor", GPIOPuerta2, "gpionumber for back door")
	flag.BoolVar(&showVersion, "version", false, "show version")
}

//var pub *pubsub.PubSub

func main() {

	flag.Parse()
	if showVersion {
		fmt.Printf("version: %v\n", version)
		os.Exit(2)
	}

	ctxroot := actor.NewActorSystem().Root

	initLogs(debug, logStd)
	logs.LogInfo.Printf("running, version: %s", version)
	pubsub.Init(ctxroot)

	time.Sleep(3 * time.Second)

	listnn := listenner.NewActor()
	listnn.AddGpioWatch(0, gpioFrontDoor)
	listnn.AddGpioWatch(1, gpioRearDoor)

	propsGsp := actor.PropsFromFunc(gps.NewActor().Receive)

	propsLisnn := actor.PropsFromFunc(listnn.Receive)

	if _, err := ctxroot.SpawnNamed(
		actor.PropsFromFunc(app.NewActor(propsGsp, propsLisnn).Receive),
		"doors-actor"); err != nil {
		time.Sleep(1 * time.Second)
		logs.LogError.Fatalln(err)
	}

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)

	for range finish {
		logs.LogWarn.Print("Finish")
		time.Sleep(300 * time.Millisecond)
		return
	}

}
