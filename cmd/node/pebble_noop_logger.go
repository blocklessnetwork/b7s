package main

type pebbleNoopLogger struct{}

func (p *pebbleNoopLogger) Infof(_ string, _ ...any) {}

func (p *pebbleNoopLogger) Fatalf(_ string, _ ...any) {}
