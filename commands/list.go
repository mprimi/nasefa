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

  bundles, err := loadBundles()
  if err != nil {
    fmt.Printf("❌ List failed: %s\n", err)
    return subcommands.ExitFailure
  }

  fmt.Printf("Found %d file bundles:\n", len(bundles))

  for _, bundle := range bundles {
    bundleSize := datasize.ByteSize(bundle.objStoreStatus.Size())
    fmt.Printf(" * %s (%d files, %s)\n",
      bundle.name,
      len(bundle.files),
      bundleSize.HumanReadable(),
    )
    for _, file := range bundle.files {
      fileSize := datasize.ByteSize(file.objInfo.Size)
      fmt.Printf("   - %s (%s) [%s]\n",
        file.fileName,
        fileSize.HumanReadable(),
        file.id,
      )
      bundleSize += fileSize
    }
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
