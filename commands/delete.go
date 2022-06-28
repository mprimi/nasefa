package commands

import (
  "context"
  "flag"
  "github.com/google/subcommands"
)

type deleteBundleCommand struct {
  bundleName        string
}

func DeleteBundleCommand() (subcommands.Command) {
  return &deleteBundleCommand{}
}

func (*deleteBundleCommand) Name() string     { return "delete" }
func (*deleteBundleCommand) Synopsis() string { return "Deletes a files bundle" }
func (*deleteBundleCommand) Usage() string { return "delete [options] \n" }
func (this *deleteBundleCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bundleName, "bundleName", "", "Name of the bundle to delete")
}

func (this *deleteBundleCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if this.bundleName == "" {
    log.err("Usage error: missing bundle name\n")
    return subcommands.ExitUsageError
  }

  err := deleteBundle(this.bundleName)
  if err != nil {
    log.err("Failed to delete bundle: %s\n", err)
    return subcommands.ExitFailure
  }
  log.info("Deleted file bundle '%s'", this.bundleName)

  log.success("Done\n")
  return subcommands.ExitSuccess
}
