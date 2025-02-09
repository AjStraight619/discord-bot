package utils

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetProjectRoot() string {
	rootDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("❌ Error getting project root: %v", err)
	}
	return rootDir
}

func EnsureFFmpeg(binDir string) {
	ffmpegPath := filepath.Join(binDir, "ffmpeg.exe")
	ffprobePath := filepath.Join(binDir, "ffprobe.exe")

	if FileExists(ffmpegPath) && FileExists(ffprobePath) {
		log.Println("✅ FFmpeg and FFprobe are already installed.")
		return
	}

	log.Println("⚠ FFmpeg or FFprobe not found, downloading...")

	ffmpegURL := "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
	zipPath := filepath.Join(binDir, "ffmpeg.zip")

	err := DownloadFile(zipPath, ffmpegURL)
	if err != nil {
		log.Fatalf("❌ Failed to download FFmpeg: %v", err)
	}

	log.Println("✅ FFmpeg downloaded, extracting...")

	err = unzip(zipPath, binDir)
	if err != nil {
		log.Fatalf("❌ Failed to extract FFmpeg: %v", err)
	}

	extractedDir := FindExtractedFFmpegDir(binDir)
	if extractedDir == "" {
		log.Fatal("❌ Failed to find extracted FFmpeg directory.")
	}

	os.Rename(filepath.Join(extractedDir, "ffmpeg.exe"), ffmpegPath)
	os.Rename(filepath.Join(extractedDir, "ffprobe.exe"), ffprobePath)

	os.Remove(zipPath)
	os.RemoveAll(extractedDir)

	log.Println("✅ FFmpeg and FFprobe installed successfully in ./bin/")
}

func DownloadFile(filepath string, url string) error {
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

func FindExtractedFFmpegDir(basePath string) string {
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

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func EnsureYTDLP(binDir string) {
	ytDLPPath := filepath.Join(binDir, "yt-dlp.exe")

	if FileExists(ytDLPPath) {
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
