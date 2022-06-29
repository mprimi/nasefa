package commands

import (
  "context"
  "flag"
  "github.com/google/subcommands"
)

const (
  kDefaultBindAddr = ":8080"
  kDefaultPrefixPath = "/"
  kDefaultAllowList = false
)

type webCommand struct {
  bindAddr    string
  certPath    string
  certKeyPath string
  prefixPath  string
  allowList   bool
}

func WebCommand() (subcommands.Command) {
  return &webCommand{}
}

func (*webCommand) Name() string     { return "web" }
func (*webCommand) Synopsis() string { return "Starts web application" }
func (*webCommand) Usage() string { return "web [options]\n" }
func (this *webCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bindAddr, "bindAddr", kDefaultBindAddr, "Address to bind")
  f.StringVar(&this.certPath, "cert", "", "Path to certificate file (for HTTPS)")
  f.StringVar(&this.certKeyPath, "certKey", "", "Path to certificate key file (for HTTPS)")
  f.StringVar(&this.prefixPath, "prefix", kDefaultPrefixPath, "HTTP path prefix")
  f.BoolVar(&this.allowList, "allowBundlesListing", kDefaultAllowList, "Allow listing of bundles via web (Danger!)")
}

func (this *webCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if len(f.Args()) > 0 {
    log.err("Usage error: unknown arguments: %v\n", f.Args())
    return subcommands.ExitUsageError
  }

  WebAppStart(this.bindAddr, this.certPath, this.certKeyPath, this.prefixPath, this.allowList)

  log.success("Done\n")
  return subcommands.ExitSuccess
}
