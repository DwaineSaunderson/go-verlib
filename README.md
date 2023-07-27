# go-verlib

`go-verlib` is a powerful and flexible library for parsing, comparing, and manipulating version numbers in Go, adhering
to the Semantic Versioning 2.0.0 standard.

## Installation

To use the package, import it in your Go code:

```go
import "github.com/DwaineSaunderson/go-verlib"
```

## Usage

Here are a few examples of what you can do with `go-verlib`:

### Parsing a version

```go
v, err := verlib.ParseVersion("1.2.3-beta.1+5678")
if err != nil {
log.Fatalf("Failed to parse version: %v", err)
}
fmt.Printf("Parsed version: %s\n", v.String())
```

### Strict parsing of a version

This ensures that the version string is valid according to Semantic Versioning 2.0.0.

```go
v, err := verlib.ParseVersionStrict("1.2.3-beta.1+5678")
if err != nil {
log.Fatalf("Failed to parse version: %v", err)
}
fmt.Printf("Parsed version: %s\n", v.StrictString())
```

### Creating and comparing versions

```go
v1 := verlib.NewVersion(1, 0, 0)
v2 := verlib.NewVersion(2, 0, 0)
fmt.Printf("v1 is less than v2: %t\n", v1.Less(v2)) // Output: true
```

### Working with constraints

```go
c := verlib.NewConstraint(verlib.GT, v1)
fmt.Printf("v2 satisfies constraint c: %t\n", v2.Satisfies(c)) // Output: true
```

### Parsing and using constraints

```go
c, err := verlib.ParseConstraint(">= 1.2.0")
if err != nil {
log.Fatalf("Failed to parse constraint: %v", err)
}
fmt.Printf("v2 satisfies constraint c: %t\n", v2.Satisfies(c)) // Output: true
```

### Checking for contradictory constraints

```go
c1, err := verlib.ParseConstraint(">= 2.0.0")
c2, err := verlib.ParseConstraint("< 1.5.0")
constraints := verlib.Constraints{c1, c2}
err = constraints.Contradicts()
if err != nil {
log.Fatalf("Constraints are contradictory: %v", err)
}
```

### Incrementing a version

```go
v := verlib.NewVersion(1, 2, 3)
v.IncrementMajor()
fmt.Printf("Incremented version: %s\n", v.String()) // Output: 2.0.0
```

### Using pre-release and build metadata

```go
v := verlib.NewVersionWithPreReleaseAndMetadata(1, 2, 3, "beta.1", "5678")
fmt.Printf("Version with pre-release and build metadata: %s\n", v.String()) // Output: 1.2.3-beta.1+5678
```

## Contributing

Contributions are welcome! Feel free to open a pull request on Github.

## License

MIT License

## Author

Dwaine Saunderson
