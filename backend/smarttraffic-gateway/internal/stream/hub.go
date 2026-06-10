package stream

import (
	"net/http"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	buffer  []byte
	clients map[chan []byte]struct{}
}

func NewHub(bufferBytes int) *Hub {
	return &Hub{
		buffer:  make([]byte, 0, bufferBytes),
		clients: make(map[chan []byte]struct{}),
	}
}

func (h *Hub) Write(chunk []byte) (int, error) {
	data := append([]byte(nil), chunk...)

	h.mu.Lock()
	h.buffer = append(h.buffer, data...)
	if extra := len(h.buffer) - cap(h.buffer); extra > 0 {
		h.buffer = append([]byte(nil), h.buffer[extra:]...)
	}

	for client := range h.clients {
		select {
		case client <- data:
		default:
			close(client)
			delete(h.clients, client)
		}
	}
	h.mu.Unlock()

	return len(chunk), nil
}

func (h *Hub) HandleFLV(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "video/x-flv")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")

	flusher, _ := w.(http.Flusher)
	client := make(chan []byte, 32)

	h.mu.Lock()
	bufferSnapshot := append([]byte(nil), h.buffer...)
	h.clients[client] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		if _, ok := h.clients[client]; ok {
			delete(h.clients, client)
			close(client)
		}
		h.mu.Unlock()
	}()

	if len(bufferSnapshot) > 0 {
		if _, err := w.Write(bufferSnapshot); err != nil {
			return
		}
		if flusher != nil {
			flusher.Flush()
		}
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case data, ok := <-client:
			if !ok {
				return
			}
			if _, err := w.Write(data); err != nil {
				return
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
	}
}
