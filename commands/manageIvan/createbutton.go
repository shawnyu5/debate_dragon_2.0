package manageIvan

import "github.com/bwmarrin/discordgo"

// createBanButton create a ban button.
// disable: whether the button should be disabled.
func createBanButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Ban",
		Style:    discordgo.DangerButton,
		Disabled: disable,
		Emoji:    discordgo.ComponentEmoji{},
		CustomID: startBanProcessID,
	}
}

// createDontBanButton create a dont ban button
// disable: whether the button should be disabled.
func createDontBanButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Dont ban...",
		Style:    discordgo.PrimaryButton,
		Disabled: disable,
		Emoji:    discordgo.ComponentEmoji{},
		CustomID: dontBanIvanID,
	}
}

// createJumpScareButton create a jump scare button
// disable: whether the button should be disabled.
func createJumpScareButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Jump scare",
		Style:    discordgo.SuccessButton,
		Disabled: disable,
		Emoji:    discordgo.ComponentEmoji{},
		CustomID: banJumpScareID,
	}
}

// createKickButton creat a button to kick a user
// disable: whether the button should be disabled.
func createKickButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Kick",
		Style:    discordgo.DangerButton,
		Disabled: disable,
		Emoji:    discordgo.ComponentEmoji{},
		CustomID: kickID,
	}
}

// CreateAllButtons create a row of all buttons for this interaction.
// disable: whether the buttons should be disabled.
// return: a row of buttons.
func CreateAllButtons(disable bool) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			createBanButton(disable),
			createDontBanButton(disable),
			createJumpScareButton(disable),
			createKickButton(disable),
		},
	}
}
