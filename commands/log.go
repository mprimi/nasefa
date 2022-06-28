package commands

import (
  "fmt"
)

type logFunc func(string, ...interface{})(int, error)

var log struct {
  debug       logFunc
  info        logFunc
  warn        logFunc
  err         logFunc
  success     logFunc
}

func InitLogger() {
  debugEmoji := "🐛 "
  infoEmoji :=  "ℹ️ "
  warnEmoji :=  "⚠️ "
  errEmoji :=   "❌ "
  successEmoji :=   "✅ "

  if options.noEmojis {
    debugEmoji =     "[DBG] "
    infoEmoji =      "[ i ] "
    warnEmoji =      "[ ! ] "
    errEmoji =       "[!!!] "
    successEmoji =   "[ ✓ ] "
  }

  mutedFunc := func(_ string, _ ...interface{}) (int, error)  {
    return 0, nil
  }

  log.debug = mutedFunc
  log.info = mutedFunc
  log.warn = mutedFunc
  log.success = mutedFunc

  log.err = func(format string, a ...interface{}) (int, error)  {
    return fmt.Printf(errEmoji + format + "\n", a...)
  }

  if options.debug {
    log.debug = func(format string, a ...interface{}) (int, error)  {
      return fmt.Printf(debugEmoji + format + "\n", a...)
    }
  }

  if options.debug || ! options.quiet {
    log.info = func(format string, a ...interface{}) (int, error)  {
      return fmt.Printf(infoEmoji + format + "\n", a...)
    }
    log.warn = func(format string, a ...interface{}) (int, error)  {
      return fmt.Printf(warnEmoji + format + "\n", a...)
    }
    log.success = func(format string, a ...interface{}) (int, error)  {
      return fmt.Printf(successEmoji + format + "\n", a...)
    }
  }
}
