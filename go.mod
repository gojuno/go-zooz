module github.com/gtforge/go-zooz

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.2.0
	github.com/stretchr/testify v1.7.0
)

retract v1.3.1 // Contains errors with deserialization
