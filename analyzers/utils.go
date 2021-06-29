package analyzers

type traversalAction = string

const (
	Leave traversalAction = "leave"
	Visit traversalAction = "visit"
)

// A _very_ simple helper to use switch cases in go ast traversal
// using the node "visit" and "leave" verbs. `push bool` is taken from
// the ast/inspector .Nodes tree traversal function.
func Action(push bool) traversalAction {
	if push {
		return Visit
	} else {
		return Leave
	}
}
