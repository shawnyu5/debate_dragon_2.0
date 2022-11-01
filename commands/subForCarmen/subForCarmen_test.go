package subforcarmen_test

import (
	"time"

	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	subforcarmen "github.com/shawnyu5/debate_dragon_2.0/commands/subForCarmen"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var _ = Describe("Subforcarmen", func() {
	Context("CheckMessageAuthor", func() {
		It("Should recognize camen message", func() {
			c := utils.Config{
				Development: false,
				Ivan: struct {
					Emotes []struct {
						Name         string "json:\"name\""
						FileLocation string "json:\"fileLocation\""
					} "json:\"emotes\""
				}{},
				SubForCarmen: struct {
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int64    "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int64    "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{
					CarmenID:          "12345",
					CoolDown:          0,
					GuildID:           "",
					MessageLimit:      0,
					SubscribersRoleID: "",
					IgnoredChannels:   []string{},
				},
			}
			mess := &discordgo.Message{
				Content: "hello",
				Author: &discordgo.User{
					ID: "12345",
				},
			}

			Expect(subforcarmen.CheckMessageAuthor(mess, c.SubForCarmen.CarmenID)).To(BeTrue())
		})

		It("should not recognize non-carmen message", func() {
			c := utils.Config{
				Development: false,
				Ivan: struct {
					Emotes []struct {
						Name         string "json:\"name\""
						FileLocation string "json:\"fileLocation\""
					} "json:\"emotes\""
				}{},
				SubForCarmen: struct {
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int64    "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int64    "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{
					CarmenID:          "jjjjj",
					CoolDown:          0,
					GuildID:           "",
					MessageLimit:      0,
					SubscribersRoleID: "",
					IgnoredChannels:   []string{},
				},
			}
			mess := &discordgo.Message{
				Content: "hello",
				Author: &discordgo.User{
					ID: "12345",
				},
			}

			Expect(subforcarmen.CheckMessageAuthor(mess, c.SubForCarmen.CarmenID)).To(BeFalse())
		})
	})

	Context("IsCoolDown", func() {
		It("should return true if the message is within the cooldown period", func() {
			// 10 mins ago from now
			tenMinsAgo := time.Now().Add(time.Duration(-10) * time.Minute)
			mess := &discordgo.Message{
				Content:   "hello",
				Timestamp: time.Now(),
				Author:    &discordgo.User{ID: "12345"},
			}

			s := subforcarmen.State{
				LastNotificationTime: time.Time{},
				LastMessageTime:      tenMinsAgo,
				Counter:              0,
			}
			subforcarmen.CarmenState = s
			Expect(subforcarmen.IsCoolDown(mess, 20)).To(BeTrue())
			// check if last message time stamp is updated
			Expect(subforcarmen.CarmenState.LastMessageTime).To(BeIdenticalTo(mess.Timestamp))
		})

		It("should return false if the message is not within the cooldown period", func() {
			// 10 mins ago from now
			tenMinsAgo := time.Now().Add(time.Duration(-10) * time.Minute)
			mess := &discordgo.Message{
				Content:   "hello",
				Timestamp: time.Now(),
				Author: &discordgo.User{
					ID: "12345",
				},
			}

			carmenState := subforcarmen.State{
				LastNotificationTime: time.Time{},
				LastMessageTime:      tenMinsAgo,
				Counter:              0,
			}
			subforcarmen.CarmenState = carmenState
			Expect(subforcarmen.IsCoolDown(mess, 5)).To(BeFalse())
		})
	})

	Context("IncreaseCounter", func() {
		It("Should increase counter", func() {
			// 5 mins ago from now
			fiveMinsAgo := time.Now().Add(time.Duration(-5) * time.Minute)
			state := subforcarmen.State{
				LastNotificationTime: time.Time{},
				LastMessageTime:      fiveMinsAgo,
				Counter:              0,
			}
			mess := &discordgo.Message{
				Timestamp: time.Now(),
			}
			subforcarmen.CarmenState = state

			Expect(subforcarmen.IncreaseCounter(mess)).To(BeTrue())
			// counter should have increased by one
			Expect(subforcarmen.CarmenState.Counter).To(Equal(1))
			Expect(subforcarmen.CarmenState.LastMessageTime).To(BeIdenticalTo(mess.Timestamp))
		})

		It("Should reset counter to 0", func() {
			// 5 mins ago from now
			tenMinsAgo := time.Now().Add(time.Duration(-15) * time.Minute)
			state := subforcarmen.State{
				LastNotificationTime: time.Time{},
				LastMessageTime:      tenMinsAgo,
				Counter:              5,
			}
			mess := &discordgo.Message{
				Timestamp: time.Now(),
			}
			subforcarmen.CarmenState = state

			Expect(subforcarmen.IncreaseCounter(mess)).To(BeFalse())
			// counter should have increased by one
			Expect(subforcarmen.CarmenState.Counter).To(Equal(0))
		})
	})
})
