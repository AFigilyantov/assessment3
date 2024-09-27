package main

import (
	"chitests/cmd/commands"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	cmd := commands.NewServeCmd()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("service is down by error: %s", err)
	}

}
