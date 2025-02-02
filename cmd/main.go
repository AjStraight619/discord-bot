package main

import (
	"archive/zip"
	"discord-bot/internal/messages"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	rootDir := getProjectRoot()
	binDir := filepath.Join(rootDir, "bin")
	os.MkdirAll(binDir, os.ModePerm)

	ensureYTDLP(binDir)
	ensureFFmpeg(binDir)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token := os.Getenv("DISCORD_KEY")
	NewsAPI := os.Getenv("NEWS_KEY")
	OpenAIAPI := os.Getenv("OPENAI_KEY")

	if Token == "" || NewsAPI == "" || OpenAIAPI == "" {
		log.Fatal("Missing environment variables.")
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	newsClient := &messages.NewsClient{APIKey: NewsAPI}
	aiClient := messages.NewAIClient(OpenAIAPI)

	msgHandler := &messages.Messages{
		Session:    dg,
		NewsClient: newsClient,
		AIClient:   aiClient,
	}

	dg.AddHandler(msgHandler.MessageHandler)

	// Open a connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error connecting to Discord: %v", err)
	}

	fmt.Println("Bot is now running! Press CTRL+C to exit.")

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Cleanup
	dg.Close()
}

func getProjectRoot() string {
	rootDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("❌ Error getting project root: %v", err)
	}
	return rootDir
}

func ensureFFmpeg(binDir string) {
	ffmpegPath := filepath.Join(binDir, "ffmpeg.exe")
	ffprobePath := filepath.Join(binDir, "ffprobe.exe")

	if fileExists(ffmpegPath) && fileExists(ffprobePath) {
		log.Println("✅ FFmpeg and FFprobe are already installed.")
		return
	}

	log.Println("⚠ FFmpeg or FFprobe not found, downloading...")

	ffmpegURL := "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
	zipPath := filepath.Join(binDir, "ffmpeg.zip")

	err := downloadFile(zipPath, ffmpegURL)
	if err != nil {
		log.Fatalf("❌ Failed to download FFmpeg: %v", err)
	}

	log.Println("✅ FFmpeg downloaded, extracting...")

	err = unzip(zipPath, binDir)
	if err != nil {
		log.Fatalf("❌ Failed to extract FFmpeg: %v", err)
	}

	extractedDir := findExtractedFFmpegDir(binDir)
	if extractedDir == "" {
		log.Fatal("❌ Failed to find extracted FFmpeg directory.")
	}

	os.Rename(filepath.Join(extractedDir, "ffmpeg.exe"), ffmpegPath)
	os.Rename(filepath.Join(extractedDir, "ffprobe.exe"), ffprobePath)

	os.Remove(zipPath)
	os.RemoveAll(extractedDir)

	log.Println("✅ FFmpeg and FFprobe installed successfully in ./bin/")
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(fPath)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func findExtractedFFmpegDir(basePath string) string {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return ""
	}

	for _, f := range files {
		if f.IsDir() && strings.HasPrefix(f.Name(), "ffmpeg") {
			return filepath.Join(basePath, f.Name(), "bin")
		}
	}
	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ensureYTDLP(binDir string) {
	ytDLPPath := filepath.Join(binDir, "yt-dlp.exe")

	if fileExists(ytDLPPath) {
		log.Println("✅ yt-dlp is already installed.")
		return
	}

	log.Println("⚠ yt-dlp not found, downloading...")

	cmd := exec.Command("curl", "-L", "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe", "-o", ytDLPPath)
	if err := cmd.Run(); err != nil {
		log.Fatal("❌ Failed to download yt-dlp:", err)
	}

	log.Println("✅ yt-dlp installed successfully in ./bin/")
}
