module github.com/lvyahui8/ellyn/ellyn

go 1.18

require (
	github.com/lvyahui8/ellyn v0.0.0
	github.com/urfave/cli/v2 v2.27.5
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/mod v0.14.0 // indirect
)

replace github.com/lvyahui8/ellyn => ../

require github.com/lvyahui8/ellyn/api v0.0.0 // indirect

replace github.com/lvyahui8/ellyn/api => ../api
