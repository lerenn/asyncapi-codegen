package main

import (
	"fmt"
	"os"

	"github.com/TheSadlig/asyncapi-codegen/pkg/codegen"
	"github.com/spf13/cobra"
)

var flags Flags

var cmd = &cobra.Command{
	Use:   "asyncapi-codegen",
	Short: "An AsyncAPI Code generator that generates all code from the broker to the application/user.",
	Long: `An AsyncAPI Golang Code generator that generates all Go code from the broker to the application/user. 
Just plug your application to your favorite message broker!

More info on README: https://github.com/TheSadlig/asyncapi-codegen
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cg, err := codegen.FromFile(flags.InputPaths[0], flags.InputPaths[1:]...)
		if err != nil {
			return err
		}

		opt, err := flags.ToCodegenOptions()
		if err != nil {
			return err
		}

		if err := cg.Generate(opt); err != nil {
			return err
		}

		return nil
	},
}

func main() {
	flags.SetToCommand(cmd)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
