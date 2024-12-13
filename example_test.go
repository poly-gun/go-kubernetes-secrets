package secrets_test

import (
	"context"

	"github.com/poly-gun/go-kubernetes-secrets"
)

func Example() {
	ctx := context.Background()

	instance := secrets.New()
	e := instance.WalkWithContext(ctx, func(o *secrets.Options) {
		o.Directory = "/etc/secrets"
	})

	if e != nil {
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
