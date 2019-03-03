package shellcmd

import "testing"

func TestKillAction(t *testing.T) {
	t.SkipNow()
	KillAction(nil, []string{"tidb-docker-compose_tikv0_1"})
}