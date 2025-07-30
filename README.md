# Gator

A command-line tool used for aggregating RSS feeds.


## Installation

Postgres and Go are required for running Gator.

```bash
  go install gator
  gator <command> <args>
```

Create and set up `.gatorconfig.json` inside your home directory to point to your database and current user like so:

```bash
{"db_url":"postgres://username:password@host:port/database?sslmode=disable","current_user_name":}
```
## Usage

Available commands:

`login <user>` 
> Log in as a user

`register <user>` 
> Register a new user

`reset` 
> Resets all users

`users` 
> Lists all users

`agg <duration>` 
> Aggregates posts from followed feeds every duration, e.g agg 5m for every 5 minutes

`feeds` 
> Lists all feeds along with the user who added it

`addfeed <title> <url>` 
> Add a feed

`follow <url>` 
> Follow a feed

`following` 
> Shows the feeds the logged-in user is following

`unfollow <url>` 
> Unfollow a feed

`browse (optional)<number>` 
> Shows <number> of posts from your followed feeds, if <number> isn't provided it defaults to 2
