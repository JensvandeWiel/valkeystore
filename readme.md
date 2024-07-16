[![Testing](https://github.com/JensvandeWiel/valkeystore/actions/workflows/test.yml/badge.svg)](https://github.com/JensvandeWiel/valkeystore/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/pkg.go.dev/github.com/JensvandeWiel/valkeystore.svg)](https://pkg.go.dev/pkg.go.dev/github.com/JensvandeWiel/valkeystore)
# [gorilla/sessions](https://github.com/gorilla/sessions) Store implementation for [valkey](https://github.com/valkey-io/valkey)
This implementation uses [valkey-io/go-valkey](https://github.com/valkey-io/valkey-go) as the client.

## Usage
```bash
go get github.com/JensvandeWiel/valkeystore
```

```go
package main

import (
	"github.com/JensvandeWiel/valkeystore"
	"github.com/gorilla/sessions"
	"github.com/valkey-io/valkey-go"
)

func main() {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal("failed to create valkey client", err)
	}

	defer client.Close()

	store, err := valkeystore.NewValkeyStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}
}
```
For more examples check the tests or godoc

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Contributing
Feel free to open a pull request.

## This project is inspired by:
- https://github.com/rbcervilla/redisstore
- https://github.com/boj/redistore