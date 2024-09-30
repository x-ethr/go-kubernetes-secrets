# `go-kubernetes-secrets`

`go-kubernetes-secrets` provides utilities for extracting volume-mounted Kubernetes secrets.

## Overview

Consider the following secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: ethr
data:
    host: ZXhhbXBsZS5ldGhyLmdnCg==
    port: ODA4MAo=
    username: c2VnbWVudGF0aW9uYWwK
    password: UEBzc3cwcmQK
type: Opaque
```

Mounted to the following deployment:

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
                        secretName: ethr
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
                            mountPath: /etc/secrets/ethr
```

The secret will be mounted `/etc/secrets`:

```
.
└── ethr
    ├── host        -> ..data/host
    ├── port        -> ..data/port
    ├── username    -> ..data/username
    └── password    -> ..data/password
```

Note that the references are *symbolic links*.

`go-kubernetes-secrets` will parse a given mount directory to return a [`Secrets`](./secrets.go) mapping, abstracting
the overhead of parsing a file-system, converting the binary contents of each key's file contents to a string, and
organizing the secret-to-key value(s):

```json
{
    "ethr": {
        "host": "example.ethr.gg",
        "port": "8080",
        "username": "segmentational",
        "password": "P@ssw0rd"
    }
}
```

## Documentation

Official `godoc` documentation (with examples) can be found at the [Package Registry](https://pkg.go.dev/github.com/x-ethr/go-kubernetes-secrets).

## Usage

###### Add Package Dependency

```bash
go get -u github.com/x-ethr/go-kubernetes-secrets
```

###### Import & Implement

`main.go`

```go
package main

import (
	"context"
	"fmt"

	"github.com/x-ethr/go-kubernetes-secrets"
)

func main() {
	ctx := context.Background()

	instance := secrets.New()
	if e := instance.Walk(ctx, "/etc/secrets"); e != nil {
		panic(e)
	}

	for secret, keys := range instance {
		for key, value := range keys {
			fmt.Println("Secret", secret, "Key", key, "Value", value)
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
- See https://pkg.go.dev/github.com/x-ethr/go-kubernetes-secrets for additional documentation.

## Contributions

See the [**Contributing Guide**](./CONTRIBUTING.md) for additional details on getting started.
