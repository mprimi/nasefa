package main

import (
  "context"
  "os"
  "flag"
  "github.com/google/subcommands"
  "nasefa/commands"
)


func main()  {
  commands.RegisterTopLevelFlags()
  subcommands.Register(commands.SendCommand(), "")
  subcommands.Register(commands.ReceiveCommand(), "")
  subcommands.Register(commands.ListCommand(), "")
  subcommands.Register(subcommands.HelpCommand(), "help")
  subcommands.Register(subcommands.FlagsCommand(), "help")
  subcommands.Register(subcommands.CommandsCommand(), "help")

  flag.Parse()
  ctx := context.Background()
  os.Exit(int(subcommands.Execute(ctx)))
}
