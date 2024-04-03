package main

import (
	"context"
	"fmt"
	"os"
	"strings"
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

	generator func(context.Context) error
	linter    *dagger.Container
	examples  map[string]*dagger.Container
	tests     map[string]*dagger.Container
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

		// Create containers
		generator = ci.Generator(client)
		linter = ci.Linter(client)
		examples = ci.Examples(client, brokers)
		tests = ci.Tests(client, brokers)

		return nil
	},
}

var allCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "Execute all CI",
	Run: func(cmd *cobra.Command, args []string) {
		execute(context.Background(), generator)
		executeContainers(context.Background(), []*dagger.Container{linter})
		executeContainers(context.Background(), utils.MapToList(tests), utils.MapToList(examples))
	},
}

var examplesCmd = &cobra.Command{
	Use:     "examples",
	Aliases: []string{"g"},
	Short:   "Execute examples step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		if exampleFlag != "" {
			_, exists := examples[exampleFlag]
			if !exists {
				panic(fmt.Errorf("example %q doesn't exist", exampleFlag))
			}
			executeContainers(context.Background(), []*dagger.Container{examples[exampleFlag]})
		} else {
			executeContainers(context.Background(), utils.MapToList(examples))
		}
	},
}

var generatorCmd = &cobra.Command{
	Use:     "generator",
	Aliases: []string{"g"},
	Short:   "Execute generator step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		execute(context.Background(), generator)
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
		if testFlag != "" {
			if !strings.HasPrefix(testFlag, "./") {
				testFlag = "./" + testFlag
			}

			_, exists := tests[testFlag]
			if !exists {
				panic(fmt.Errorf("test %q doesn't exist in %+v", testFlag, tests))
			}
			executeContainers(context.Background(), []*dagger.Container{tests[testFlag]})
		} else {
			executeContainers(context.Background(), utils.MapToList(tests))
		}
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

func executeContainers(ctx context.Context, containers ...[]*dagger.Container) {
	// Regroup arg
	funcs := make([]func(context.Context) error, 0)
	for _, l1 := range containers {
		for _, l2 := range l1 {
			if l2 == nil {
				continue
			}

			// Note: create a new local variable to store value of actual l2
			callback := l2

			fn := func(ctx context.Context) error {
				_, err := callback.Stderr(ctx)
				return err
			}
			funcs = append(funcs, fn)
		}
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
