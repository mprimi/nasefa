package commands

import (
  "context"
  "flag"
  "strings"
  "time"
  "github.com/google/subcommands"
)

type sendCommand struct {
  bundleName        string
  ttl               time.Duration
  recipients        recipientTags
}

func SendCommand() (subcommands.Command) {
  return &sendCommand{}
}

type recipientTags []string
func (this *recipientTags) Set(tag string) error {
  *this = append(*this, tag)
  return nil
}
func (this *recipientTags) String() string {
  return strings.Join([]string(*this), ",")
}

func (*sendCommand) Name() string     { return "send" }
func (*sendCommand) Synopsis() string { return "Send a bundle of files" }
func (*sendCommand) Usage() string { return "send [options] <file> ... \n" }
func (this *sendCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bundleName, "bundleName", "", "Unique ID for this file bundle, used for download (randomly generated if not provided)")
  f.DurationVar(&this.ttl, "expire", 0, "Automatically delete this file bundle after a certain amount of time")
  f.Var(&this.recipients, "to", "Tag name for recipients in auto-download (can be repeated)") //TODO: Crappy description, review later
}

func (this *sendCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  filePaths := flagSet.Args()
  numFiles := len(filePaths)
  if numFiles < 1 {
    log.err("Usage error: must provide at least one file\n")
    return subcommands.ExitUsageError
  }

  bundle, err := newBundle(this.bundleName, this.ttl)
  if err != nil {
    log.err("Failed to create bundle: %s\n", err)
    return subcommands.ExitFailure
  }
  log.info("Created file bundle '%s'", bundle.name)

  for i, filePath := range filePaths {
    log.info("Uploading file %d/%d: %s", i+1, numFiles, filePath)
    _, err := addFileToBundle(bundle, filePath)
    if err != nil {
      log.err("Add file error '%s': %s\n", filePath, err)
      return subcommands.ExitFailure
    }
  }

  if len(this.recipients) > 0 {
    err := notifyRecipients(bundle, this.recipients...)
    if err != nil {
      log.err("Failed to notify bundle recipients: %s\n", err)
      return subcommands.ExitFailure
    }
  }

  log.success("Done\n")
  return subcommands.ExitSuccess
}
