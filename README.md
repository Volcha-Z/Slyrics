<div align="center">

# Slyrics

Synchronized lyrics in your terminal.

![preview](./preview.png)

</div>

## Features

- Spotify, MPD, Mopidy, MPRIS, or any browser (including YouTube Music) via the [WebNowPlaying](https://wnp.keifufu.dev/extension/getting-started) extension
- Three lyrics sources ([lrclib](https://lrclib.net), NetEase, QQ Music) queried in parallel, so a track missing from one usually turns up on another

## Build from source

Requires Go 1.26+.

```sh
git clone https://github.com/Volcha-Z/Slyrics.git && cd Slyrics && ./install.sh
```

Installs `slyrics` to `~/.local/bin` (override with `BINDIR=... ./install.sh`).

## Quick start

A config file is created on first run, at `~/.config/slyrics/config.yaml`.

### YouTube Music (or any other browser player)

1. Install the [WebNowPlaying extension](https://wnp.keifufu.dev/extension/getting-started) in your browser.
2.  Put the player as browser in the config file.
3.  Play something, then run `slyrics`.

### Spotify, MPD, Mopidy, MPRIS, or local `.lrc` files

Set `player` in the config to `spotify`, `mpd`, `mopidy`, `mpris`, or leave it unset and fill in `local.folder` to scan a folder of `.lrc` files instead. Each player's own settings live under its matching key in the config — see the comments in the generated file.

Spotify needs a one-time login: create an app at [developer.spotify.com](https://developer.spotify.com/dashboard) with redirect URI `http://127.0.0.1:8888/callback`, then run `slyrics login` and follow the prompts.

### Changing the color

`slyrics` (the wrapper installed by `install.sh`) accepts a color flag: `slyrics -blue`, `-red`, `-orange`, `-yellow`, `-green`, `-cyan`, `-purple`, `-pink`, `-white`, or any hex code directly, e.g. `slyrics -#ff00ff`.

## Style

```yaml
style:
  blockFont: true      # big ASCII-block letters instead of normal text
  effect3d: false      # adds a 3d effect for the ASCII-block
  onlyCurrent: true     # show only the current line, centered
  uppercase: true
  hAlignment: center
  current:
    bold: true
    foreground: "#3b82f6"   # HEX or ANSI 0-255
```

Style flags also override the config on the command line, e.g. `slyrics --current "bold,#3b82f6"`. Run `slyrics --help` for the full list.

> **Note:** `blockFont` needs enough terminal width to draw a whole word at once, a window that's too narrow for a given line falls back to plain text Wider is better for block.
## License

MIT — see [LICENSE](./LICENSE).
