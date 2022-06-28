package commands

import (
  "context"
  "flag"
  "github.com/google/subcommands"
)

const (
  kDefaultBindAddr = ":8080"
  kDefaultPrefixPath = "/"
)

type webCommand struct {
  bindAddr    string
  prefixPath  string
}

func WebCommand() (subcommands.Command) {
  return &webCommand{}
}

func (*webCommand) Name() string     { return "web" }
func (*webCommand) Synopsis() string { return "Starts web application" }
func (*webCommand) Usage() string { return "web [options]\n" }
func (this *webCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bindAddr, "bindAddr", kDefaultBindAddr, "Address to bind")
  f.StringVar(&this.prefixPath, "prefix", kDefaultPrefixPath, "HTTP path prefix")
}

func (this *webCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if len(f.Args()) > 0 {
    log.err("Usage error: unknown arguments: %v\n", f.Args())
    return subcommands.ExitUsageError
  }

  WebAppStart(this.bindAddr, this.prefixPath)

  log.success("Done\n")
  return subcommands.ExitSuccess
}
