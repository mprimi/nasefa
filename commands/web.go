package commands

import (
  "context"
  "flag"
  "fmt"
  "github.com/google/subcommands"
)

const (
  kDefaultBindAddr = ":8080"
)

type webCommand struct {
  bindAddr    string
}

func WebCommand() (subcommands.Command) {
  return &webCommand{}
}

func (*webCommand) Name() string     { return "web" }
func (*webCommand) Synopsis() string { return "Starts web application" }
func (*webCommand) Usage() string { return "web [options]\n" }
func (this *webCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bindAddr, "bindAddr", kDefaultBindAddr, "Address to bind")
}

func (this *webCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if len(f.Args()) > 0 {
    fmt.Printf("⚠️ Usage error: unknown arguments: %v\n", f.Args())
    return subcommands.ExitUsageError
  }

  WebAppStart(this.bindAddr)

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
