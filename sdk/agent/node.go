package agent

type node struct {
	methodId uint32
	blocks   []bool
	cost     int64
	args     []any
	results  []any
}
