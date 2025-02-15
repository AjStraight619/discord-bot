package bot

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// Song represents a queued song with its download status.
type Song struct {
	URL         string // The original YouTube URL.
	FilePath    string // The local file path once downloaded.
	Downloading bool   // True if the song is still downloading.
	DownloadErr error  // Holds any error that occurred during download.
}

// SongCommand is the command that triggers playing songs.
type SongCommand struct{}

func (sc SongCommand) Execute(b *BotController, msg *discordgo.MessageCreate, options []string) {
	if len(options) < 1 {
		b.displayCmdError(msg.ChannelID, "⚠ Usage: `!play <music_link>`")
		return
	}
	// Process all URLs passed in (allows multiple songs at once)
	b.Play(options, msg)
}

func (sc SongCommand) Help() string {
	return "!play <music_link> - Plays the specified YouTube music link(s)."
}

// Play adds one or more songs to the queue and starts downloading them concurrently.
func (b *BotController) Play(options []string, msg *discordgo.MessageCreate) {
	// Loop over all provided URLs.

	for _, youtubeURL := range options {
		log.Println("🎥 YouTube URL Received:", youtubeURL)

		// Create a new Song instance and mark it as downloading.
		song := &Song{
			URL:         youtubeURL,
			Downloading: true,
		}

		// Append the song to the queue (FIFO order is preserved).
		b.musicQueue = append(b.musicQueue, song)
		b.displayCmdError(msg.ChannelID, fmt.Sprintf("🎵 Added to queue: %s", youtubeURL))

		// Spawn a goroutine to download this song concurrently.
		go func(s *Song) {
			filePath, err := downloadYouTubeAudio(s.URL)
			if err != nil {
				s.DownloadErr = err
			} else {
				s.FilePath = filePath
			}
			s.Downloading = false
		}(song)
	}

	// If playback is not already running, start playing the queue.
	if !b.isPlaying {
		go b.startPlaying(msg)
	}
}

// startPlaying processes the music queue and plays songs sequentially.
func (b *BotController) startPlaying(msg *discordgo.MessageCreate) {
	if len(b.musicQueue) == 0 {
		b.isPlaying = false
		b.displayCmdError(msg.ChannelID, "🎵 Queue is empty.")
		return
	}

	b.isPlaying = true

	// Process the queue in FIFO order.
	for len(b.musicQueue) > 0 {
		// Dequeue the first song.
		song := b.musicQueue[0]
		b.musicQueue = b.musicQueue[1:]

		b.displayCmdError(msg.ChannelID, fmt.Sprintf("🎶 Now playing: %s", song.URL))

		// Wait until the song finishes downloading.
		for song.Downloading {
			log.Println("⏳ Waiting for download to finish for:", song.URL)
			time.Sleep(1 * time.Second)
		}

		// If there was an error during download, skip the song.
		if song.DownloadErr != nil || song.FilePath == "" {
			log.Println("❌ Error downloading song:", song.DownloadErr)
			b.displayCmdError(msg.ChannelID, fmt.Sprintf("⚠ Error downloading song: %s", song.URL))
			continue
		}

		// Ensure the file exists.
		if _, err := os.Stat(song.FilePath); os.IsNotExist(err) {
			log.Println("❌ Error: File does not exist!", song.FilePath)
			b.displayCmdError(msg.ChannelID, "⚠ Error: Downloaded file not found.")
			continue
		}

		// Join the voice channel.
		vc, err := b.joinUserChannel(msg.GuildID, msg.Author.ID, false, true)
		if err != nil {
			b.displayCmdError(msg.ChannelID, "⚠ Failed to join voice channel.")
			return
		}

		log.Println("✅ Bot joined voice channel. Starting playback...")
		time.Sleep(2 * time.Second) // Short delay before streaming.

		StreamAudio(vc, song.FilePath)
	}

	b.isPlaying = false
}

// StreamAudio streams the specified audio file to Discord.
func StreamAudio(vc *discordgo.VoiceConnection, filename string) {
	log.Println("🎵 Preparing to stream audio...")

	if vc == nil {
		log.Println("❌ Voice connection is nil. Cannot play audio.")
		return
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		log.Println("❌ Error getting absolute path:", err)
		return
	}

	log.Println("✅ Checking file:", absPath)
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		log.Println("❌ Error: Downloaded file does not exist!", absPath)
		return
	}

	log.Println("📂 File Size:", fileInfo.Size(), "bytes")
	log.Println("🔍 File Format:", filepath.Ext(absPath))
	if fileInfo.Size() < 1000 {
		log.Println("❌ File is too small! Probably an empty/corrupt MP3.")
		return
	}

	log.Println("✅ Voice connection established.")
	vc.Speaking(true)
	defer vc.Speaking(false)

	log.Println("🔄 Encoding file with DCA:", absPath)
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "audio"
	options.Volume = 10.0
	options.FrameRate = 48000
	options.BufferedFrames = 100

	encodeSession, err := dca.EncodeFile(absPath, options)
	if err != nil {
		log.Println("❌ Error encoding file with DCA:", err)
		return
	}
	defer encodeSession.Cleanup()

	log.Println("✅ Audio file encoded, starting playback...")
	frameCount := 0
	for {
		frame, err := encodeSession.OpusFrame()
		if err != nil {
			log.Println("🎵 Finished playing, disconnecting...")
			break
		}
		if len(frame) == 0 {
			log.Println("⚠ Warning: Empty audio frame, DCA might be broken.")
			continue
		}
		vc.OpusSend <- frame
		frameCount++
	}
	log.Println("✅ Total frames sent:", frameCount)
	vc.Disconnect()
}

// downloadYouTubeAudio downloads and converts audio using yt-dlp and FFmpeg.
func downloadYouTubeAudio(url string) (string, error) {
	projectRoot, err := filepath.Abs(".")
	if err != nil {
		log.Println("❌ Error getting project root:", err)
		return "", err
	}

	audioDir := filepath.Join(projectRoot, "audio")
	os.MkdirAll(audioDir, os.ModePerm)
	binDir := filepath.Join(projectRoot, "bin")
	ytDLPPath := filepath.Join(binDir, "yt-dlp.exe")
	ffmpegPath := filepath.Join(binDir, "ffmpeg.exe")
	ytdlpOutputTemplate := filepath.Join(audioDir, "%(title)s.%(ext)s")

	log.Println("📥 Downloading YouTube audio with yt-dlp...")
	cmd := exec.Command(ytDLPPath, "-f", "bestaudio", "-o", ytdlpOutputTemplate, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Println("❌ Error downloading YouTube audio:", err)
		return "", err
	}

	files, err := os.ReadDir(audioDir)
	if err != nil {
		log.Println("❌ Error reading audio directory:", err)
		return "", err
	}

	var downloadedFile string
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".webm" || ext == ".m4a" {
			downloadedFile = filepath.Join(audioDir, file.Name())
			break
		}
	}

	if downloadedFile == "" {
		log.Println("❌ No valid downloaded file found in:", audioDir)
		return "", fmt.Errorf("no valid downloaded file found")
	}
	log.Println("✅ Download complete:", downloadedFile)

	finalAudioFile := downloadedFile[:len(downloadedFile)-len(filepath.Ext(downloadedFile))] + ".mp3"
	log.Println("🎵 Converting to MP3 using FFmpeg...")
	ffmpegCmd := exec.Command(ffmpegPath, "-y", "-i", downloadedFile, "-q:a", "0", "-map", "a", finalAudioFile)
	ffmpegCmd.Stdout = os.Stdout
	ffmpegCmd.Stderr = os.Stderr
	err = ffmpegCmd.Run()
	if err != nil {
		log.Println("❌ Error converting audio with FFmpeg:", err)
		return "", err
	}

	if _, err := os.Stat(finalAudioFile); os.IsNotExist(err) {
		log.Println("❌ Conversion failed. MP3 file not found:", finalAudioFile)
		return "", fmt.Errorf("conversion failed: %s", finalAudioFile)
	}
	log.Println("✅ Conversion complete:", finalAudioFile)
	os.Remove(downloadedFile)

	return finalAudioFile, nil
}
