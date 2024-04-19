package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"dagger.io/dagger"
	"github.com/lerenn/asyncapi-codegen/pkg/ci"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	client *dagger.Client

	exampleFlag string
	tagFlag     string
	testFlag    string

	brokers map[string]*dagger.Service
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
		brokers = ci.Brokers(client)

		return nil
	},
}

var allCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "Execute all CI",
	Run: func(cmd *cobra.Command, args []string) {
		execute(context.Background(), ci.Generator(client))
		executeContainers(context.Background(), ci.Linter(client))
		executeContainers(context.Background(), ci.Tests(client, brokers, "./..."))
		executeContainers(context.Background(), ci.Tests(client, brokers, "./..."))
	},
}

var examplesCmd = &cobra.Command{
	Use:     "examples",
	Aliases: []string{"g"},
	Short:   "Execute examples step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		examples := ci.Examples(client, brokers)

		if exampleFlag != "" {
			_, exists := examples[exampleFlag]
			if !exists {
				panic(fmt.Errorf("example %q doesn't exist", exampleFlag))
			}
			executeContainers(context.Background(), examples[exampleFlag])
		} else {
			executeContainers(context.Background(), utils.MapToList(examples)...)
		}
	},
}

var generatorCmd = &cobra.Command{
	Use:     "generator",
	Aliases: []string{"g"},
	Short:   "Execute generator step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		execute(context.Background(), ci.Generator(client))
	},
}

var linterCmd = &cobra.Command{
	Use:     "linter",
	Aliases: []string{"g"},
	Short:   "Execute linter step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), ci.Linter(client))
	},
}

var publishCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"p"},
	Short:   "Tag and publish to different repositories.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ci.Publish(context.Background(), client, tagFlag)
	},
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"g"},
	Short:   "Execute tests step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		if testFlag == "" {
			testFlag = "./..."
		}

		executeContainers(context.Background(), ci.Tests(client, brokers, testFlag))
	},
}

func main() {
	rootCmd.AddCommand(allCmd)

	examplesCmd.Flags().StringVarP(&exampleFlag, "example", "e", "", "Example to execute")
	rootCmd.AddCommand(examplesCmd)

	rootCmd.AddCommand(generatorCmd)

	rootCmd.AddCommand(linterCmd)

	publishCmd.Flags().StringVarP(&tagFlag, "tag", "t", "", "Tag used to tag this version")
	rootCmd.AddCommand(publishCmd)

	testCmd.Flags().StringVarP(&testFlag, "test", "t", "", "Test to execute")
	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func executeContainers(ctx context.Context, containers ...*dagger.Container) {
	// Regroup arg
	funcs := make([]func(context.Context) error, 0)
	for _, l1 := range containers {
		if l1 == nil {
			continue
		}

		// Note: create a new local variable to store value of actual l2
		callback := l1

		fn := func(ctx context.Context) error {
			_, err := callback.Stderr(ctx)
			return err
		}
		funcs = append(funcs, fn)
	}

	execute(ctx, funcs...)
}

func execute(ctx context.Context, funcs ...func(context.Context) error) {
	// Excute containers
	var wg sync.WaitGroup
	for _, fn := range funcs {
		go func(callback func(context.Context) error) {
			if err := callback(ctx); err != nil {
				panic(err)
			}
			wg.Done()
		}(fn)

		wg.Add(1)
	}

	wg.Wait()
}
