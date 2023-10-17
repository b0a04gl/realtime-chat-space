package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestHandleConnections(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(handleConnections))
    defer server.Close()

    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fail()
    }
    defer ws.Close()

    mu.Lock()
    if len(clients) != 1 {
        t.Fail()
    }
    mu.Unlock()
}
