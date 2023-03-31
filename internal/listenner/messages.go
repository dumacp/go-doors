package listenner

import "github.com/dumacp/go-doors/pkg/gpiosysfsne"

type event struct {
	pin int
	evt *gpiosysfsne.Event
}
