package main_test

import (
	"net/http"
	"testing"

	"github.com/posener/wstest"
)

func TestHandler(t *testing.T) {
	h := &Handler{}
	d := wstest.NewDialer(h, nil)

	c, resp, err := d.Dial("ws://localhost:4000", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("resp.StatusCode: %q, want: %q", got, want)
	}

	err = c.WriteJSON("test")
	if err != nil {
		t.Fatal(err)
	}
}
