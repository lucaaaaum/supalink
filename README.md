# supalink

Ever wanted to symlink multiple files using RegEx patterns? I sure do, all the time!

Well, guess what, now you can! *supalink* makes it easy and straightforward to do multiple symlinks in one go, providing some interesting flags to allow flexibility. It also goes fine with existing scripts, if you want to integrate it into a larger workflow.

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
supalink "/path/to/downloads/\[TorrentMaintainer\] Video/.*\.mkv" "/path/to/library/Video/Season \$STEP/Video S\$STEPE\$STEP_COUNT.mkv" --step 2 --step 2
```

> But this looks ridiculous! What on Earth is going on?

### Explanation

Yeah so RegEx is a funny thing. You can write it, but it's really, really hard to read it. What we're doing there basically tells *supalink* to:

1. Go to the directory `/path/to/downloads/[TorrentMaintainer] Video`
2. Search for all files that match the RegEx `.*\.mkv` (all MKV files, basically)
3. For each file found, create a symlink in `/path/to/library/Video/Season $STEP/Video S$STEPE$STEP_COUNT.mkv`, where:
   - `$STEP` is the current season number (starts counting from 1)
   - `$STEP_COUNT` is the total number of episodes processed so far

You'll understand more the more you use it. So, for testing purposes, you can always use

```bash
supalink --dry-run
```

which runs supalink but does not actually make any symlinks.

And, if you're unsure on how the links will actually end-up like, you can also just run

```bash
supalink --confirm
```

which will ask for confirmation before creating any symlinks.

### I still need more explanation

You can always check the available flags and their descriptions with

```bash
supalink --help
```

## Installation

Because I'm a [nix](https://nixos.org/) nerd, you can install *supalink* using the [nix flake](./flake.nix), just like below:

```nix
{
  inputs.supalink.url = "github:lucaaaaum/supalink";

  yourPackageSource* = [
    inputs.supalink.packages.${yourPackageSource.system}.supalink
  ];
}
```

*\* this being your **system.Packages** or **home.Packages** or whatever*

## Roadmap

- [ ] Use [Charm's](https://github.com/charmbracelet) toolset to make *supalink* prettier;
- [X] Actually implement Step functionality (it's on the example but it's not working lol);
- [ ] Improve user experience:
    - [ ] Make it so the user doesn't have to add quote marks on input;
    - [ ] Make it so the user doesn't have to escape RegEx characters that shouldn't be treated as RegEx (brackets on torrent maintainers names, dots, dashes, etc.);
- [ ] Release it on nixpkgs and other package managers;
- [ ] Create a snippet using Charm's [VHS](https://github.com/charmbracelet/vhs).

## Quick Q/A

***Q: why?***

**A:** torrent mantainers don't all follow the same naming conventions, making it hard to utilize tools such as [Jellyfin](https://jellyfin.org/), which do in fact rely on naming conventions. What I'm currently doing is separating my the torrenting downloads from my media library. So all files are actually on the downloads directory, while the media library simply points to those files.

***Q: couldn't you just write a bash script?***

**A:** yeah but I don't like writing bash. Sorry. Better than whatever powershell has going on though.

***Q: is this cross-platform?***

**A:** if you manage to install this on windows or mac, it should work. I only made the binary available for nix though.

***Q: are all the features implemented?***
**A:** nope, this is still WIP. Feel free to open issues or PRs if you want to contribute!
