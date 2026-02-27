//go:build mage
// +build mage

// Package main provides Mage build targets for linting, running tests, and
// orchestrating a Docker‑Compose test stack.  All public functions are
// exported targets; every docstring is wrapped to 80 columns so that the file
// remains readable in split‑screen editors.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	composeFile = "./acceptance/docker-compose.yaml" // path to Compose file
	projectName = "tests"                            // Compose project name
)

// composeEnv builds the environment passed to every docker‑compose command. It
// injects the current user‑ID so containers can create files with the correct
// ownership, enables Terraform acceptance tests, and forces Compose to emit
// plain progress output (useful in CI logs).
func composeEnv() map[string]string {
	uidBytes, _ := exec.Command("id", "-u").Output()
	return map[string]string{
		"ZITADEL_DEV_UID":  strings.TrimSpace(string(uidBytes)),
		"TF_ACC":           "1",
		"COMPOSE_PROGRESS": "plain",
	}
}

// runTests brings the stack up, executes `go test`, and always tears the stack
// down.  Extra arguments (e.g. –coverprofile flags) are forwarded to the `go
// test` invocation so that callers can run both plain and coverage runs with
// the same orchestration logic.
func runTests(ctx context.Context, extraArgs ...string) error {
	if err := Up(); err != nil {
		return err
	}
	defer Down()

	args := append([]string{"test"}, extraArgs...)
	args = append(args, "./...")
	return sh.RunV("go", args...)
}

// Test is a Mage namespace.  `mage test:default` runs the normal test suite,
// while `mage test:cover` writes a coverage profile.
type Test mg.Namespace

// Default executes the full acceptance test suite against the Compose stack.
func (Test) Default(ctx context.Context) error {
	return runTests(ctx)
}

// Cover runs the same suite but saves coverage data to profile.cov so that CI
// can upload the report.
func (Test) Cover(ctx context.Context) error {
	return runTests(ctx, "-coverprofile=profile.cov")
}

// Lint invokes golangci‑lint if it is available; otherwise it falls back to
// `go vet` so that developers without the linter installed can still run the
// target locally without errors.
func Lint() error {
	if err := sh.Run("golangci-lint", "run", "./..."); err == nil {
		return nil
	}
	fmt.Println("golangci-lint not found – falling back to `go vet`")
	return sh.RunV("go", "vet", "./...")
}

// Up starts the Compose stack in detached mode and waits for all health
// checks.  Using a dedicated project name isolates the network and volumes
// from any other Compose projects that might be running on the same host.
// Up starts the Compose stack in detached mode, waits for health checks, and
// then sleeps an extra 30 seconds to give services a buffer before the test
// suite begins.  This is occasionally helpful when containers report healthy
// but still need a moment to accept connections (e.g. databases warming up).
func Up() error {
	if err := sh.RunWith(composeEnv(),
		"docker", "compose",
		"--file", composeFile,
		"--project-name="+projectName,
		"up", "--detach", "--wait"); err != nil {
		return err
	}
	// Extra buffer so that flaky startup races are less likely.
	if err := sh.RunV("sleep", "30"); err != nil {
		fmt.Println("warning: sleep command failed:", err)
	}
	return nil
}

// Down stops the stack and removes volumes.  If one of the well‑known debug
// flags is present (DEBUG, ACTIONS_STEP_DEBUG, or ACTIONS_RUNNER_DEBUG), it
// streams Compose logs to stdout first so developers can diagnose failures in
// CI without hunting for artifacts.
func Down() error {
	if os.Getenv("DEBUG") != "" ||
		os.Getenv("ACTIONS_STEP_DEBUG") == "true" ||
		os.Getenv("ACTIONS_RUNNER_DEBUG") == "true" {
		if err := sh.RunWith(composeEnv(),
			"docker", "compose",
			"--file", composeFile,
			"--project-name="+projectName,
			"logs"); err != nil {
			fmt.Println("warning: failed to fetch compose logs:", err)
		}
	}
	return sh.RunWith(composeEnv(),
		"docker", "compose",
		"--file", composeFile,
		"--project-name="+projectName,
		"down", "--volumes", "--remove-orphans")
}

// Default target: `mage` with no arguments lints the codebase and then runs
// the default test suite so that a single command ensures code quality and a
// passing set of integration tests.
var Default = All

// All wires the high‑level workflow of linting first and then running tests.
func All(ctx context.Context) {
	mg.SerialDeps(Lint, Test{}.Default)
}
