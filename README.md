# supalink

Ever wanted to symlink multiple files using RegEx patterns? I sure do. All the time!

Well guess what, now you can! *supalink* makes it easy and straightforward to do multiple symlinks in one go, providing some interesting flags to allow flexibility. It goes really 
well with existing scripts as well, if you want to integrate into a larger workflow.

## Usage

Say we have this file structure:

```
.
├── [TorrentMaintainer] Video - Episode 1.mkv
├── [TorrentMaintainer] Video - Episode 2.mkv
├── [TorrentMaintainer] Video - Episode 3.mkv
└── [TorrentMaintainer] Video - Episode 4.mkv
```

But you actually want this:

```
.
├── Season 1
│   ├── Video S01E01.mkv
│   └── Video S01E02.mkv
└── Season 2
    ├── Video S02E01.mkv
    └── Video S02E02.mkv
```

You can do this with *supalink* like so:

```bash
supalink "/path/to/downloads/[TorrentMaintainer] Video/.*Episode \d+\.mkv" "/path/to/library/Video/Season \$STEP/Video S\$STEPE\$STEP_COUNT.mkv -step 2 2"
```

> But this looks ridiculous! What on Earth is going on?

Yeah you so RegEx is a funny thing. You can write it, but it's really, really hard to read it. What we're doing there basically tells supalink to:

### Explanation

1. Go to the directory `/path/to/downloads/[TorrentMaintainer] Video`
2. Search for all files that match the RegEx `.*Episode \d+\.mkv` (which is all files that end with "Episode X.mkv")
3. For each file found, create a symlink in `/path/to/library/Video/Season $STEP/Video S$STEPE$STEP_COUNT.mkv`, where:
   - `$STEP` is the current season number (starting from 1)
   - `$STEP_COUNT` is the total number of episodes processed so far

### I still need more explanation

You can always check the available flags and their descriptions with

```bash
supalink --help
```

or even

```bash
bashsupalink -h
```

## Installation

Because I'm a [nix](https://nixos.org/) nerd, you can install *supalink* using the [nix flake](./flake.nix), just like below:

```nix
{
  inputs.supalink.url = "github:lucaaaaum/supalink";

  yourPackageSource* = [
    inputs.supalink.packages.${pkgs.system}.supalink
  ];
}
```

*\* this being your **system.Packages** or **home.Packages** or whatever*

## Quick Q/A

***Q: why?***

**A:** torrent mantainers don't all follow the same naming conventions, making it hard to utilize tools such as [Jellyfin](https://jellyfin.org/), which do in fact rely on naming conventions. What I'm currently doing is separating my the torrenting downloads from my media library. So all files are actually on the downloads directory, while the media library simply points to those files.

***Q: couldn't you just write a bash script?***

**A:** yeah but I don't like writing bash. Sorry. Better than whatever powershell has going on though.

***Q: is this cross-platform?***

**A:** if you manage to install this on windows or mac, it should work. I only made the binary available for nix though.

***Q: are all the features implemented?***
**A:** nope, this is still WIP. Feel free to open issues or PRs if you want to contribute!
