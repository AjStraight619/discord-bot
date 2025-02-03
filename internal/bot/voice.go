package bot

import (
	"log"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

// TODO: Maybe implement struct here to control channels more often

type PersonVoice struct {
	userID   string
	username string
	// recvChan chan *discordgo.Packet
	// sendChan chan *discordgo.Packet
}

func (v *PersonVoice) NewVoice(userID, username string) *PersonVoice {
	return &PersonVoice{
		userID:   userID,
		username: username,
		// recvChan: make(chan *discordgo.Packet, 2),
		// sendChan: make(chan *discordgo.Packet, 2),
	}

}

func (b *BotController) Echo() {

	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(b.VoiceConn, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(b.VoiceConn, send)

	b.VoiceConn.Speaking(true)
	defer b.VoiceConn.Speaking(false)

	for {

		p, ok := <-recv
		if !ok {
			return
		}

		send <- p.PCM
	}
}

func (b *BotController) ListenVoice(msg *discordgo.MessageCreate) {
	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(b.VoiceConn, recv)

	err := b.VoiceConn.Speaking(true)
	if err != nil {
		log.Printf("Error setting speaking state: %v", err)
		return
	}
	defer func() {
		_ = b.VoiceConn.Speaking(false)
	}()

	log.Println("🎙️ Now listening for incoming audio...")

	for {
		select {
		case p, ok := <-recv:
			if !ok {
				log.Println("Voice receive channel closed.")
				return
			}

			log.Printf("Received %d PCM samples", len(p.PCM))
			if len(p.PCM) > 0 {
				log.Printf("First sample: %d", p.PCM[0])
			}
		case <-time.After(30 * time.Second):
			log.Println("No audio received for 30 seconds, stopping listening.")
			return
		}
	}
}
