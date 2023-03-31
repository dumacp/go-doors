package listenner

import (
	"context"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-doors/pkg/gpiosysfsne"
	"github.com/dumacp/go-logs/pkg/logs"
)

type Actor struct {
	gpioDoor  map[int]int
	lastEvent map[int]*gpiosysfsne.Event
	cancel    func()
}

func NewActor() *Actor {
	a := &Actor{
		gpioDoor:  make(map[int]int),
		lastEvent: make(map[int]*gpiosysfsne.Event),
	}
	return a
}

func (a *Actor) AddGpioWatch(id, gpionumber int) {
	a.gpioDoor[id] = gpionumber
}

func (a *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		ctxroot := ctx.ActorSystem().Root
		self := ctx.Self()
		contxt, cancel := context.WithCancel(context.TODO())
		a.cancel = cancel
		for id, v := range a.gpioDoor {
			id := id
			fmt.Printf("id: %d\n", id)
			Listen(contxt, v, func(evt *gpiosysfsne.Event) {
				ctxroot.RequestWithCustomSender(self, &Event{ID: id, Evt: evt}, self)
			})
		}
		logs.LogInfo.Printf("actor %q started", ctx.Self().GetId())
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
		logs.LogInfo.Printf("actor %q stopping", ctx.Self().GetId())
	case *actor.Stopped:
	case *actor.Terminated:

	case *Event:
		if ctx.Parent() != nil {
			ctx.Request(ctx.Parent(), msg)
		}
	}
}
