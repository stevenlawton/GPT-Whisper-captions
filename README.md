# GPT-Whisper-Captions

## Description

This repository contains a Go module for video transcription using OpenAI's Whisper API. It allows you to extract audio from a video file, send it to the Whisper API for transcription, and then embed the transcribed text back into the video as subtitles.

## Installation

1. Make sure you have Go installed on your system. If not, you can download and install it from [here](https://golang.org/dl/).

2. Clone this repository:
    ```bash
    git clone https://github.com/stevenlawton/GPT-Whisper-captions.git
    ```

3. Navigate to the project folder and install the dependencies:
    ```bash
    cd GPT-Whisper-captions
    go mod download
    ```

4. Make sure FFmpeg is installed on your system. If not, you can download it from [here](https://ffmpeg.org/download.html).

## Usage


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

