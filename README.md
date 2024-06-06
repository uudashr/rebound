[![Go Reference](https://pkg.go.dev/badge/github.com/uudashr/rebound.svg)](https://pkg.go.dev/github.com/uudashr/rebound)

# Rebound
Rebound is simple event handling based on registered name. The handling is automatically unmarshaled to the registered type.

Think like: you want to dispatch a serialized event, then you want to handle it in a type-safe way.

It can be use to
1. Send and consume through the messaging/streaming service
2. Store/retrieve the event to and from the database (outbox pattern)

## Installation

```shell
go get github.com/uudashr/rebound
```

## Usage
```go
package main

import (
	"fmt"

	"github.com/uudashr/rebound"
)

type OrderCompleted struct {
    OrderID string
}

func main() {
    rb := &rebound.Rebound{}

	// Register the event
	rb.ReactTo("order.completed", func(event OrderCompleted) error {
        // Handle the "order.completed" event
		fmt.Printf("Order %q is completed\n", event.OrderID)
		return nil
	})

    // Dispatch the "order.completed" event
	rb.Dispatch("order.completed", []byte(`{"OrderID":"123"}`))
}
```
