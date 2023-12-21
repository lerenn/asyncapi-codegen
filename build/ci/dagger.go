///usr/local/bin/dagger run go run $0 $@ ; exit

package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"dagger.io/dagger"
	"github.com/lerenn/asyncapi-codegen/pkg/pipeline"
	"github.com/spf13/cobra"
)

var (
	client *dagger.Client

	brokers map[string]*dagger.Service

	generator *dagger.Container
	linter    *dagger.Container
	examples  []*dagger.Container
	tests     []*dagger.Container
)

var rootCmd = &cobra.Command{
	Use:   "./build/ci/dagger.go",
	Short: "A simple CLI to execute asyncapi-codegen project CI/CD with dagger",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Initialize Dagger client
		client, err = dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stderr))
		if err != nil {
			return err
		}
		defer client.Close()

		// Create services
		brokers = pipeline.Brokers(client)

		// Create containers
		generator = pipeline.Generator(client)
		linter = pipeline.Linter(client)
		examples = pipeline.Examples(client, brokers)
		tests = pipeline.Tests(client, brokers)

		return nil
	},
}

var allCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "Execute all CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), []*dagger.Container{generator, linter})
		executeContainers(context.Background(), tests, examples)
	},
}

var examplesCmd = &cobra.Command{
	Use:     "examples",
	Aliases: []string{"g"},
	Short:   "Execute examples step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), examples)
	},
}

var generatorCmd = &cobra.Command{
	Use:     "generator",
	Aliases: []string{"g"},
	Short:   "Execute generator step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), []*dagger.Container{generator})
	},
}

var linterCmd = &cobra.Command{
	Use:     "linter",
	Aliases: []string{"g"},
	Short:   "Execute linter step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), []*dagger.Container{linter})
	},
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"g"},
	Short:   "Execute tests step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), tests)
	},
}

func main() {
	rootCmd.AddCommand(allCmd)
	rootCmd.AddCommand(examplesCmd)
	rootCmd.AddCommand(generatorCmd)
	rootCmd.AddCommand(linterCmd)
	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func executeContainers(ctx context.Context, containers ...[]*dagger.Container) {
	// Regroup arg
	rContainers := make([]*dagger.Container, 0)
	for _, c := range containers {
		rContainers = append(rContainers, c...)
	}

	// Excute containers
	var wg sync.WaitGroup
	for _, ec := range rContainers {
		go func(e *dagger.Container) {
			_, err := e.Stderr(ctx)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(ec)

		wg.Add(1)
	}

	wg.Wait()
}
