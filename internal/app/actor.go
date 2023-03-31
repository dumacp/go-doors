package app

import (
	"encoding/json"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-doors/api"
	"github.com/dumacp/go-doors/internal/gps"
	"github.com/dumacp/go-doors/internal/listenner"
	"github.com/dumacp/go-doors/internal/pubsub"
	"github.com/dumacp/go-logs/pkg/logs"
	"google.golang.org/protobuf/encoding/protojson"
)

type Actor struct {
	propsGsp       *actor.Props
	pidGps         *actor.PID
	propsListenner *actor.Props
	lastEvents     map[int]*listenner.Event
}

func NewActor(gps, listen *actor.Props) actor.Actor {
	a := &Actor{
		propsGsp:       gps,
		propsListenner: listen,
		lastEvents:     make(map[int]*listenner.Event),
	}
	return a
}

func (a *Actor) Receive(ctx actor.Context) {
	logs.LogBuild.Printf("Message arrived in %s: %s, %T, %s", ctx.Self().GetId(),
		ctx.Message(), ctx.Message(), ctx.Sender())
	switch msg := ctx.Message().(type) {
	case *actor.Started:

		if a.propsGsp != nil {
			ctx.SpawnNamed(a.propsGsp, "gps")
		}
		if a.propsListenner != nil {
			ctx.SpawnNamed(a.propsListenner, "listenner")
		}
		pubsub.Subscribe("DOORS", ctx.Self(), ParseRequest)
		logs.LogInfo.Printf("actor %q started", ctx.Self().GetId())
	case *api.DoorsStateRequest:
		for id, v := range a.lastEvents {
			val := &api.DoorState{
				Id: int32(id),
				Value: func() int32 {
					if v.Evt.RisingEdge {
						return int32(1)
					}
					return int32(0)
				}(),
			}
			data, err := json.Marshal(val)
			if err != nil {
				logs.LogWarn.Printf("error sending event to pubsub: %s", err)
				break
			}
			pubsub.Publish(msg.Topic, data)
		}
	case *listenner.Event:
		id := msg.ID
		coord := ""
		a.lastEvents[id] = msg
		if a.pidGps != nil {
			req := ctx.RequestFuture(a.pidGps, &gps.MsgGPS{}, 1*time.Second)
			if err := req.Wait(); err != nil {
				logs.LogWarn.Print("empty gps frame in ignition event")
			}
			frame, err := req.Result()
			if err != nil {
				logs.LogWarn.Printf("empty gps frame in ignition event, error: %s", err)
			} else {
				if frame != nil {
					if data, ok := frame.([]byte); ok {
						coord = string(data)
					}
				} else {
					logs.LogWarn.Print("empty gps frame in ignition event")
				}
			}
		}
		val := &api.ValuePubSubEvent{
			Id: int32(id),
			State: func() int32 {
				if msg.Evt.RisingEdge {
					return int32(1)
				}
				return int32(0)
			}(),
			Coord: coord,
		}
		extmss := &api.PubSubEvent{
			Value:     val,
			Timestamp: float64(msg.Evt.Time.UnixNano() / 1000000000),
			Type:      "DOOR",
		}

		m := protojson.MarshalOptions{EmitUnpopulated: true}
		data, err := m.Marshal(extmss)
		if err != nil {
			logs.LogWarn.Printf("error sending event to pubsub: %s", err)
			break
		}
		pubsub.Publish(api.TOPIC_EVENTS, data)

	case *actor.Stopped:
	case *actor.Terminated:
	}
}
