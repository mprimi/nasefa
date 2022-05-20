package main

import (
  "log"
  "os"
  "nasefa/helpers"
)

func send(sourceFile string)  {
  err := helpers.UploadFile(sourceFile)
  if err != nil {
    log.Printf("File upload failed: %s", err)
    os.Exit(1)
  }
}

func recv(fileId, destinationDirectory string)  {
  err := helpers.DownloadFiles(fileId, destinationDirectory)
  if err != nil {
    log.Printf("Files download failed: %s", err)
    os.Exit(1)
  }
}

func main()  {
  if len(os.Args) < 3 {
    log.Printf("Usage:")
    log.Printf(" nasefa send <file>")
    log.Printf(" nasefa recv <file-id> <directory>")
    os.Exit(1)
  }

  if os.Args[1] == "send" {
    send(os.Args[2])
  } else if os.Args[1] == "recv" {
    recv(os.Args[2], os.Args[3])
  } else {
    log.Printf("Invalid options")
    os.Exit(1)
  }
}
