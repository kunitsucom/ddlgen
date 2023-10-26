package main

import (
	"context"
	"log"

	"github.com/kunitsucom/ddlgen/pkg/ddlgen"
)

func main() {
	ctx := context.Background()

	if err := ddlgen.DDLGen(ctx); err != nil {
		log.Fatalf("ddlgen: %+v", err)
	}
}
