package commands

import (
  "context"
  "flag"
  "fmt"
  "github.com/google/subcommands"
)

type receiveCommand struct {
}

func ReceiveCommand() (subcommands.Command) {
  return &receiveCommand{}
}

func (*receiveCommand) Name() string     { return "receive" }
func (*receiveCommand) Synopsis() string { return "Receive one or more file bundles" }
func (*receiveCommand) Usage() string { return "receive [options] <destination_directory> <bundle> ...\n" }
func (p *receiveCommand) SetFlags(f *flag.FlagSet) {
}

func (p *receiveCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  numBundles := len(flagSet.Args()) - 1

  if numBundles < 1 {
    fmt.Printf("⚠️ Usage error: destination directory and at least one bundle id required\n")
    return subcommands.ExitUsageError
  }

  destinationDirectory := flagSet.Args()[0]
  bundleNames := flagSet.Args()[1:]

  totalFiles := 0
  for _, bundleName := range bundleNames {
    bundle, err := downloadBundle(destinationDirectory, bundleName)
    if err != nil {
      fmt.Printf("❌ Error receiving bundle '%s': %s\n", bundleName, err)
      return subcommands.ExitFailure
    }
    totalFiles += len(bundle.files)
  }

  fmt.Printf("✅ Received %d files in %d bundles\n", totalFiles, numBundles)
  return subcommands.ExitSuccess
}
