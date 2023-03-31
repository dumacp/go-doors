package main

import (
	"flag"
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

var debug bool
var gpioFrontDoor, gpioRearDoor int

func init() {
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.IntVar(&gpioFrontDoor, "frontaldoor", 162, "gpionumber for frontal door")
	flag.IntVar(&gpioFrontDoor, "backdoor", 163, "gpionumber for back door")
}

func main() {
	flag.Parse()

	ctxroot := actor.NewActorSystem().Root

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
