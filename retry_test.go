package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExponential(t *testing.T) {
	exp := Exponential(2)
	if d := exp(0); d != 0 {
		t.Fatal(d)
	}

	if d := exp(1); d != 2 {
		t.Fatal(d)
	}

	if d := exp(10); d != 20 {
		t.Fatal(d)
	}

	exp = Exponential(1.8)
	if d := exp(10); d != 18 {
		t.Fatal(d)
	}
}

func TestEnsure(t *testing.T) {
	r := New(WithBaseDelay(1 * time.Millisecond))

	val := 0
	do := func() error {
		val++
		t.Log(val)
		if val == 5 {
			return nil
		}
		return Retriable(errors.New("please retry"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := r.Ensure(ctx, do); err != nil {
		t.Fatal(err)
	}

	if val != 5 {
		t.Fatal(val)
	}
}

func TestWithBaseDelay(t *testing.T) {
	r := &Retry{}
	opt := WithBaseDelay(1)
	opt(r)

	if r.base != 1 {
		t.Fatal(r.base)
	}
}
func TestWithBackoff(t *testing.T) {
	r := &Retry{}
	opt := WithBackoff(nil)
	opt(r)
	if r.backoff != nil {
		t.Fatal(r.backoff)
	}

	opt = WithBackoff(Exponential(2))
	opt(r)
	if r.backoff == nil {
		t.Fatal(r.backoff)
	}
}

func TestEnsureNWithContext(t *testing.T) {
	r = New(WithBaseDelay(1 * time.Millisecond))

	val := 0
	do := func() error {
		val++
		t.Log(val)
		if val == 5 {
			return nil
		}
		return Retriable(errors.New("please retry"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err := r.EnsureNWithContext(ctx, 3, do)
	cancel()

	if err == nil {
		t.Fatal("should be error")
	}
	if val != 3 {
		t.Fatal(val)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Millisecond)
	err = r.EnsureNWithContext(ctx, 5, do)
	cancel()

	if err == nil || err != context.DeadlineExceeded {
		t.Fatal("should be error")
	}
	if val == 5 {
		t.Fatal(val)
	}
}
