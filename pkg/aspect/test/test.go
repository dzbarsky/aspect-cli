/*
Copyright © 2021 Aspect Build Systems Inc

Not licensed for re-use.
*/

package test

import (
	"fmt"

	"aspect.build/cli/pkg/aspecterrors"
	"aspect.build/cli/pkg/bazel"
	"aspect.build/cli/pkg/ioutils"
	"aspect.build/cli/pkg/plugin/system/bep"
)

type Test struct {
	ioutils.Streams
	bzl bazel.Bazel
}

func New(streams ioutils.Streams, bzl bazel.Bazel) *Test {
	return &Test{
		Streams: streams,
		bzl:     bzl,
	}
}

func (t *Test) Run(args []string, besBackend bep.BESBackend) (exitErr error) {
	besBackendFlag := fmt.Sprintf("--bes_backend=grpc://%s", besBackend.Addr())
	bazelCmd := []string{"test", besBackendFlag}
	bazelCmd = append(bazelCmd, args...)

	exitCode, bazelErr := t.bzl.Spawn(bazelCmd)

	// Process the subscribers errors before the Bazel one.
	subscriberErrors := besBackend.Errors()
	if len(subscriberErrors) > 0 {
		for _, err := range subscriberErrors {
			fmt.Fprintf(t.Streams.Stderr, "Error: failed to run test command: %v\n", err)
		}
		exitCode = 1
	}

	if exitCode != 0 {
		err := &aspecterrors.ExitError{ExitCode: exitCode}
		if bazelErr != nil {
			err.Err = bazelErr
		}
		return err
	}

	return nil
}