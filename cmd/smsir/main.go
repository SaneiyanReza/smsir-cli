package main

import (
	"github.com/SaneiyanReza/smsir-cli/cmd/smsir/commands"
)

func main() {
	// Delegate entirely to cobra; users can run `smsir menu` to start UI
	commands.Execute()
}
