package main

import "testing"

func TestDockerCompose(t *testing.T) {
	t.SkipNow()
	dockerCompose("/home/dqyuan/language/Docker/projects/tidb-docker-compose/docker-compose.yml", "up")
}
