# Lost Grace Client Application

## What is this?

This is the companion application for the [lostgrace](https://skulle.xyz/lost_grace.php) service. A service that is
designed to make syncing your Seamless Coop save easy.

You can also install the [Seamless Coop mod](https://github.com/LukeYui/EldenRingSeamlessCoopRelease) automagically with
this as well.

## Why did you make it?

Plugging a thumb drive into something is not as straight forward as it should be sometimes (Steam Deck, ROG Ally, etc.).
Also wanted more reps in [GO](https://go.dev) with GUI programming. Even though my day to day is in GO, I wanted to
write something fun.

I am also addicted to making tools for some reason. This one actually seems useful.

## Supported OS

- Windows 11
- Ubuntu 24.04
- I will add Steam Deck once I test it

## Running and building

### Run

Get the latest [release](https://github.com/rskulles/lostgrace/releases/latest).

You can take a look at [lostgrace](https://skulle.xyz/lost_grace.php) for some outdated instructions that still work.

But the gist is:

- Register using the link above.
- Download client to your machines.
- Put configuration info into client. Make sure to save.
- Sync your save file, or install the latest Seamless Coop mod.
- You will also be able to do everything from command line, so you can use my tool in your tool, so you can tool while
  you are tooling.

### Build

- Have `go` installed. Currently on 1.22.
- Clone repo.
- This uses Gio UI. You need to install its [prerequisites](https://gioui.org/doc/install). Follow the guide for the OS
  you are building on.
- Run `go build .` and voilà. Too easy. Gio UI has tools for deployment, but I use `go build .` and `go run .` while
  developing.

## Is this affiliated with FromSoftware or Seamless Coop?

Nope and nope. 

