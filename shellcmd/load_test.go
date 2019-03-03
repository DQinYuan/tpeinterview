package shellcmd

import "testing"

func TestLoadAction(t *testing.T) {
	LoadAction(nil, []string{"*"})
}
