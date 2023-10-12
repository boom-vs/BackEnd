package controllers

import (
	"crm-backend/types"
	"crypto/sha1"
	"fmt"
	"io"
)

type ControllerTest struct {
}

func (ct *ControllerTest) Connection(context *types.RequestContext) {
	maxSize := 10 * 1024 * 1024
	currentSize := 0
	counter := 0

	reset := context.Response.Data

	for currentSize < maxSize {

		counter += 1
		hashHandler := sha1.New()
		_, err := io.WriteString(hashHandler, fmt.Sprint(counter))
		if err != nil {
			return
		}

		hash := fmt.Sprintf("%x", hashHandler.Sum(nil))

		context.Response.Data = reset
		context.Response.Data = append(context.Response.Data, struct {
			Counter int
			Hash    string
		}{
			Counter: counter,
			Hash:    hash,
		})

		currentSize += len(hash)

		context.Socket.WriteJSON(context.Response)

	}

}

func (ct *ControllerTest) Ping(context *types.RequestContext) {
	context.Response.Data = append(context.Response.Data, struct {
		Pong string
	}{
		Pong: "",
	})
}
