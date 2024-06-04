package secrets_test

import (
	"context"

	"github.com/x-ethr/go-kubernetes-secrets"
)

func Example() {
	ctx := context.Background()

	instance := secrets.New()
	if e := instance.Walk(ctx, "/etc/secrets"); e != nil {
		panic(e)
	}

	for secret := range instance {
		keys := instance[secret]
		for key := range keys {
			_ = keys[key] // --> secret key's value
		}
	}

	service := instance["service"]

	_ = service["port"]
	_ = service["hostname"]
	_ = service["username"]
	_ = service["password"]
}
