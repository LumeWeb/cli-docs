package docs

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestToJSONFull(t *testing.T) {
	cmd := buildExtendedTestCommand(t)
	res, err := ToJSON(cmd)
	require.NoError(t, err)
	expectFileContent(t, "testdata/expected-json-full.json", string(res))
}

func TestToJSONAllFlagTypes(t *testing.T) {
	cmd := buildFlagTest(t)
	res, err := ToJSON(cmd)
	require.NoError(t, err)
	expectFileContent(t, "testdata/expected-json-all-flag-types.json", string(res))
}

func TestToJSONNoFlags(t *testing.T) {
	cmd := buildExtendedTestCommand(t)
	cmd.Flags = nil
	res, err := ToJSON(cmd)
	require.NoError(t, err)
	expectFileContent(t, "testdata/expected-json-no-flags.json", string(res))
}

func TestToJSONNoCommands(t *testing.T) {
	cmd := buildExtendedTestCommand(t)
	cmd.Commands = nil
	res, err := ToJSON(cmd)
	require.NoError(t, err)
	expectFileContent(t, "testdata/expected-json-no-commands.json", string(res))
}

func TestToJSONEmptyCommand(t *testing.T) {
	cmd := &cli.Command{Name: "empty"}
	res, err := ToJSON(cmd)
	require.NoError(t, err)
	expectFileContent(t, "testdata/expected-json-empty.json", string(res))
}

func TestToJSONHiddenCommand(t *testing.T) {
	cmd := &cli.Command{
		Name:   "root",
		Hidden: false,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "visible-flag", Usage: "shown"},
			&cli.StringFlag{Name: "hidden-flag", Usage: "not shown", Hidden: true},
		},
		Commands: []*cli.Command{
			{Name: "visible-cmd", Usage: "shown"},
			{Name: "hidden-cmd", Usage: "not shown", Hidden: true},
		},
	}
	res, err := ToJSON(cmd)
	require.NoError(t, err)
	expectFileContent(t, "testdata/expected-json-hidden.json", string(res))
}
