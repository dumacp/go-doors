package listenner

import (
	"context"
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-doors/api"
	"github.com/dumacp/go-doors/internal/pubsub"
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
			Listen(contxt, v, func(evt *gpiosysfsne.Event) {
				ctxroot.RequestWithCustomSender(self, &Event{ID: id, Evt: evt}, self)
			})
		}
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
	case *actor.Stopped:
	case *actor.Terminated:

	case *Event:
		log.Printf("event: %v", msg)
		id, ok := a.gpioDoor[msg.ID]
		if !ok {
			logs.LogWarn.Printf("unknown door id -> pin (%d), id (%d)", msg.ID, id)
			break
		}
		a.lastEvent[msg.ID] = msg.Evt

		val := &api.ValuePubSubEvent{
			Id: int32(id),
			State: func() int32 {
				if msg.Evt.RisingEdge {
					return int32(1)
				}
				return int32(0)
			}(),
			Coord: "",
		}
		extmss := &api.PubSubEvent{
			Value:     val,
			Timestamp: float64(msg.Evt.Time.UnixNano() / 1000000000),
			Type:      "DOOR",
		}

		data, err := json.Marshal(extmss)
		if err != nil {
			logs.LogWarn.Printf("error sending event to pubsub: %s", err)
			break
		}
		pubsub.Publish(api.TOPIC_EVENTS, data)

	}
}
