package metrics

type Metrics struct {
	Resources string
	Operation Operation
}

type Operation struct {
	Push string
}
