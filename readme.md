## TODO
- [ ] Global debug flag
    - Print output from mutagen and other exec.Cmd calls
- [ ] SSH config edits ~/.ssh/config to add an alias?!
    - Prints out helpful message "Add this yourself, or press 'y' for us to do it"
- [ ] InitCommands to run on `vessel start`
- [ ] `vessel stop` command to cleanup / stop mutagen sessions
- [ ] Starting a session twice results in double sessions (sync + forward)
    - Use labels and --label-selector, although machine-readable output is a coming feature of mutagen