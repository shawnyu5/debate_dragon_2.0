# Debate dragon 2.0

A blazingly fast discord bot written in Go to burn your debate foes to the ground

## Commands


`/dd text:<text>` - generate a dragon drawing, with `text` imposed into the speech bubble

<img src="media/img/dragon_drawing.png" alt="dragon drawing" width="300">

`/insult user:<userName> anonymous: <true | false>`: Send an insult to `userName`. `anonymous` determines if you'll be shown as the person that executed the command

`/rmp profname: <name>`: search for a professor by name from Seneca college. If more than 1 professor by that name at Seneca, a select menu will be displayed to prompt to select the prof to display the ratings of. Other wise, the prof's ratings will be displayed.

`/subforcarmen subscribe: true` - Give the user the subscriber role

`/subforcarmen subscribe: false` - Remove the subscriber role from the user

This is a command that will ping a subscriber role, if a selected user is talking too much, for everyone to witness the drama that is happening right now.

The user and guild which this command watches can be configured in `config.json`, under `subForCarmen.carmenID`, and `subForCarmen.guildID`. The subscriber role is under `subForCarmen.subscriberRoleID`

`subForCarmen.on` can be used to toggle notifications


## Development

The bot requires a `config.json` to set environment variables. Refer to `config.example.jsonc` for all values needed in `config.json`.

**Note** the actual config should be a json file. The `config.example.jsonc` is `jsonc` only to allow comments
