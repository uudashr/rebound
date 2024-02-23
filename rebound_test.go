package rebound_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/uudashr/rebound"
)

func ExampleRebound() {
	rb := &rebound.Rebound{}

	type OrderCompleted struct {
		OrderID string
	}

	rb.ReactTo("order.completed", func(event OrderCompleted) error {
		fmt.Printf("Order %q is completed\n", event.OrderID)
		return nil
	})

	// the "order.completed" can be derived from nats streaming subject, kafka topic, etc.
	rb.Dispatch("order.completed", []byte(`{"OrderID":"123"}`))

	// Output:
	// Order "123" is completed
}

func TestNatsSubjectParsing(t *testing.T) {
	subject := "sales.events.private.order.123.completed"
	prefix := "sales.events.private."
	if got, want := strings.HasPrefix(subject, prefix), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}

	eventQualifiedName := subject[len(prefix):]
	if got, want := eventQualifiedName, "order.123.completed"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	parts := strings.Split(eventQualifiedName, ".")
	if got, want := len(parts), 3; got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	eventName := fmt.Sprintf("%s.%s", parts[0], parts[2])
	if got, want := eventName, "order.completed"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
