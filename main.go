package GPTWhisperCaptions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

type TranscriptionResponse struct {
	Text string `json:"text"`
}

type TimedText struct {
	Start float64 // start time in seconds
	End   float64 // end time in seconds
	Text  string  // the transcribed text
}

// CheckFFmpegInstallation checks if FFmpeg is installed on the system
func CheckFFmpegInstallation() error {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return errors.New("FFmpeg is not installed on this system")
	}
	return nil
}

// ExtractAudio Extracts audio from a video file using FFmpeg
func ExtractAudio(videoFilename string, audioFilename string) error {
	cmd := exec.Command("ffmpeg", "-i", videoFilename, "-q:a", "0", "-map", "a", audioFilename, "-y")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func SegmentAudio(audioFilename string, segmentLength int) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", audioFilename, // Input file
		"-f", "segment", // Enable segmenting
		"-segment_time", fmt.Sprintf("%d", segmentLength), // Segment length in seconds
		"-c", "copy", // Copy streams without re-encoding
		"segment_%03d.mp3", // Output filename pattern
	)
	return cmd.Run()
}

func EmbedSubtitles(videoFilename string, srtFilename string, outputVideoFilename string) error {
	//subtitleOptions := "FontName=Arial:FontSize=16:Fontcolor=white:Box=1:Boxcolor=black@0.5"
	subtitleOptions := "Fontname=Consolas,BackColour=&H80000000,Spacing=0.2,Outline=0,Shadow=0.75"

	cmd := exec.Command(
		"ffmpeg",
		"-y",                // Force overwrite output file without asking
		"-i", videoFilename, // Input video file
		"-vf", fmt.Sprintf("subtitles=%s:force_style='%s'", srtFilename, subtitleOptions),
		"-c:a", "copy", // Copy audio stream without re-encoding
		"-c:v", "libx264", // Video codec to use
		"-crf", "23", // Constant rate factor for video quality
		"-preset", "fast", // Preset for speed/quality trade-off
		outputVideoFilename, // Output video file
	)

	fmt.Println(cmd.String())

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Capture the standard error output from the command
		return fmt.Errorf("ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// SendToWhisper Sends the audio file to OpenAI's Whisper API
func SendToWhisper(audioFilename string, apiKey string) (string, error) {
	// Create a buffer to hold the form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Open the audio file
	f, err := os.Open(audioFilename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create a form field and write the file into that field
	fw, err := w.CreateFormFile("file", audioFilename)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return "", err
	}

	// Add other fields
	w.WriteField("model", "whisper-1")

	// Close the writer
	w.Close()

	// Create a request with the form data
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &b)
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("openai whisper API returned non-200 status code: %s", respBytes)
	}

	// Unmarshal JSON into a Go struct
	var transcriptionResponse TranscriptionResponse
	if err := json.Unmarshal(respBytes, &transcriptionResponse); err != nil {
		return "", err
	}

	return transcriptionResponse.Text, nil
}

func GenerateSRT(timedTexts []TimedText, srtFilename string) error {
	var srtContent string
	for i, tt := range timedTexts {
		start := fmt.Sprintf("%02d:%02d:%02d,000",
			int(tt.Start)/3600, int(tt.Start)/60%60, int(tt.Start)%60)
		end := fmt.Sprintf("%02d:%02d:%02d,000",
			int(tt.End)/3600, int(tt.End)/60%60, int(tt.End)%60)
		srtContent += fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, start, end, tt.Text)
	}

	return ioutil.WriteFile(srtFilename, []byte(srtContent), 0644)
}
