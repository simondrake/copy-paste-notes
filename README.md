# Copy Paste Notes

copy-paste-notes is a CLI that is designed to make it possible to manage your notes from the command line. You can add, list, delete and copy notes directly to your system clipboard.

It is in very early development, so it comes with no assurances. It has only been tested on Ubuntu Linux. It _may_ work with your system, but it is in no way guaranteed.

# Setup/Installation

* Install - `go install github.com/simondrake/copy-paste-notes@latest`
* Create the database file - `touch ~/cpn.db`

# Platform Specific Details

copy-paste-notes relies on the `golang.design/x/clipboard` package, please refer to [their platform specific details](golang.design/x/clipboard) otherwise you may encounter errors.

## Linux

### Wayland

In addition to installing `libx11-dev` or `xorg-dev` or `libX11-devel`, you'll also need to install [wl-clipboard](https://github.com/bugaevc/wl-clipboard). The `x/clipboard` integration only seems to work with X11 (not Wayland).


# TODO

* [ ] Tests 🙈
* [ ] See if there's a way of making this work without the `os/exec` / `wl-clipboard` hack.
* [ ] Test on different platforms.
* [ ] Support releases/tags
