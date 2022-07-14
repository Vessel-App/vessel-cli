## TODO
- [x] Global debug flag (env var)
    - Print output from mutagen and other exec.Cmd calls
- [x] SSH config edits ~/.ssh/config to add an alias?!
    - Prints out helpful message "Add this yourself, or press 'y' for us to do it"
- [ ] InitCommands to run on `vessel start` over SSH
- [x] `vessel stop` command to cleanup / stop mutagen sessions
- [ ] Starting a session twice results in double sessions (sync + forward)
    - Use labels and --label-selector, although machine-readable output is a coming feature of mutagen
- [ ] GitHub Actions: Releaser (goreleaser/goreleaser-action@v2)
- [ ] Download a set version of mutagen for the current OS to embed (~12M)