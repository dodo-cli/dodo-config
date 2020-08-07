package command

import (
	"fmt"

	"github.com/oclaussen/dodo/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config plugin subcommands",
	}

	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewValidateCommand())

	return cmd
}

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Short:                 "List available backdrop configurations",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			backdrops := map[string]*types.Backdrop{}
			configfiles.GimmeConfigFiles(&configfiles.Options{
				Name:                      "dodo",
				Extensions:                []string{"yaml", "yml", "json"},
				IncludeWorkingDirectories: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					d := decoder.New(configFile.Path)
					d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
						"backdrops": decoder.Map(types.NewBackdrop(), &backdrops),
					})

					return false
				},
			})

			for name := range backdrops {
				// TODO filename
				fmt.Println(name)
			}

			return nil
		},
	}
}

func NewValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "validate",
		Short:                 "Validate configuration files for syntax errors",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			backdrops := map[string]*types.Backdrop{}
			configfiles.GimmeConfigFiles(&configfiles.Options{
				FileGlobs:        args,
				UseFileGlobsOnly: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					d := decoder.New(configFile.Path)
					d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
						"backdrops": decoder.Map(types.NewBackdrop(), &backdrops),
					})

					for _, err := range d.Errors() {
						fmt.Println(err)
					}

					return false
				},
			})

			return nil
		},
	}
}