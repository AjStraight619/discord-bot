package bot

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"log"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

// TODO: Maybe implement struct here to control channels more often

type PersonVoice struct {
	userID   string
	username string
}

type VoiceCommandHandler struct {
	Bot *BotController
}

type ListenCommand struct{}

func (lc ListenCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	b.ListenVoice(msg)
}

func (lc ListenCommand) Help() string {
	return "!listen - Listen to voice channel."
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

func (b *BotController) VoiceToTextStream() {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create speech client: %v", err)
	}
	defer client.Close()

	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatalf("Failed to open streaming recognize: %v", err)
	}

	// Send the initial configuration.
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000, // Adjust to match your PCM sample rate.
					LanguageCode:    "en-US",
				},
				InterimResults: true, // Set to false if you only want final results.
			},
		},
	}); err != nil {
		log.Fatalf("Could not send config: %v", err)
	}

	// Create a channel to receive the recognized text.
	textChan := make(chan string)

	// Process the responses from the API in a separate goroutine.
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error receiving recognition: %v", err)
				break
			}
			// Loop through the results and send the transcript to textChan.
			for _, result := range resp.Results {
				for _, alt := range result.Alternatives {
					textChan <- alt.Transcript
				}
			}
		}
		close(textChan)
	}()

	// Now, capture audio from Discord and stream it to the API.
	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(b.VoiceConn, recv)

	b.VoiceConn.Speaking(true)
	defer b.VoiceConn.Speaking(false)

	// In a loop, convert and send audio chunks.
	go func() {
		for packet := range recv {
			// Convert the PCM samples (int16 slice) to bytes.
			audioBytes := ConvertPCMToBytes(packet.PCM)
			// Send the audio data chunk.
			err := stream.Send(&speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
					AudioContent: audioBytes,
				},
			})
			if err != nil {
				log.Printf("Error sending audio chunk: %v", err)
				return
			}
		}
		// After channel is closed, close the stream for the API.
		err := stream.CloseSend()
		if err != nil {
			log.Printf("Error closing stream: %v", err)
		}
	}()

	// Consume recognized text and do something with it (e.g., send to Discord channel).
	for transcript := range textChan {
		log.Printf("Recognized: %s", transcript)
		// For example, you could send this transcript as a message:
		// b.Session.ChannelMessageSend("YOUR_CHANNEL_ID", transcript)
	}
}

func ConvertPCMToBytes(samples []int16) []byte {
	buf := new(bytes.Buffer)
	// Write each int16 sample as little-endian.
	for _, sample := range samples {
		// Check error if needed, here we assume no error.
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}
	return buf.Bytes()
}

// func (b *BotController) ProcessVoiceCommand(transcript, channelID string) {
// 	// Convert transcript to lower-case for easier matching.
// 	t := strings.ToLower(transcript)
//
// 	// Check for a specific command trigger.
// 	if strings.HasPrefix(t, "command:") {
// 		// Remove the "command:" part and process as a text command.
// 		cmdText := strings.TrimSpace(t[len("command:"):])
// 		b.ProcessCommand(cmdText, channelID)
// 		return
// 	}
//
// 	// Alternatively, do keyword matching.
// 	if strings.Contains(t, "news") {
// 		// You might need to extract the country code from the transcript.
// 		// For a very simple version:
// 		var countryCode string
// 		// For example, if the transcript is "news us", then:
// 		parts := strings.Fields(t)
// 		for _, part := range parts {
// 			if len(part) == 2 { // crude check: a two-letter country code.
// 				countryCode = part
// 				break
// 			}
// 		}
// 		if countryCode != "" {
// 			// Construct the command string like "!news us"
// 			cmdText := fmt.Sprintf("!news %s", countryCode)
// 			b.ProcessCommand(cmdText, channelID)
// 			return
// 		}
// 	}
//
// 	if strings.Contains(t, "ping") {
// 		b.ProcessCommand("!ping", channelID)
// 		return
// 	}
//
// 	// Handle more commands similarly...
// 	// If no commands are recognized, you could send feedback to the user.
// 	b.Session.ChannelMessageSend(channelID, "Sorry, I didn't understand that command.")
// }
