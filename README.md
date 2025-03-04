# `go-kubernetes-secrets`

`go-kubernetes-secrets` is a zero-dependency package that provides utilities for extracting volume-mounted Kubernetes secrets.

By using `go-kubernetes-secrets`, Kubernetes workload(s) receive automatic updates whenever a `Secret` is
modified -- avoiding restarts, redeployments, and base64 decoding.

## Overview

Consider the `Secret`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: test-data
data:
    host: ZXhhbXBsZS5ldGhyLmdnCg==
    port: ODA4MAo=
    username: c2VnbWVudGF0aW9uYWwK
    password: UEBzc3cwcmQK
type: Opaque
```

Configured with the `Deployment`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
    name: test-secrets
    labels:
        app: test-secrets
spec:
    replicas: 1
    selector:
        matchLabels:
            app: test-secrets
    template:
        metadata:
            labels:
                app: test-secrets
        spec:
            volumes:
                -   name: secret-volume
                    secret:
                        secretName: test-data
                        optional: false
            containers:
                -   name: test-secrets
                    image: service:latest
                    imagePullPolicy: Always
                    ports:
                        -   containerPort: 8080
                    volumeMounts:
                        -   name: secret-volume
                            readOnly: true
                            mountPath: /etc/secrets/test-data
```

The `Secret` will be mounted within the `/etc/secrets` directory:

```
.
└── test-data
    ├── host        -> ..data/host
    ├── port        -> ..data/host
    ├── username    -> ..data/host
    └── password    -> ..data/password
```

> [!NOTE]
> The `host`, `port`, `username`, and `password` file(s) are *symbolic links*, and represent the individual key-value pairs
> of the `test-data`'s `Secret`.

`go-kubernetes-secrets` will parse a given mount directory to return a [`Secrets`](./secrets.go) mapping, abstracting
the overhead of parsing a file-system, converting the binary contents of each key's file contents to a string, and
organizing the secret-to-key value(s):

```json
{
    "test-data": {
        "host": "polygun.com",
        "port": "8080",
        "username": "segmentational",
        "password": "P@ssw0rd"
    }
}
```

## Documentation

Official `godoc` documentation (with examples) can be found at the [Package Registry](https://pkg.go.dev/github.com/poly-gun/go-kubernetes-secrets).

## Usage

###### Add Package Dependency

```bash
go get -u github.com/poly-gun/go-kubernetes-secrets
```

##### Import & Implement

`main.go`

###### Directory Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/poly-gun/go-kubernetes-secrets"
)

func main() {
    ctx := context.Background()

    instance := secrets.New()
    e := instance.WalkWithContext(ctx, func(o *secrets.Options) {
        o.Directory = "/etc/secrets" // Directory Usage Example
    })

    if e != nil {
        panic(e)
    }

    for secret, keys := range instance {
        for key, value := range keys {
            fmt.Println("Secret", secret, "Key", key, "Item", value)
        }
    }

    service := instance["service"]

    port := service["port"]
    hostname := service["hostname"]
    username := service["username"]
    password := service["password"]

    fmt.Println("Port", port, "Hostname", hostname, "Username", username, "Password", password)
}

```

###### FS Usage

```go
package main

import (
    "context"
    "embed"
    "fmt"

    "github.com/poly-gun/go-kubernetes-secrets"
)

//go:embed test-data
var filesystem embed.FS

func main() {
    ctx := context.Background()

    instance := secrets.New()
    e := instance.WalkWithContext(ctx, func(o *secrets.Options) {
        o.FS = filesystem // FS Usage Example
    })

    if e != nil {
        panic(e)
    }

    for secret, keys := range instance {
        for key, value := range keys {
            fmt.Println("Secret", secret, "Key", key, "Item", value)
        }
    }

    service := instance["service"]

    port := service["port"]
    hostname := service["hostname"]
    username := service["username"]
    password := service["password"]

    fmt.Println("Port", port, "Hostname", hostname, "Username", username, "Password", password)
}
```

- Please refer to the [code examples](./example_test.go) for additional usage and implementation details.
- See https://pkg.go.dev/github.com/poly-gun/go-kubernetes-secrets for additional documentation.

## Contributions

See the [**Contributing Guide**](./CONTRIBUTING.md) for additional details on getting started.
