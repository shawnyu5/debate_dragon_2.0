package snipe_test

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shawnyu5/debate_dragon_2.0/commands/snipe"
)

var _ = Describe("Snipe", func() {
	It("should keep track of a new message", func() {
		mess := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content: "hello world",
			},
		}
		snipe.TrackMessage(mess)
		Expect(len(snipe.AllMessages)).To(Equal(1))
	})

	It("should not have more than 100 messages", func() {

		for i := 0; i < 105; i++ {
			mess := &discordgo.MessageCreate{
				Message: &discordgo.Message{
					ID:      strconv.Itoa(i),
					GuildID: "12345",
					Content: "hello world",
				},
			}
			snipe.TrackMessage(mess)
		}
		Expect(len(snipe.AllMessages["12345"])).To(Equal(100))
	})

	It("should remove the oldest message", func() {
		for i := 0; i < 105; i++ {
			mess := &discordgo.MessageCreate{
				Message: &discordgo.Message{
					ID:        strconv.Itoa(i),
					GuildID:   "12345",
					Content:   "hello world",
					Timestamp: time.Now(),
				},
			}
			snipe.TrackMessage(mess)
		}

		// add a really old message.
		oldMessage := discordgo.MessageCreate{
			Message: &discordgo.Message{
				ID:        "old",
				GuildID:   "12345",
				Content:   "IM OLD",
				Timestamp: time.Now().Add(time.Duration(-3) * time.Hour),
			},
		}
		snipe.TrackMessage(&oldMessage)

		// the old message should be gone
		Expect(snipe.AllMessages["12345"]["old"]).To(BeEquivalentTo(discordgo.Message{}))
		Expect(len(snipe.AllMessages["12345"])).To(Equal(100))
	})
})
