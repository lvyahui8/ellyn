module github.com/lvyahui8/ellyn

go 1.18

require (
	github.com/emirpasic/gods v1.18.1
	github.com/stretchr/testify v1.9.0
	golang.org/x/mod v0.14.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require github.com/lvyahui8/ellyn/sdk v0.0.0

replace github.com/lvyahui8/ellyn/sdk => ./sdk

require github.com/lvyahui8/ellyn/api v0.0.0 // indirect

replace github.com/lvyahui8/ellyn/api => ./api
