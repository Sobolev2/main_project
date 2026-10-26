package events

type MessageSent struct {
	Event      string `json:"event"`
	MessageID  int    `json:"message_id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
}
