package commands

import (
  "context"
  "flag"
  "fmt"
  "os"
  "path"
  "github.com/nats-io/nats.go"
  "github.com/google/subcommands"
  "github.com/google/uuid"
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

  objStore, err := getObjStore(p.bucketName)
  if err != nil {
    fmt.Printf("❌ %s\n", err)
    return subcommands.ExitFailure
  }

  for i, filePath := range f.Args() {
    fmt.Printf("📤 Sending file %d/%d: %s\n", i+1, numFiles, filePath)

    // _, err := objStore.PutFile(filePath)

    f, err := os.Open(filePath)
    if err != nil {
      fmt.Printf("❌ Send error, failed to open '%s': %s\n", filePath, err)
      return subcommands.ExitFailure
    }
    defer f.Close()

    fileId := uuid.NewString()
    filename := path.Base(filePath)
    objMeta := nats.ObjectMeta{
      Name: fileId,
      Description: fmt.Sprintf("Nasefa file object (%s)", filename),
    }
    setFilename(&objMeta, filename)
    objInfo, err := objStore.Put(&objMeta, f)

    logDebug("Uploaded: %s/%s (%s) => %v", p.bucketName, objMeta.Name, objInfo.NUID, objInfo)

    if err != nil {
      fmt.Printf("❌ Send error: %s\n", err)
      return subcommands.ExitFailure
    }
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
