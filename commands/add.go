package commands

import (
  "context"
  "flag"
  "fmt"
  "time"
  "github.com/google/subcommands"
)

type addFileCommand struct {
  bundleName        string
  ttl               time.Duration
  recipients        recipientTags
}

func AddFileCommand() (subcommands.Command) {
  return &addFileCommand{}
}

func (*addFileCommand) Name() string     { return "add" }
func (*addFileCommand) Synopsis() string { return "Add files to an existing bundle" }
func (*addFileCommand) Usage() string { return "add [options] <file> ... \n" }
func (this *addFileCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bundleName, "bundleName", "", "Bundle name")
}

func (this *addFileCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  filePaths := flagSet.Args()
  numFiles := len(filePaths)
  if numFiles < 1 {
    fmt.Printf("⚠️ Usage error: must provide at least one file\n")
    return subcommands.ExitUsageError
  }

  bundle, err := loadBundle(this.bundleName)
  if err != nil {
    fmt.Printf("❌ Failed to load bundle: %s\n", err)
    return subcommands.ExitFailure
  }
  logInfo("Loaded file bundle '%s'", bundle.name)

  for i, filePath := range filePaths {
    logInfo("Uploading file %d/%d: %s", i+1, numFiles, filePath)
    _, err := addFileToBundle(bundle, filePath)
    if err != nil {
      fmt.Printf("❌ Add file error '%s': %s\n", filePath, err)
      return subcommands.ExitFailure
    }
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
