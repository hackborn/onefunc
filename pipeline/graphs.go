package pipeline

var pipeline1 = `
graph (
	src -> go -> folder
	src - > sql -d> go
)
`

var fullTestPipeline1 = `
graph (
	src -> go -> folder
	src - > sql -d> go -> "->me" -> go/home
)
`

var driverPipeline = `
graph (
	src -> go -> folder
	src -> sql -> go
)

vars (
)
`
