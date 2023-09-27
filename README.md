# GPT-Whisper-Captions

## Description

This repository contains a Go module for video transcription using OpenAI's Whisper API. It allows you to extract audio from a video file, send it to the Whisper API for transcription, and then embed the transcribed text back into the video as subtitles.

## Installation

1. Make sure you have Go installed on your system. If not, you can download and install it from [here](https://golang.org/dl/).

2. Make sure FFmpeg is installed on your system. If not, you can download it from [here](https://ffmpeg.org/download.html).

As a library
```bash
go get github.com/stevenlawton/GPT-Whisper-captions
```

## Usage
In your Go app you can do something like

```go
package main

import (
    "log"
    "os"

    captions "github.com/stevenlawton/GPT-Whisper-captions"
)

func main() {
   // Check if FFmpeg is installed
   err = captions.CheckFFmpegInstallation()
   if err != nil {
      log.Fatal("FFmpeg is not installed: ", err)
      return
   }

}
```
## Features

- Audio extraction from video using FFmpeg
- Audio segmentation for more manageable transcription
- Transcription using OpenAI's Whisper API
- SRT file generation for subtitles
- Embedding subtitles back into the video

## Contributing

Feel free to open issues or submit pull requests. Your contributions are welcome!

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

