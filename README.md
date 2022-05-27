# Nasefa
## Send and receive files using NATS

Nasefa is a utility to send and receive files between computers.

### Why..?

1. It's 2022, and somehow it's still hard to move files between my computers (running Linux, Windows and macOS).
2. I've been itching to try building something on NATS

*Nasefa* is a girl name, and it also vaguely sounds like **NA**TS **SE**nd **F**iles.

### Use cases

I decided to build something that would help me with the the following scenarios:

 - Send files to a computer on the home network, regardless of platform/OS
 - Send files to a computer even if it currently is offline/asleep
 - Send files to one or more remote computers and have them automatically download it at a given location
 - Send a file to a friend without using Dropbox/GDrive/...

## Usage:

`$ nasefa send -fileId taxes tax-return.pdf`

`$ nasefa receive taxes ~/Desktop`
