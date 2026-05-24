package websocket

import (
	"context"
	"net/http"
	"testing"
)

type onlineTrackerStub struct {
	setCalls   []string
	clearCalls []string
}

func (s *onlineTrackerStub) MarkOnline(_ context.Context, uid string) error {
	s.setCalls = append(s.setCalls, uid)
	return nil
}

func (s *onlineTrackerStub) MarkOffline(_ context.Context, uid string) error {
	s.clearCalls = append(s.clearCalls, uid)
	return nil
}

func TestServerTracksOnlineStateOnAddAndClose(t *testing.T) {
	tracker := &onlineTrackerStub{}
	srv := NewServer("127.0.0.1:0", WithOnlineTracker(tracker))
	conn := &Conn{s: srv, done: make(chan struct{})}

	req := mustRequest(t, "u100")
	srv.addConn(conn, req)

	if len(tracker.setCalls) != 1 || tracker.setCalls[0] != "u100" {
		t.Fatalf("expected MarkOnline to be called with u100, got %#v", tracker.setCalls)
	}

	srv.Close(conn)

	if len(tracker.clearCalls) != 1 || tracker.clearCalls[0] != "u100" {
		t.Fatalf("expected MarkOffline to be called with u100, got %#v", tracker.clearCalls)
	}
}

func mustRequest(t *testing.T, userID string) *http.Request {
	t.Helper()
	req, err := http.NewRequest("GET", "http://example.com/ws?userId="+userID, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	return req
}
