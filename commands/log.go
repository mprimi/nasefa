package commands

import (
  "fmt"
)

func logDebug(format string, a ...interface{}) (int, error)  {
  return fmt.Printf(" 🐛 " + format + "\n", a...)
}

func logInfo(format string, a ...interface{}) (int, error)  {
  return fmt.Printf(" ℹ️  " + format + "\n", a...)
}

func logWarn(format string, a ...interface{}) (int, error)  {
  return fmt.Printf(" ⚠️ " + format + "\n", a...)
}
