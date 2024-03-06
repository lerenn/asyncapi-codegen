package ci

import (
	"context"

	"dagger.io/dagger"
)

// Generator returns a container that generates code.
func Generator(client *dagger.Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		_, err := client.Container().
			// Add base image
			From(GolangImage).
			// Add source code as work directory
			With(sourceAsWorkdir(client)).
			// Add command to generate code
			WithExec([]string{"go", "generate", "./..."}).
			// Export directory
			Directory(".").
			Export(ctx, ".")

		return err
	}
}
