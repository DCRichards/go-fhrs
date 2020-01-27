# go-fhrs

> 🍽  A Go client library for the UK [Food Hygiene Rating Scheme API](https://api.ratings.food.gov.uk/help).

## Usage

Simply `go get` or add it to your imports and run any `go` command to auto-update your modules.

```bash
go get -u github.com/dcrichards/go-fhrs
```

```go
package main

import (
        "github.com/dcrichards/go-fhrs/fhrs"
)

func main() {
        client, err := fhrs.NewClient()
        if err != nil {
                // Handle err
        }
        
        est, err := client.Establishments.GetByID("82940")
        if err != nil {
                // Handle err
        }
        
        // Do stuff with est
}
```

## Examples

An example can be found in the `example` directory.

## Testing

Tests use a mock server and do not require access to the API.

**Docker**
```bash
docker-compose run --rm go test ./...
# With coverage
docker-compose run --rm go test -coverprofile=coverage.out ./...
```

**No Docker**

```bash
CGO_ENABLED=0 go test ./...
# With coverage
CGO_ENABLED=0 go test -coverprofile=coverage.out ./...
```
