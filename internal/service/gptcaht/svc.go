package gptcaht

import "github.com/gone-io/gone"

type svc struct {
	gone.Goner
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}
