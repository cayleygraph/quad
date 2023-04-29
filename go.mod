module github.com/cayleygraph/quad

go 1.19

replace google.golang.org/protobuf => github.com/aperturerobotics/protobuf-go v1.30.1-0.20230428014030-7089409cbc63 // aperture

require (
	github.com/piprate/json-gold v0.5.0
	github.com/stretchr/testify v1.8.2
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
