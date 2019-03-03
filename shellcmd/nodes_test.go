package shellcmd

import "testing"

func TestNodesAction(t *testing.T) {
	NodesAction(nil, make([]string, 0))
}
