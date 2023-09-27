package main

import (
	"bytes"
	"encoding/json"
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

// Extracts audio from a video file using FFmpeg
func extractAudio(videoFilename string, audioFilename string) error {
	cmd := exec.Command("ffmpeg", "-i", videoFilename, "-q:a", "0", "-map", "a", audioFilename, "-y")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func segmentAudio(audioFilename string, segmentLength int) error {
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

func embedSubtitles(videoFilename string, srtFilename string, outputVideoFilename string) error {
	subtitleOptions := "FontName=Arial:FontSize=16:Fontcolor=white:Box=1:Boxcolor=black@0.5"

	cmd := exec.Command(
		"ffmpeg",
		"-i", videoFilename, // Input video file
		"-vf", fmt.Sprintf("subtitles=%s:force_style='%s'", srtFilename, subtitleOptions), // Input subtitle file with styling
		"-c:a", "copy", // Copy audio stream without re-encoding
		"-c:v", "libx264", // Video codec to use
		"-crf", "23", // Constant rate factor for video quality
		"-preset", "fast", // Preset for speed/quality trade-off
		outputVideoFilename, // Output video file
	)
	return cmd.Run()
}

// Sends the audio file to OpenAI's Whisper API
func sendToWhisper(audioFilename string, apiKey string) (string, error) {
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
		return "", fmt.Errorf("Whisper API returned non-200 status code: %s", respBytes)
	}

	// Unmarshal JSON into a Go struct
	var transcriptionResponse TranscriptionResponse
	if err := json.Unmarshal(respBytes, &transcriptionResponse); err != nil {
		return "", err
	}

	return transcriptionResponse.Text, nil
}

func generateSRT(timedTexts []TimedText, srtFilename string) error {
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

func main() {
	videoFilename := "input_video.mp4"
	audioFilename := "output_audio.mp3"
	srtFilename := "output.srt"
	outputVideoFilename := "video_with_subtitles.mp4"
	apiKey := "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"
	segmentLength := 5

	// Step 1: Extract Audio
	if err := extractAudio(videoFilename, audioFilename); err != nil {
		fmt.Println("Error extracting audio:", err)
		return
	}

	// Step 2: Segment Audio
	if err := segmentAudio(audioFilename, segmentLength); err != nil {
		fmt.Println("Error segmenting audio:", err)
		return
	}

	// Step 3: Transcribe Each Segment
	var timedTexts []TimedText
	for i := 0; ; i++ {
		segmentFilename := fmt.Sprintf("segment_%03d.mp3", i)
		if _, err := os.Stat(segmentFilename); os.IsNotExist(err) {
			break // No more segments
		}
		transcript, err := sendToWhisper(segmentFilename, apiKey)
		if err != nil {
			fmt.Println("Error sending to Whisper:", err)
			return
		}
		startTime := float64(i * segmentLength)
		endTime := startTime + float64(segmentLength)
		timedTexts = append(timedTexts, TimedText{Start: startTime, End: endTime, Text: transcript})
	}

	// Step 4: Generate SRT
	if err := generateSRT(timedTexts, srtFilename); err != nil {
		fmt.Println("Error generating SRT:", err)
		return
	}

	fmt.Println("Subtitles generated.")

	// Step 5: Embed Subtitles
	if err := embedSubtitles(videoFilename, srtFilename, outputVideoFilename); err != nil {
		fmt.Println("Error embedding subtitles:", err)
		return
	}

	fmt.Println("Subtitles embedded.")
}
