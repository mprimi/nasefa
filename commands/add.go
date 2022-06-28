package commands

import (
  "context"
  "flag"
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

  if this.bundleName == "" {
    log.err("Usage error: missing bundle name\n")
    return subcommands.ExitUsageError
  }

  filePaths := flagSet.Args()
  numFiles := len(filePaths)
  if numFiles < 1 {
    log.err("Usage error: must provide at least one file\n")
    return subcommands.ExitUsageError
  }

  bundle, err := loadBundle(this.bundleName)
  if err != nil {
    log.err("Failed to load bundle: %s\n", err)
    return subcommands.ExitFailure
  }
  log.info("Loaded file bundle '%s'", bundle.name)

  for i, filePath := range filePaths {
    log.info("Uploading file %d/%d: %s", i+1, numFiles, filePath)
    _, err := addFileToBundle(bundle, filePath)
    if err != nil {
      log.err("Add file error '%s': %s\n", filePath, err)
      return subcommands.ExitFailure
    }
  }

  log.success("Done\n")
  return subcommands.ExitSuccess
}
