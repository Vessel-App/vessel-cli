## Remote Dev Environments
Remote dev environments that feel local. Powered by [Fly.io](https://fly.io).

<!-- gif showing it off here -->

## What?

Installing tons of development stuff onto laptop sucks. Let's not but say we did.

This project:

1. Creates an app on Fly.io that has your dependencies in it.
2. Syncs your local code into the app on-the-fly
3. Forwards defined ports (e.g. open `localhost:8000` in your browser)

> The Fly app turns off automatically when you disconnect from it, so it's not running 24/7.
>
> It'll turn back on automatically when you re-connect.

## Get Started

1️⃣ First, install dependencies:

1. [Install `flyctl`](https://fly.io/docs/getting-started/installing-flyctl/) and create a Fly account
2. [Install `mutagen`](https://mutagen.io/documentation/introduction/installation) to sync code / forward ports
3. Install this thing by uhhhhh, downloading the right binary from [releases](https://github.com/Vessel-App/vessel-cli/releases). I'll figure something nicer out.
    - M1/M2 Macs should grab the `Darwin ARM64` file
    - Intel Macs should grab the `Darwin AMD64` file
    - Linux desktop users already know what to pick. You poor, wretched sods

2️⃣ Then install `vessel` from the [releases](https://github.com/Vessel-App/vessel-cli/releases) page. (I'll improve this soon).

- M1/M2 Macs should grab the `Darwin ARM64` file
- Intel Macs should grab the `Darwin AMD64` file
- Linux desktop users already know what to pick. You poor, wretched sods

3️⃣ Then authenticate.

```bash
# If you logged in with `fly auth login`
vessel auth

# If you have a Fly API token
vessel auth -t YOUR_TOKEN_HERE
```

4️⃣ Then head to a code base and run initialize your project.

```bash
cd ~/Code/some-laravel-project
vessel init
```

🔁 Once that's finished, run the `start` command to enable file syncing / port forwarding

```bash
# Start syncing/port forwarding
vessel start

# Run some commands to get dependencies in your server
## (Dependencies aren't synced)
## You can also run `vessel cmd <some command here>`
vessel -- composer install
vessel -- npm install && npm run build

# Open localhost:8000 in your browser
vessel open

# When you're done:
vessel stop
```

## I'm a Windows user

I haven't had a chance to test this on Windows.

Theoretically I've made the code work for some subset of Windows users (based on [Mutagen's requirements](https://mutagen.io/documentation/transports/ssh#windows)).

## What is this thing exactly?

This is the result of some fun I had using Fly's Machines API to make remote development environment.

It's inspired by:

1. Me having a new computer and not wanting to install so much crap into it
2. Amos's [article on remote dev on Fly](https://fasterthanli.me/articles/remote-development-with-rust-on-fly-io) (This project is a bit different, but I totally stole his Rust code)

## TODO
- [ ] Download a set version of mutagen for the current OS to embed (~12M)
- [ ] Config for base URL
- [ ] Forward all defined ports in `.vessel.yml`