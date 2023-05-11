package main

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/tracer"

	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal"
)

func main() {
	tracer.SetTraceId("", func() {
		gone.Serve(internal.MasterPriest)
	})
}
