package commands

import (
  "context"
  "flag"
  "fmt"
  "github.com/google/subcommands"
  "github.com/c2h5oh/datasize"
)

type listCommand struct {
  bucketName    string
}

func ListCommand() (subcommands.Command) {
  return &listCommand{}
}

func (*listCommand) Name() string     { return "list" }
func (*listCommand) Synopsis() string { return "List available files" }
func (*listCommand) Usage() string { return "list [options] <file> ...\n" }
func (this *listCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&this.bucketName, "bucket", defaultBucketName, "Name of the bucket to list")
}

func (p *listCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if len(f.Args()) > 0 {
    fmt.Printf("⚠️ Usage error: unknown arguments: %v\n", f.Args())
    return subcommands.ExitUsageError
  }

  objStore, err := getObjStore(p.bucketName)
  if err != nil {
    fmt.Printf("❌ %s\n", err)
    return subcommands.ExitFailure
  }

  objsInfo, err := objStore.List()
  if err != nil {
    fmt.Printf("❌ List error: %s\n", err)
    return subcommands.ExitFailure
  }

  fmt.Printf("Listing bucket: %s\n", p.bucketName)

  skipped := []string{}
  filesCount := 0
  filesSize := datasize.ByteSize(0)
  for _, objInfo := range objsInfo {
    filename := getFilename(objInfo)
    if filename == "" {
      skipped = append(skipped, objInfo.Name)
    } else {
      fileId := objInfo.Name
      fileSize := datasize.ByteSize(objInfo.Size)
      fmt.Printf(" - %s [%s] [%s]\n", fileId, filename, fileSize.HumanReadable())
      filesCount += 1
      filesSize += fileSize
    }
  }
  fmt.Printf("Total: %d files (%s)\n", filesCount, filesSize.HumanReadable())

  if len(skipped) > 0 {
    fmt.Printf("⚠️ Ignored %d objects that are not Nasefa files (%v)\n", len(skipped), skipped)
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
