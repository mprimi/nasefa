package commands

import (
  "context"
  "flag"
  "fmt"
  "nasefa/helpers"
  "github.com/google/subcommands"
)

type sendCommand struct {
  bucketName    string
}

func SendCommand() (subcommands.Command) {
  return &sendCommand{}
}

func (*sendCommand) Name() string     { return "send" }
func (*sendCommand) Synopsis() string { return "Send a file" }
func (*sendCommand) Usage() string { return "send [options] <file> ...\n" }
func (this *sendCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bucketName, "bucket", defaultBucketName, "Name of the bucket where file is stored")
}

func (p *sendCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  numFiles := len(f.Args())

  if numFiles < 1 {
    fmt.Printf("⚠️ Usage error: no files provided\n")
    return subcommands.ExitUsageError
  }

  for i, filePath := range f.Args() {
    fmt.Printf("📤 Sending file %d/%d: %s\n", i+1, numFiles, filePath)
    err := helpers.UploadFile(filePath)
    if err != nil {
      fmt.Printf("❌ File upload failed: %s\n", err)
      return subcommands.ExitFailure
    }
  }

  return subcommands.ExitSuccess
}
