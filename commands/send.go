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
  fileId        string
}

func SendCommand() (subcommands.Command) {
  return &sendCommand{}
}

func (*sendCommand) Name() string     { return "send" }
func (*sendCommand) Synopsis() string { return "Send a file" }
func (*sendCommand) Usage() string { return "send [options] <file>\n" }
func (this *sendCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bucketName, "bucket", defaultBucketName, "Name of the bucket where file is stored")
  f.StringVar(&this.fileId, "fileId", "", "ID of the file, used for download (randomly generated if not provided)")
}

func (p *sendCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if len(flagSet.Args()) != 1 {
    fmt.Printf("⚠️ Usage error: must provide a single file path\n")
    return subcommands.ExitUsageError
  }

  filePath := flagSet.Args()[0]

  objStore, err := getObjStore(p.bucketName)
  if err != nil {
    fmt.Printf("❌ %s\n", err)
    return subcommands.ExitFailure
  }

  file, err := os.Open(filePath)
  if err != nil {
    fmt.Printf("❌ Send error, failed to open '%s': %s\n", filePath, err)
    return subcommands.ExitFailure
  }
  defer file.Close()

  fileId := p.fileId
  if fileId == "" {
    fileId = uuid.NewString()
  }

  fmt.Printf("📤 Sending file %s => %s\n", filePath, fileId)

  filename := path.Base(filePath)
  objMeta := nats.ObjectMeta{
    Name: fileId,
    Description: fmt.Sprintf("Nasefa file object (%s)", filename),
  }
  setFilename(&objMeta, filename)
  objInfo, err := objStore.Put(&objMeta, file)

  logDebug("Uploaded: %s/%s (%s) => %v", p.bucketName, objMeta.Name, objInfo.NUID, objInfo)

  if err != nil {
    fmt.Printf("❌ Send error: %s\n", err)
    return subcommands.ExitFailure
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
