# Kavaca

A discord bot that forwards DMs sent to it to a configured user. Source code is meant to be simple, and hence it pulls in only one direct dependency.

While it does not reveal any information about the owner of the bot, or who configured it, it DOES NOT promise to anonymize your presence on discord.

All code in this repository is licensed under the GPLv3 license.

## Slash Commands

Slash command integration has arrived. The bot function solely using slash commands, which means that my substandard argument parsing is no longer needed. This also eliminates the need for a command prefix. The [legacy-no-slashcommands](https://github.com/fisik-yum/kavaca/tree/legacy-no-slashcommands) branch serves as a snapshot of the initial version in case someone really wants to use that.

## Setup

clone the repository using `git clone`, and compile it using `go build`

Set up your discord application and create a config.json beside the binary and format it like this

Leave defaultchannel blank `""` to use the owner's DMs instead. 

```
{
"id": "<your user ID>",
"token": "<bot token>",
"defaultChannel": "<default message channel>"
}
```


And set up bindings.json like this

`[]`

## Usage

The bot NEEDS the owner's ID and bot token to function. A prefix is recommended.
There are a few basic commands.

`/bind <userID> <channelID>` binds a users messages to one channel that the bot has access to.

`/listbinds`shows the current binds in place in the console. mainly for testing.

`/reset <userID>` reset user's bind to default channel.

`/savebinds`force save binds. kavaca saves binds automatically at shutdown.

Reply to a message from a user by replying to that message. The message will DMed to the user.
