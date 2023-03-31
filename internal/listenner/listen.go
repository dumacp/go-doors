package listenner

import (
	"context"
	"fmt"
	"time"

	"github.com/dumacp/go-doors/pkg/gpiosysfsne"
)

func Listen(ctx context.Context, gpionumber int, fn func(evt *gpiosysfsne.Event)) error {

	pin, err := gpiosysfsne.OpenPinWithEvents(gpionumber)
	if err != nil {
		return err
	}

	ch := pin.Events()

	go func() {
		defer pin.Close()
		timeout := 2 * time.Second
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		backupevent := new(gpiosysfsne.Event)
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-ch:
				fmt.Printf("event: %v\n", v)
				if time.Since(backupevent.Time) > timeout {
					fn(v)
				} else {
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}
					timer.Reset(timeout)
				}
				backupevent = new(gpiosysfsne.Event)
				backupevent.RisingEdge = v.RisingEdge
				backupevent.Time = v.Time

			case <-timer.C:
				fn(backupevent)
			}
		}
	}()
	return nil
}
