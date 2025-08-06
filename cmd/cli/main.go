package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "altalune",
		Short: "altalune - go template",
		Long:  "altalune - altalune is a fullstack template for building go applications",
	}

	registerFlags(cmd)
	registerCommands(cmd)

	execCtx, execCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer execCancel()

	go func() {
		<-execCtx.Done()
		log.Println("🚨 interrupt/terminate signal received")
	}()

	if err := cmd.ExecuteContext(execCtx); err != nil {
		log.Printf("🔥 error execute the command: %v", err)
		os.Exit(1)
	}
}

func registerFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("config", "c", "config.yaml", "Configuration file path")
}

func registerCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewServeCommand(cmd),
		NewMigrateCommand(cmd),
	)
}
