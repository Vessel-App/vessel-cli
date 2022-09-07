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

1. [Install `flyctl`](https://fly.io/docs/getting-started/installing-flyctl/) and create a Fly account via `fly auth signup`
2. Install `vessel`

For Mac/Linux, you can install `vessel` this way:

```bash
# Don't forget to follow instructions to add ~/.vessel/bin to your $PATH
curl https://vessel.fly.dev/stable/install.sh | sh
```

> 😭 I haven't built nor tested this on Windows yet.

2️⃣ Then authenticate Vessel

```bash
# If you already logged in with `fly auth login`
vessel auth

# Else, if you have a Fly API token
vessel auth -t YOUR_TOKEN_HERE
```

This gives Vessel access to the Fly API token that you'd like to use. If the API token is specific to your default organization, we'll use that. Otherwise, we'll prompt to ask which organization to use.
Each Fly.io organization is billed separately.

3️⃣ Then head to a code base and run initialize your project.

```bash
# Probably a Laravel project
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
vessel -- "npm install && npm run build"

# Open http://localhost:8000 in your browser
vessel open

# When you're done:
vessel stop
```

## Project Configuration

Your project will contain a `vesssel.yml` file. You can customize this to forward additional ports. By default, it will forward
localhost port `8000` to the development environment's port `80`.

## I'm a Windows user

I haven't had a chance to test this on Windows.

Theoretically I've made the code work for some subset of Windows users (based on [Mutagen's requirements](https://mutagen.io/documentation/transports/ssh#windows)).

## What is this thing exactly?

This is the result of some fun I had using Fly's Machines API to make remote development environment.

It's inspired by:

1. Me having a new computer and not wanting to install so much crap into it
2. Amos's [article on remote dev on Fly](https://fasterthanli.me/articles/remote-development-with-rust-on-fly-io) (This project is a bit different, but I totally stole his Rust code)
