package commands

import (
  "context"
  "flag"
  "time"
  "github.com/google/subcommands"
)

type createBundleCommand struct {
  bundleName        string
  ttl               time.Duration
  recipients        recipientTags
}

func CreateBundleCommand() (subcommands.Command) {
  return &createBundleCommand{}
}

func (*createBundleCommand) Name() string     { return "create" }
func (*createBundleCommand) Synopsis() string { return "Create a new empty bundle (useful for accepting file uploads via web)" }
func (*createBundleCommand) Usage() string { return "create_bundle [options] \n" }
func (this *createBundleCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bundleName, "bundleName", "", "Unique ID for this (empty) bundle of files (randomly generated if not provided)")
  f.DurationVar(&this.ttl, "expire", 0, "Automatically delete this bundle after a specified amount of time")
}

func (this *createBundleCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if this.bundleName == "" {
    log.err("Usage error: missing bundle name\n")
    return subcommands.ExitUsageError
  }

  bundle, err := newBundle(this.bundleName, this.ttl)
  if err != nil {
    log.err("Failed to create bundle: %s\n", err)
    return subcommands.ExitFailure
  }
  log.info("Created empty file bundle '%s'", bundle.name)

  log.success("Done\n")
  return subcommands.ExitSuccess
}
