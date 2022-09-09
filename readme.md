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

## How to Use This

Here's how to install and use Vessel.

### 1️⃣ Install

First, install dependencies:

1. [Install `flyctl`](https://fly.io/docs/getting-started/installing-flyctl/) and create a Fly account via `fly auth signup`
2. Install `vessel`

For Mac/Linux, you can install `vessel` this way:

```bash
# 👉 Don't forget to follow instructions to add ~/.vessel/bin to your $PATH
curl https://vessel.fly.dev/stable/install.sh | sh
```

### 2️⃣ Authenticate
Then authenticate Vessel

```bash
# If you already logged in with `fly auth login`
vessel auth

# Else, if you have a Fly API token
vessel auth -t YOUR_TOKEN_HERE
```

This gives Vessel access to the Fly API token that you'd like to use. If the API token is specific to your default organization, we'll use that. Otherwise, we'll prompt to ask which organization to use.
Each Fly.io organization is billed separately.

### 3️⃣ Start a project

Head to a code base and run initialize your project.

```bash
# We need to be able to communicate to Fly's API
# Run this in another terminal, as it's a long-running process
# It's only needed for the init step
fly machine api-proxy

# Probably a Laravel project
cd ~/Code/some-laravel-project
vessel init
```

### 🔁 Usage

Once that's finished, run the `start` command to enable file syncing / port forwarding

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

Your project will contain a `vesssel.yml` file. You can customize this to forward additional ports.

By default, it will forward localhost port `8000` to the development environment's port `80`.

## SSH and Commands

You can run one-off commands and SSH into your environments.

> File syncing (and port forwarding) is done via [mutagen](https://github.com/mutagen-io/mutagen), which also works over SSH. 

### One-Off Commands

Vessel supports one-off commands two ways:

1. `vessel -- composer install` - This will run `composer install` after connecting to your dev env
2. `vessel cmd npm install` - Similarly, This will run `npm install` after connecting to the dev env

Commands are run from the `~/app` directory within the dev environment. If you run a one-off command without first syncing, you may get
errors about the `~/app` working directory not existing.

### SSH

This project configures an easy way to SSH into the dev environment. 
After you run `vessel init`, check your `~/.ssh/config` file to see a new alias created there. You should be able to SSH to the server without the `vessel` command:

```bash
# SSH via vessel
vessel ssh

# SSH outside of Vessel
# This will match a host set in ~/.ssh/config
ssh vessel-<my-project-name>
```

## Global Configuration

You'll find global configuration and a debug log file in `~/.vessel`:

* `~/.vessel/config.yml` - Configuration including your Fly API token and the Fly organization used
* `~/.vessel/debug.log` - Logs to help troubleshoot issues
* `~/.vessel/<your-project>` - A directory containing SSH keys used to access your dev environment

## Debugging

Use `LOG_LEVEL=debug vessel ...` to get a bit more information output to your `~/.vessel/debug.log` file.

Try adding the `-v` flag to any `vessel` command to get complete errors output directly to your console, e.g. `vessel -v init`.

## I'm a Windows user

I haven't had a chance to test this on Windows.

Theoretically I've made the code work for some subset of Windows users (based on [Mutagen's requirements](https://mutagen.io/documentation/transports/ssh#windows)).

## Ephemeral Environments

The environments will shut down after 5 minutes of inactivity. In activity means no active SSH connections (file syncing via `vessel start` counts as an active SSH session).

This means you won't be charged for usage when environments are not in use. However, the environments are ephemeral. When you start an environment back up,
**it's as if you're starting from a blank slate**. Your files will get re-synced when you run `vessel start` next, but you'll need to ensure any data you need
is put back into place (for now!).

I use `sqlite` for all development in this fashion (for as long as I can get away with it!), as it lets me easily have my "state" synced to the dev environment.

## Making API Calls to Fly.io

During the `init` step, we used `fly machine api-proxy`. This proxies requests from `localhost:4280` to `_api.internal:4280`.
This `_api.internal` address is actually a private network address that works from within Fly.io's private networks.

The other way to talk to Fly.io's Machines API is to log into your organizations private network via VPN.
Instructions on setting up [Fly.io's Private Networt VPN are found here](https://fly.io/docs/reference/private-networking/#private-network-vpn).

If you use that method instead, you can:

1. Talk to your VM's directly via `*.internal` hostnames
2. Use `vessel` by setting the `FLY_HOST` environment variable

```bash
 export FLY_HOST="_api.internal"
 vessel init
```

## What is this thing exactly?

This is the result of some fun I had using Fly's Machines API to make remote development environment.

It's inspired by:

1. Me having a new computer and not wanting to install so much crap into it
2. Amos's [article on remote dev on Fly](https://fasterthanli.me/articles/remote-development-with-rust-on-fly-io) (This project is a bit different, but I totally stole his Rust code)
