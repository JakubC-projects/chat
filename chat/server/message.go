package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JakubC-projects/chat/chat"
	"github.com/google/uuid"
)

func (s *Server) handleSendMessage(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	req.ParseMultipartForm(1024)

	sendTime := time.Now()

	msg := chat.Message{
		ChatUid:   uuid.MustParse(req.FormValue("chat_uid")),
		Content:   req.FormValue("content"),
		Type:      chat.MessageType(req.FormValue("type")),
		SentAt:    sendTime,
		ChangedAt: sendTime,
	}

	_, err := s.store.SendMessage(ctx, msg)
	if err != nil {
		http.Error(w, "cannot send message: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMessageEvents(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	req.ParseForm()
	chatUid := uuid.MustParse(req.FormValue("chat_uid"))

	var afterDate time.Time
	if reqAfterDateStr := req.FormValue("after_date"); reqAfterDateStr != "" {
		t, err := time.Parse(time.RFC3339, reqAfterDateStr)
		if err == nil {
			afterDate = t
		}
	}

	eventSrc, err := s.publisher.Subscribe(ctx, "message")
	if err != nil {
		http.Error(w, "cannot listen to events", http.StatusInternalServerError)
		return
	}
	defer eventSrc.Close(ctx)

	messages, err := s.store.GetMessagesAfterDate(ctx, chatUid, afterDate)
	for _, msg := range messages {
		msgBytes, _ := json.Marshal(msg)
		fmt.Fprintf(w, "data: %s\n\n", string(msgBytes))
	}
	w.(http.Flusher).Flush()

	for {
		select {
		case <-req.Context().Done():
			return
		default:
			event, err := eventSrc.NextEvent(ctx)
			if err != nil {
				fmt.Println("error listening", err)
			}
			fmt.Fprintf(w, "data: %s\n\n", event.Payload)
			w.(http.Flusher).Flush()

		}
	}
}
