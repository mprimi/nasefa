package main

import (
  "context"
  "os"
  "flag"
  "github.com/google/subcommands"
  "nasefa/commands"
)


const (
  kSendReceiveGroup = "sending and receiving bundles of files"
  kFileAndBundleOpsGroup = "file and bundle management"
  kHelpGroup = "help"
  kWebGroup = "web"
)

func main()  {
  defer commands.ClientCleanup()
  commands.RegisterTopLevelFlags()

  commandsMap := map[string][]subcommands.Command{
    kSendReceiveGroup: []subcommands.Command{
      commands.SendCommand(),
      commands.ReceiveCommand(),
      commands.AutoreceiveCommand(),
    },
    kFileAndBundleOpsGroup: []subcommands.Command{
      commands.ListCommand(),
      commands.CreateBundleCommand(),
      commands.DeleteBundleCommand(),
      commands.AddFileCommand(),
    },
    kWebGroup: []subcommands.Command{
      commands.WebCommand(),
    },
    kHelpGroup: []subcommands.Command{
      subcommands.HelpCommand(),
      subcommands.FlagsCommand(),
      subcommands.CommandsCommand(),
    },
  }

  for groupName, commands := range commandsMap {
    for _, command := range commands {
      subcommands.Register(command, groupName)
    }
  }

  flag.Parse()
  ctx := context.Background()
  os.Exit(int(subcommands.Execute(ctx)))
}
