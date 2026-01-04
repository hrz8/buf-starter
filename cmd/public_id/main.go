package main

import (
	"fmt"
	"os"

	"github.com/hrz8/altalune/internal/shared/nanoid"
	"github.com/spf13/cobra"
)

var (
	batch bool
	count int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "publicid",
		Short: "Generate Public IDs",
		Long:  "Generate NanoID-based Public IDs in single or batch mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if batch {
				ids, err := nanoid.GeneratePublicIDBatch(count)
				if err != nil {
					return err
				}

				for _, id := range ids {
					fmt.Println(id)
				}
				return nil
			}

			id, err := nanoid.GeneratePublicID()
			if err != nil {
				return err
			}

			fmt.Println(id)
			return nil
		},
	}

	rootCmd.Flags().BoolVarP(
		&batch,
		"batch",
		"b",
		false,
		"Generate IDs in batch mode",
	)

	rootCmd.Flags().IntVarP(
		&count,
		"count",
		"c",
		1,
		"Number of IDs to generate (used with --batch)",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
