# Nasefa

Nasefa is a utility to send and receive files on a local network, or over the internet.

![ThePlague]

It's 2022, computers can do incredible things, and yet it is still non-trivial to send/receive a file to your friends (without also sharing it with `$MEGACORP`).

Nasefa is primarily a fun & learning project (Golang, NATS).
But it turned into an utility that solves some real use-cases.

Here's how I use Nasefa:

## Use case 1: Send myself files (between computers on LAN)

Between me and family, we have computers running Linux, macOS, iOS, Android, Windows, and more.

Rather than email/Dropbox/Google Drive/USB stick/... I can now send myself files using:

```
$ nasefa send -bundleName tax_documents tax-return.pdf tax-summary.xls receipts.zip
```

Later on, on a different computer, I download the files with:

```
$ nasefa receive ~/Desktop tax_documents
```

This is convenient (for me), because the receiving computer does not have to be on/awake.

If I forget the bundle name (`tax_documents` in the example), I can look it up with `nasefa list`

## Use case 2: Deploy simple file changes

Between work and hobby projects, I often find myself using `scp` to copy files around. It is kind of a pain sometimes.

Start Nasefa as background process in auto-receive mode on selected hosts, example:

```
$ nasefa auto-receive ~/bin app-foo cluster-bar host-baz
```

This will automatically download any file bundle tagged `app-foo` or `cluster-bar` or `host-baz`.

To distribute an updated version of a binary to a selected group of servers, issue:

```
$ nasefa send -bundleName app-foo-binaries -to cluster-bar app-foo.exe
```

Any hosts tagged `cluster-bar` will auto-download the latest version of `app-foo.exe` into `~/bin`.

## Use case 3: Receive files from friends

Exchanging files with friends is a particular pain if you want to avoid the all-seeing-eye of `$MEGACORP`, or the file size is just too large.

Nasefa has a simple built-in web interface that anyone can understand and use.

On your (internet-exposed) host, you can run:

```
$ nasefa web -bindAddr :8080
```

If a friend wants to send you a collections of pictures and videos, you can create an empty bundle with:

```
$ nasefa create -bundleName holiday_pictures
```

Then send a link to your friend: https://nasefa.example.com:8080/upload/holiday_pictures/.
They'll be able to upload one or more files through a simple web form.

## Use case 4: Send files to friends

If you are the one sharing, you can upload via CLI or web, then share a link such as: https://nasefa.example.com:8080/holiday_pictures/.

They'll be able to download on any device using any web browser.

#### Auto-expiring file bundles

Especially for files exchanged via web, it can be useful to:
 - Automatically delete shared files after some time
 - Automatically block upload to a bundle after some time

This is easily done, for example:

```
$ nasefa send -expire 3d -bundleName hiking-pics hiking-pictures.zip
```

The `hiking-pics` bundle will automatically disappear 3 days after upload.

# Deployment

Nasefa is built on (open-source) [NATS](https://nats.io/).
It inherits its flexibility and adaptability to the most diverse scenarios.
It supports various kinds of authentication & authorization and it can be exposed directly to the internet.

Running a NATS server for Nasefa is as simple as:

```
$ nats-server -js
```

It runs on most modern OSs and images are available for, Docker, K8S.
If you want to enable TLS, isolate accounts or applications, set resource quotas, or set up federation, you may need to [configure](https://docs.nats.io/running-a-nats-service/configuration) a few more things. It's all pretty easy and intuitive.

Here's a few ways how I personally self-host:

### Home network

NATS server can run on nimble devices like Raspberry Pi.
I personally run on an always-on old computer with a static local IP.

In this settings, I don't bother locking down NATS.
However I do still create a separate [account](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/accounts) for Nasefa. This allows me to run multiple "apps" on the same server in complete isolation.

### VPN

To have access to the server while on the go, I make service reachable via VPN. This is trivial when using something like [Tailscale](https://tailscale.com/).

### Internet-facing (web only)

Sending and receiving files from friends requires the web interface to be reachable, but not NATS.

Personally, I run NATS server bound to localhost, so it's port is closed to the world. On the same host, i run `nasefa web`, and this is reachable to anyone without authentication.

I don't advertise the host:port, and even if a malicious actor discovered the web UI, they'd still need to guess the bundle names in order to upload or download something.

I realize this is security by obscurity, but this is good enough for me, given the use case. If you wanted to block access to the web interface, you could put it behind a proxy with an authentication method of your choice.

### Internet-facing (full access)

If your friends (or collaborators) are comfortable with the command-line, you may run NATS exposed to the internet.

A good starting point is setting up a server-side TLS certificate, so clients can verify who they are connecting to.

To authenticate and authorize clients, [you have options](https://docs.nats.io/running-a-nats-service/configuration/securing_nats): tokens, passwords, JWT, NKEY, Mutual TLS, and more.

### Federation

Discussing multi-server deployment is way beyond the scope of this tiny utility, but it's absolutely possible.

In my case, the internet facing server and the home network server are configured so that files uploaded via web are automatically mirrored to my home server, but not the other way around.
This all works by just configuring NATS, without any change to Nasefa.

If you are feeling adventurous, you can even run NATS server in a WASM-capable browser. But that's a story for a different time.

---

# F.A.Q.

### Why didn't you use ___ instead?

I'm well aware of the existence of FTP, NAS, rsync, IPFS, S3, Wormhole, Firefox Send, Keybase, Freenet, NextCloud, Dropbox, Storj, git, NFS, etc, etc, etc.

Nasefa is primarily a fun & learning project. I'm not trying to create a product, and I'm not claiming this is the best solution.

p.s. if I wanted the same utility, I would have probably gone with FTP.
Simple to self-host, supported by most browsers and does all the things I need it to do (minus self-expiring files, which i could have implemented client-side).

### Nasefa is close to what I want, but not quite right, do you know anything else similar?

The answer to the previous question contains a few pointers.

[Awesome Self-Hosted](https://github.com/awesome-selfhosted/awesome-selfhosted) has numerous great options.

Of course, nothing matches the reliability of [S4](http://www.supersimplestorageservice.com/) and nothing is cheaper than [YouTubeDrive](https://github.com/dzhang314/YouTubeDrive).


### What does Nasefa stand for?

*Nasefa* is a female first name I always kinda liked. And it kinda sounds like **NA**TS **SE**nd **F**iles.

---

# Thoughts, comments, feedback, critiques, ...

I'd love to hear anything you have to say about Nasefa. Hit me up!

(Contact information is in my GitHub profile).



[ThePlague]: /theplague.png "It's not that easy, Mr ThePlague!"
