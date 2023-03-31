package listenner

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/dumacp/go-doors/api"
	"github.com/dumacp/go-doors/pkg/gpiosysfsne"
)

type Actor struct {
	gpioDoor []int
	lastEvent  *gpiosysfsne.Event
	evs        *eventstream.EventStream
	cancel     func()
}

func (a *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:

		ctxroot := ctx.ActorSystem().Root
		self := ctx.Self()

		contxt, cancel := context.WithCancel(context.TODO())
		a.cancel = cancel
		Listen(contxt, a.gpioDoor_0, func(evt *gpiosysfsne.Event) {
			ctxroot.RequestWithCustomSender(self, &event{pin: a.gpioDoor_0, evt: evt}, self)
		})
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
	case *actor.Stopped:
	case *actor.Terminated:

	case *event:
		a.lastEvent = msg
		extmss := &api.DoorState{
			Id: ,
		}

	}
}
