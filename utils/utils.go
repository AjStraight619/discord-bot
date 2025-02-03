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

// func FindExtractedLibOpusDir(basePath string) string {
// 	files, err := os.ReadDir(basePath)
// 	if err != nil {
// 		return ""
// 	}

// 	for _, f := range files {
// 		if f.IsDir() && strings.HasPrefix(strings.ToLower(f.Name()), "libopus") {
// 			return filepath.Join(basePath, f.Name())
// 		}
// 	}
// 	return ""
// }

// func EnsureLibOpus(binDir string) {
// 	// Adjust the DLL filename as needed.
// 	dllPath := filepath.Join(binDir, "libopus-0.dll")
// 	if FileExists(dllPath) {
// 		log.Println("✅ libopus-0.dll is already installed.")
// 		return
// 	}

// 	log.Println("⚠ libopus-0.dll not found, downloading...")

// 	// Replace this URL with the actual URL where you can download a precompiled libopus ZIP for Windows.
// 	libopusURL := "https://example.com/path/to/libopus-windows.zip"
// 	zipPath := filepath.Join(binDir, "libopus.zip")

// 	err := DownloadFile(zipPath, libopusURL)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to download libopus: %v", err)
// 	}

// 	log.Println("✅ libopus downloaded, extracting...")

// 	err = unzip(zipPath, binDir)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to extract libopus: %v", err)
// 	}

// 	// Find the extracted directory (this example looks for one that starts with "libopus").
// 	extractedDir := FindExtractedLibOpusDir(binDir)
// 	if extractedDir == "" {
// 		log.Fatal("❌ Failed to find extracted libopus directory.")
// 	}

// 	// Move the DLL from the extracted directory into binDir.
// 	err = os.Rename(filepath.Join(extractedDir, "libopus-0.dll"), dllPath)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to move libopus DLL: %v", err)
// 	}

// 	// Clean up: remove the zip file and the extracted directory.
// 	os.Remove(zipPath)
// 	os.RemoveAll(extractedDir)

// 	log.Println("✅ libopus-0.dll installed successfully in", binDir)
// }
