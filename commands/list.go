package commands

import (
  "context"
  "flag"
  "fmt"
  "github.com/google/subcommands"
  "github.com/c2h5oh/datasize"
)

type listCommand struct {
}

func ListCommand() (subcommands.Command) {
  return &listCommand{}
}

func (*listCommand) Name() string     { return "list" }
func (*listCommand) Synopsis() string { return "List available files" }
func (*listCommand) Usage() string { return "list [options] <file> ...\n" }
func (this *listCommand) SetFlags(f *flag.FlagSet) {
}

func (p *listCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  if len(f.Args()) > 0 {
    log.err("Usage error: unknown arguments: %v\n", f.Args())
    return subcommands.ExitUsageError
  }

  bundles, err := loadBundles()
  if err != nil {
    log.err("List failed: %s\n", err)
    return subcommands.ExitFailure
  }

  fmt.Printf("Found %d file bundles:\n", len(bundles))

  for _, bundle := range bundles {
    bundleSize := datasize.ByteSize(bundle.objStoreStatus.Size())
    bundleExpiration := bundle.objStoreStatus.TTL()
    expiration := "never"
    if bundleExpiration.Nanoseconds() > 0 {
      expiration = bundleExpiration.String()
    }
    fmt.Printf(" * %s (%d files, %s, expires: %s)\n",
      bundle.name,
      len(bundle.files),
      bundleSize.HumanReadable(),
      expiration,
    )
    for _, file := range bundle.files {
      fileSize := datasize.ByteSize(file.objInfo.Size)
      fmt.Printf("   - %s (%s)\t[%s]\n",
        file.fileName,
        fileSize.HumanReadable(),
        file.id,
      )
      bundleSize += fileSize
    }
  }

  log.success("Done\n")
  return subcommands.ExitSuccess
}
