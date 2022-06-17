package commands

import (
  "context"
  "flag"
  "fmt"
  "os"
  "path"
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

    file, err := os.Open(filePath)
    if err != nil {
      fmt.Printf("❌ Failed to open %s: %s\n", filePath, err)
      return subcommands.ExitFailure
    }
    defer file.Close()

    fileName := path.Base(filePath)

    bundleFile, err := addFileToBundle(bundle, file, fileName, "")
    file.Close()
    if err != nil {
      fmt.Printf("❌ Send error '%s': %s\n", filePath, err)
      return subcommands.ExitFailure
    }
    logInfo("Added file '%s' => '%s'", bundleFile.fileName, bundleFile.id)
  }

  if len(this.recipients) > 0 {
    err := notifyRecipients(bundle, this.recipients...)
    if err != nil {
      fmt.Printf("❌ Failed to notify bundle recipients: %s\n", err)
      return subcommands.ExitFailure
    }
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
