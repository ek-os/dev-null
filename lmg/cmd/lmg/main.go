package main

import (
	"context"
	"os"

	"github.com/ek-os/lmg/internal/lmg"
)

func main() {
	lmg.Run(context.Background(), os.Getenv, os.Stdout)
}
