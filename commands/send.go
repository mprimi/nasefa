package commands

import (
  "context"
  "flag"
  "fmt"
  "time"
  "github.com/google/subcommands"
)

type sendCommand struct {
  bundleName        string
  ttl               time.Duration
}

func SendCommand() (subcommands.Command) {
  return &sendCommand{}
}

func (*sendCommand) Name() string     { return "send" }
func (*sendCommand) Synopsis() string { return "Send a bundle of files" }
func (*sendCommand) Usage() string { return "send [options] <file> ... \n" }
func (this *sendCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bundleName, "bundleName", "", "Unique ID for this file bundle, used for download (randomly generated if not provided)")
  f.DurationVar(&this.ttl, "expire", 0, "Automatically delete this file bundle after a certain amount of time")
}

func (this *sendCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  filePaths := flagSet.Args()
  numFiles := len(filePaths)
  if numFiles < 1 {
    fmt.Printf("⚠️ Usage error: must provide at least one file\n")
    return subcommands.ExitUsageError
  }

  bundle, err := newBundle(this.bundleName, this.ttl)
  if err != nil {
    fmt.Printf("❌ Failed to create bundle: %s\n", err)
    return subcommands.ExitFailure
  }
  logInfo("Created file bundle '%s'", bundle.name)

  for i, filePath := range filePaths {
    logInfo("Uploading file %d/%d: %s", i+1, numFiles, filePath)
    bundleFile, err := addFileToBundle(bundle, filePath, "")
    if err != nil {
      fmt.Printf("❌ Send error '%s': %s\n", filePath, err)
      return subcommands.ExitFailure
    }
    logInfo("Added file '%s' => '%s'", bundleFile.fileName, bundleFile.id)
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
