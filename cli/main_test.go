package main

import (
	"github.com/bu3/rebrickable-cli/cli/cmd"
	"github.com/rogpeppe/go-internal/testscript"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"cli": cmd.Execute,
	}))
}

func TestCli(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:   "../testdata", //TODO pass data with bazel
		Setup: setupTestEnv,
	})
}

func setupTestEnv(env *testscript.Env) error {
	// Enable new `devbox run` so we can use it in tests. This is temporary,
	// and should be removed once we enable this feature flag.
	env.Setenv("REBRICKABLE_USERNAME", os.Getenv("REBRICKABLE_USERNAME"))
	env.Setenv("REBRICKABLE_PASSWORD", os.Getenv("REBRICKABLE_PASSWORD"))
	env.Setenv("REBRICKABLE_API_KEY", os.Getenv("REBRICKABLE_API_KEY"))

	return nil
}
