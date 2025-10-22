package cmd

import (
	"easy-attend-service/configs"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Migrate Command
func Migrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "migrate",
		Args: NotReqArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configs.ConnectDatabase()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			migrateUp().Run(cmd, args)
		},
	}
	cmd.AddCommand(migrateUp())
	cmd.AddCommand(migrateDown())
	cmd.AddCommand(migrateRefresh())
	cmd.AddCommand(migrateIntID())
	return cmd
}

func migrateUp() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "up",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := modelUp(); err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Println("Migration up completed successfully!")
			os.Exit(0)
		},
	}
	return cmd
}

func migrateDown() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "down",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := modelDown(); err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Println("Migration down completed successfully!")
			os.Exit(0)
		},
	}
	return cmd
}

func migrateRefresh() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "refresh",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := modelDown(); err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			if err := modelUp(); err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Println("Migration refresh completed successfully!")
			os.Exit(0)
		},
	}
	return cmd
}

func migrateIntID() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "int-id",
		Args: NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := migrateToIntID(); err != nil {
				fmt.Printf("Migration to int ID failed: %s", err)
				os.Exit(1)
			}
			fmt.Println("Migration to int ID completed successfully!")
			os.Exit(0)
		},
	}
	return cmd
}
