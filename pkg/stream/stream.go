package stream

import (
	"context"
	"net/http"
	"sync"
)

type Stream struct {
	Content   chan *StreamComponent
	Queue     []func() *StreamComponent
	WaitGroup sync.WaitGroup
}

func NewStream() *Stream {
	return &Stream{
		Content: make(chan *StreamComponent),
		Queue:   make([]func() *StreamComponent, 0),
	}
}

func (s *Stream) Stream(handler func() *StreamComponent) {
	s.Queue = append(s.Queue, func() *StreamComponent {
		content := handler()
		s.Content <- content
		s.WaitGroup.Done()
		return content
	})
	s.WaitGroup.Add(1)
}

func (s *Stream) Wait() {
	for _, handler := range s.Queue {
		go handler()
	}
	s.WaitGroup.Wait()
	close(s.Content)
}

func StreamHandlerFunc(handler func(w http.ResponseWriter, r *http.Request, s *Stream)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		flusher, _ := w.(http.Flusher)
		stream := NewStream()

		go stream.Wait()

		handler(w, r, stream)

		flusher.Flush()

		for {
			select {
			case content := <-stream.Content:
				if content == nil {
					return
				}
				if err := content.Render(context.Background(), w); err != nil {
					return
				}
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})
}
