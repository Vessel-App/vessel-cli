## TODO
- [x] Global debug flag (env var)
    - Print output from mutagen and other exec.Cmd calls
- [x] SSH config edits ~/.ssh/config to add an alias?!
    - Prints out helpful message "Add this yourself, or press 'y' for us to do it"
- [x] `vessel stop` command to cleanup / stop mutagen sessions
- [x] Starting a session twice results in double sessions (sync + forward)
- [ ] GitHub Actions: Releaser (goreleaser/goreleaser-action@v2)
- [ ] Download a set version of mutagen for the current OS to embed (~12M)
- [ ] Select (default) environment - php version, node version
- [ ] InitCommands to run on `vessel start` over SSH
- [ ] Global command to see "status" - what apps have mutagen sessions open?
- [ ] Hit the Vessel API to get current user and user's teams
    - Populate the team GUID (with a yaml comment of the team name, we're not animals)
    - Prompt user to ask which team this is related to (if more than one team associated with their account)
- [ ] Hit the Vessel API during init to create a machine within a team's app

> This no longer really needs to be Golang as we're not using Mutagen libraries but rather shelling out out to the Mutagen CLI.
> 
> But that doesn't seem like a useful refactoring at this point.