package grpc

import (
	"testing"
	"time"

	"github.com/AndreyVLZ/curly-octo/server/pkg/jwt"
)

type mockFileStore struct {
}

func TestStartStop(t *testing.T) {
	jwt := jwt.New("KEY", time.Minute)
	gServer := NewGRPCServere(":3300", jwt, nil, nil, 10000)

	/*
		t.Cleanup(func() {
			gServer.Stop()
		})
	*/
	go func() {
		time.Sleep(3 * time.Second)
		gServer.Stop()
	}()

	if err := gServer.Start(); err != nil {
		t.Logf("gSrv start: %v\n", err)
	}
}
