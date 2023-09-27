package GPTWhisperCaptions

import (
	"testing"
)

func TestExtractAudio(t *testing.T) {
	err := ExtractAudio("test_video.mp4", "test_audio.mp3")
	if err != nil {
		t.Errorf("Failed to extract audio: %s", err)
	}
}

func TestSegmentAudio(t *testing.T) {
	err := SegmentAudio("test_audio.mp3", 5)
	if err != nil {
		t.Errorf("Failed to segment audio: %s", err)
	}
}

func TestEmbedSubtitles(t *testing.T) {
	err := EmbedSubtitles("test_video.mp4", "test.srt", "test_output.mp4")
	if err != nil {
		t.Errorf("Failed to embed subtitles: %s", err)
	}
}

func TestGenerateSRT(t *testing.T) {
	timedTexts := []TimedText{
		{Start: 0, End: 5, Text: "Hello"},
		{Start: 6, End: 10, Text: "World"},
	}

	err := GenerateSRT(timedTexts, "test.srt")
	if err != nil {
		t.Errorf("Failed to generate SRT: %s", err)
	}
}

func TestSendToWhisper(t *testing.T) {
	apiKey := ""

	text, err := SendToWhisper("test_segment.mp3", apiKey)
	if err != nil {
		t.Errorf("Failed to send to Whisper: %s", err)
	}
	if text == "" {
		t.Errorf("Received empty text from Whisper")
	}
}
