package bot

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type TimeoutCommand struct{}

func (tc TimeoutCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	if len(options) != 1 {
		b.displayCmdError(msg.ChannelID, "Please specify a timeout duration: !timeout 30 (This will set a timeout for 30 minutes)")
		return
	}
	timeout := options[0]
	num, err := strconv.Atoi(timeout)

	if err != nil {
		b.displayCmdError(msg.ChannelID, "Please input an number")
		return
	}

	b.TimeoutDuration = time.Duration(num) * time.Minute
	b.ResetTimeout()
	b.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Timeout duration set to %d minutes.", num))

}

func (tc TimeoutCommand) Help() string {
	return "!timeout - Set a timeout for the bot that will trigger it to leave on inactivity"
}

func (b *BotController) ResetTimeout() {
	log.Println("ResetTimeout: Called to restart the inactivity timer.")

	// If there's an existing timer, stop it.
	if b.inactivityTimer != nil {
		log.Println("ResetTimeout: An existing timer was found; attempting to stop it.")
		// Stop returns false if the timer has already expired.
		if !b.inactivityTimer.Stop() {
			log.Println("ResetTimeout: Timer had already expired; draining the timer's channel.")
			// Drain the timer's channel if needed.
			select {
			case <-b.inactivityTimer.C:
				log.Println("ResetTimeout: Drained one value from the timer's channel.")
			default:
				log.Println("ResetTimeout: No value to drain from the timer's channel.")
			}
		} else {
			log.Println("ResetTimeout: Timer stopped successfully.")
		}
	} else {
		log.Println("ResetTimeout: No existing timer found. Creating a new one.")
	}

	endTime := time.Now().Add(b.TimeoutDuration)

	// Start a new timer with the configured duration.
	b.inactivityTimer = time.AfterFunc(b.TimeoutDuration, func() {
		log.Println("Timeout reached. Executing timeout action.")
		b.OnTimeout()
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			remaining := time.Until(endTime)
			if remaining <= 0 {
				break
			}
			log.Printf("Countdown: %d seconds remaining", int(remaining.Seconds()))
			// Wait for the next tick.
			<-ticker.C
		}
	}()

	log.Printf("ResetTimeout: New timer started with a timeout duration of %v.\n", b.TimeoutDuration)
}

func (b *BotController) OnTimeout() {
	log.Println("Timeout reached (inactivity)")
	b.LeaveVoiceChannel()
}
