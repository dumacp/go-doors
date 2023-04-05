package doorsmsg

type DoorState struct {
	ID    int `json:"id"`
	Value int `json:"value"`
}
type DoorsStateRequest struct {
	TopicResponse string `json:"topic"`
}
type DoorsState struct {
	States map[int]DoorState `json:"states"`
}
