package commands

import (
  "context"
  "flag"
  "strings"
  "github.com/google/subcommands"
)

type autoreceiveCommand struct {
}

func AutoreceiveCommand() (subcommands.Command) {
  return &autoreceiveCommand{}
}

func (*autoreceiveCommand) Name() string     { return "auto-receive" }
func (*autoreceiveCommand) Synopsis() string { return "Listen and receive file bundles based on recipient tags" }
func (*autoreceiveCommand) Usage() string { return "auto-receive [options] <destination_directory> <receiver_tag> ...\n" }
func (p *autoreceiveCommand) SetFlags(f *flag.FlagSet) {
}

func (p *autoreceiveCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  numRecipient := len(flagSet.Args()) - 1

  if numRecipient < 1 {
    log.err("Usage error: destination directory and at least one recipient tag required\n")
    return subcommands.ExitUsageError
  }

  destinationDirectory := flagSet.Args()[0]
  recipientTagNames := flagSet.Args()[1:]

  bundleNotificationsCh, err := watchBundles(recipientTagNames...)
  if err != nil {
    log.err("Error subscribing to bundle notifications: %s\n", err)
    return subcommands.ExitFailure
  }

  log.success("Watching for file bundles tagged: %s\n", strings.Join(recipientTagNames, "|"))

  for bundleName := range bundleNotificationsCh {
    bundle, err := downloadBundle(destinationDirectory, bundleName)
    if err != nil {
      log.err("Error receiving bundle '%s': %s\n", bundleName, err)
      return subcommands.ExitFailure
    }
    log.success("Received bundle '%s' (%d files)\n", bundleName, len(bundle.files))
  }

  return subcommands.ExitSuccess
}
