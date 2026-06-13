package docs

import (
	"encoding/json"
	"strings"

	"github.com/urfave/cli/v3"
)

type CLICommand struct {
	Name        string       `json:"name"`
	Aliases     []string     `json:"aliases,omitempty"`
	Usage       string       `json:"usage,omitempty"`
	UsageText   []string     `json:"usageText,omitempty"`
	Description string       `json:"description,omitempty"`
	ArgsUsage   string       `json:"argsUsage,omitempty"`
	Category    string       `json:"category,omitempty"`
	Version     string       `json:"version,omitempty"`
	Hidden      bool         `json:"hidden,omitempty"`
	Flags       []CLIFlag    `json:"flags,omitempty"`
	Commands    []CLICommand `json:"commands,omitempty"`
}

type CLIFlag struct {
	Name       string   `json:"name"`
	Aliases    []string `json:"aliases,omitempty"`
	Usage      string   `json:"usage,omitempty"`
	Type       string   `json:"type,omitempty"`
	Default    string   `json:"default,omitempty"`
	TakesValue bool     `json:"takesValue"`
	EnvVars    []string `json:"envVars,omitempty"`
	Hidden     bool     `json:"hidden,omitempty"`
	Required   bool     `json:"required,omitempty"`
}

func ToJSON(cmd *cli.Command) ([]byte, error) {
	out := prepareJSONCommand(cmd)
	return json.MarshalIndent(out, "", "  ")
}

func prepareJSONCommand(cmd *cli.Command) CLICommand {
	tt := tabularTemplate{}

	return CLICommand{
		Name:        cmd.Name,
		Aliases:     cmd.Aliases,
		Usage:       tt.PrepareMultilineString(cmd.Usage),
		UsageText:   prepareUsageTextLines(cmd.UsageText),
		Description: tt.PrepareMultilineString(cmd.Description),
		ArgsUsage:   tt.PrepareMultilineString(cmd.ArgsUsage),
		Category:    cmd.Category,
		Version:     cmd.Version,
		Hidden:      cmd.Hidden,
		Flags:       prepareJSONFlags(cmd.Flags),
		Commands:    prepareJSONCommands(cmd.Commands),
	}
}

func prepareJSONCommands(commands []*cli.Command) []CLICommand {
	result := make([]CLICommand, 0, len(commands))
	for _, cmd := range commands {
		result = append(result, prepareJSONCommand(cmd))
	}
	return result
}

func prepareJSONFlags(flags []cli.Flag) []CLIFlag {
	result := make([]CLIFlag, 0, len(flags))
	tt := tabularTemplate{}

	for _, f := range flags {
		flag, ok := f.(cli.DocGenerationFlag)
		if !ok {
			continue
		}

		value, defaultText := getFlagDefaultValue(flag)
		defaultValue := ""
		if defaultText != "" {
			defaultValue = defaultText
		} else if value != "" {
			defaultValue = value
		}

		jf := CLIFlag{
			Usage:      tt.PrepareMultilineString(flag.GetUsage()),
			EnvVars:    flag.GetEnvVars(),
			TakesValue: flag.TakesValue(),
			Default:    defaultValue,
			Type:       flag.TypeName(),
		}

		names := f.Names()
		for i, name := range names {
			name = strings.TrimSpace(name)
			if i == 0 {
				jf.Name = "--" + name
				continue
			}
			if len(name) > 1 {
				name = "--" + name
			} else {
				name = "-" + name
			}
			jf.Aliases = append(jf.Aliases, name)
		}

		if vf, ok := f.(cli.VisibleFlag); ok {
			jf.Hidden = !vf.IsVisible()
		}

		if rf, ok := f.(interface{ IsRequired() bool }); ok {
			jf.Required = rf.IsRequired()
		}

		result = append(result, jf)
	}

	return result
}

func prepareUsageTextLines(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.Trim(s, "\n")
	lines := strings.Split(s, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		filtered = append(filtered, trimmed)
	}
	return filtered
}
