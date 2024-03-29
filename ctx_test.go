package zerolog

import (
	"context"
	"io/ioutil"
	"testing"
)

func TestCtx(t *testing.T) {
	log := New(ioutil.Discard)
	ctx := log.WithContext(context.Background())
	log2 := Ctx(ctx)
	if !log.IsEqualToLogger(*log2) {
		t.Error("Ctx did not return the expected logger")
	}

	// update
	log = log.Level(InfoLevel)
	ctx = log.WithContext(ctx)
	log2 = Ctx(ctx)
	if !log.IsEqualToLogger(*log2) {
		t.Error("Ctx did not return the expected logger")
	}

	log2 = Ctx(context.Background())
	if log2 != disabledLogger {
		t.Error("Ctx did not return the expected logger")
	}

	DefaultContextLogger = &log
	t.Cleanup(func() { DefaultContextLogger = nil })
	log2 = Ctx(context.Background())
	if log2 != &log {
		t.Error("Ctx did not return the expected logger")
	}
}

func TestCtxDisabled(t *testing.T) {
	dl := New(ioutil.Discard).Level(Disabled)
	ctx := dl.WithContext(context.Background())
	if ctx != context.Background() {
		t.Error("WithContext stored a disabled logger")
	}

	l := New(ioutil.Discard).With().Str("foo", "bar").Logger()
	ctx = l.WithContext(ctx)
	if !l.IsEqualToLogger(*Ctx(ctx)) {
		t.Error("WithContext did not store logger")
	}

	l.UpdateContext(func(c Context) Context {
		return c.Str("bar", "baz")
	})
	ctx = l.WithContext(ctx)
	if !l.IsEqualToLogger(*Ctx(ctx)) {
		t.Error("WithContext did not store updated logger")
	}

	l = l.Level(DebugLevel)
	ctx = l.WithContext(ctx)
	if !l.IsEqualToLogger(*Ctx(ctx)) {
		t.Error("WithContext did not store copied logger")
	}

	ctx = dl.WithContext(ctx)
	if !dl.IsEqualToLogger(*Ctx(ctx)) {
		t.Error("WithContext did not override logger with a disabled logger")
	}
}
