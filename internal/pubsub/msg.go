package pubsub

type Msg struct {
	RoomId  string `json:"room_id"`
	Content []byte `json:"content"`
}
