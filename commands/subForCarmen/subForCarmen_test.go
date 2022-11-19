package subforcarmen_test

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	subforcarmen "github.com/shawnyu5/debate_dragon_2.0/commands/subForCarmen"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
	"github.com/spf13/afero"
)

// CreateMockConfig creates a mock config.json
// fs     : the file system to write to
// c      : the config object to create
// returns: error if any
func CreateMockConfig(fs afero.Fs, c utils.Config) error {
	j, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return afero.WriteFile(fs, "config.json", j, 0644)
}

var _ = Describe("Subforcarmen", func() {
	// Context("Listen()", func() {
	// It("Should recognize carmen message", func() {
	// // use mock fs
	// fs := afero.NewMemMapFs()
	// utils.AppFs = fs
	// c := utils.Config{
	// SubForCarmen: struct {
	// CarmenID          string   "json:\"carmenId\""
	// CoolDown          int64    "json:\"coolDown\""
	// GuildID           string   "json:\"guildID\""
	// MessageLimit      int      "json:\"messageLimit\""
	// SubscribersRoleID string   "json:\"subscribersRoleID\""
	// IgnoredChannels   []string "json:\"ignoredChannels\""
	// }{CoolDown: },
	// }
	// CreateMockConfig(fs)
	// carmenID := "12345"
	// mess := &discordgo.Message{
	// Content: "hello",
	// GuildID: "apple",
	// Author: &discordgo.User{
	// ID: "12345",
	// },
	// }

	// Expect(subforcarmen.Listen(mess, carmenID, "apple")).To(BeTrue())
	// })

	// It("should not recognize non-carmen message", func() {
	// carmenID := "12345"
	// mess := &discordgo.Message{
	// Content: "hello",
	// GuildID: "apple",
	// Author: &discordgo.User{
	// ID: "12345",
	// },
	// }

	// Expect(subforcarmen.Listen(mess, carmenID, "apple")).To(BeFalse())
	// })
	// })

	Context("IsValidMessage()", func() {
		utils.AppFs = afero.NewMemMapFs()
		// It("Should recognize a carmen message", func() {
		// c := utils.Config{
		// Token:       "",
		// TokenDev:    "",
		// ClientID:    "",
		// GuildID:     "",
		// LogLevel:    "",
		// Development: false,
		// Ivan: struct {
		// Emotes []struct {
		// Name         string "json:\"name\""
		// FileLocation string "json:\"fileLocation\""
		// } "json:\"emotes\""
		// }{},
		// SubForCarmen: struct {
		// On                bool     "json:\"on\""
		// CarmenID          string   "json:\"carmenId\""
		// CoolDown          int      "json:\"coolDown\""
		// GuildID           string   "json:\"guildID\""
		// MessageLimit      int      "json:\"messageLimit\""
		// SubscribersRoleID string   "json:\"subscribersRoleID\""
		// IgnoredChannels   []string "json:\"ignoredChannels\""
		// }{CarmenID: "12345"},
		// }
		// })

	})
	Context("IsCoolDown()", func() {
		It("should return true if the message is within the cooldown period", func() {
			c := utils.Config{
				SubForCarmen: struct {
					On                bool     "json:\"on\""
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int      "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int      "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{CoolDown: 60},
			}
			// last notification is sent 10 mins ago
			subforcarmen.CarmenState.LastNotificationTime = time.Now().Add(time.Duration(-10) * time.Minute)
			CreateMockConfig(utils.AppFs, c)

			mess := &discordgo.Message{
				Content:   "hello",
				Timestamp: time.Now(),
				Author:    &discordgo.User{ID: "12345"},
			}

			// cool down is 60 mins
			Expect(subforcarmen.IsCoolDown(mess)).To(BeTrue())
		})

		It("should return false if the message is not within the cooldown period", func() {
			c := utils.Config{
				SubForCarmen: struct {
					On                bool     "json:\"on\""
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int      "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int      "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{CoolDown: 1},
			}
			CreateMockConfig(utils.AppFs, c)

			mess := &discordgo.Message{
				Content:   "hello",
				Timestamp: time.Now(),
				Author: &discordgo.User{
					ID: "12345",
				},
			}

			// last notification is sent 10 hours ago
			subforcarmen.CarmenState.LastNotificationTime = time.Now().Add(time.Duration(-10) * time.Hour)
			// cool down is 1 mins
			Expect(subforcarmen.IsCoolDown(mess)).To(BeFalse())
		})
	})

	Context("IncreaseCounter()", func() {
		It("Should increase counter when message is sent within 5 mins", func() {
			c := utils.Config{
				SubForCarmen: struct {
					On                bool     "json:\"on\""
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int      "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int      "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{MessageLimit: 5},
			}
			CreateMockConfig(utils.AppFs, c)

			// last message was sent 5 mins ago
			subforcarmen.CarmenState.LastMessageTime = time.Now().Add(time.Duration(-4) * time.Minute)
			subforcarmen.CarmenState.Counter = 0
			mess := &discordgo.Message{
				Timestamp: time.Now(),
			}
			Expect(subforcarmen.IncreaseCounter(mess)).To(BeTrue())

			// counter should have increased by one
			Expect(subforcarmen.CarmenState.Counter).To(Equal(1))
		})

		It("Should reset counter to 0 when message is more than 5 mins ago", func() {
			c := utils.Config{
				SubForCarmen: struct {
					On                bool     "json:\"on\""
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int      "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int      "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{MessageLimit: 5},
			}
			CreateMockConfig(utils.AppFs, c)

			state := subforcarmen.State{
				LastNotificationTime: time.Time{},
				LastMessageTime:      time.Now().Add(time.Duration(-4) * time.Minute),
				Counter:              5,
			}
			subforcarmen.CarmenState = state

			mess := &discordgo.Message{
				Timestamp: time.Now(),
			}

			Expect(subforcarmen.IncreaseCounter(mess)).To(BeFalse())
			// counter should have increased by one
			Expect(subforcarmen.CarmenState.Counter).To(Equal(0))
		})
	})
	Context("ShouldTriggerNotification()", func() {
		It("Should return true when counter < messageLimit", func() {
			subforcarmen.CarmenState.Counter = 5
			Expect(subforcarmen.ShouldTriggerNotification(6))
		})
	})

	Context("IsIgnoredChannel()", func() {
		It("Should detect an ignored channel", func() {
			c := utils.Config{
				SubForCarmen: struct {
					On                bool     "json:\"on\""
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int      "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int      "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{
					IgnoredChannels: []string{"12345", "jflkdsjf"},
				},
			}
			CreateMockConfig(utils.AppFs, c)

			Expect(subforcarmen.IsIgnoredChannel("12345")).To(BeTrue())
		})

		It("Should detect an non ignored channel", func() {
			c := utils.Config{
				SubForCarmen: struct {
					On                bool     "json:\"on\""
					CarmenID          string   "json:\"carmenId\""
					CoolDown          int      "json:\"coolDown\""
					GuildID           string   "json:\"guildID\""
					MessageLimit      int      "json:\"messageLimit\""
					SubscribersRoleID string   "json:\"subscribersRoleID\""
					IgnoredChannels   []string "json:\"ignoredChannels\""
				}{
					IgnoredChannels: []string{"12345", "jflkdsjf"},
				},
			}
			CreateMockConfig(utils.AppFs, c)

			Expect(subforcarmen.IsIgnoredChannel("jjjjj")).To(BeFalse())
		})
	})
})
