package stream

type StreamHandler struct {
	Stream *Stream
}

func NewStreamHandler() *StreamHandler {
	return &StreamHandler{
		Stream: NewStream(),
	}
}
