package listenner

import "github.com/dumacp/go-doors/pkg/gpiosysfsne"

type Event struct {
	ID  int
	Evt *gpiosysfsne.Event
}
