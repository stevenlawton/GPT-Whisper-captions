# Video Subtitler with Whisper API

## Overview

This is a Go-based application that automates the process of generating subtitles for a given video file. It leverages OpenAI's Whisper ASR API for speech recognition.

## Features

- Extracts audio from video files.
- Segments the audio for more manageable processing.
- Transcribes the audio segments using OpenAI's Whisper ASR API.
- Generates an SRT subtitle file from the transcriptions.
- Embeds the SRT subtitles back into the video.

## Requirements

- Go 1.x.x
- FFmpeg installed and available in your PATH
- An OpenAI API key for Whisper

## Installation

1. Clone this repository:

    ```bash
    git clone https://github.com/stevenlawton/GPT-Whisper-captions.git
    ```

2. Navigate to the project directory:

    ```bash
    cd GPT-Whisper-captions
    ```

3. Install dependencies (if any):

    ```bash
    go get -u ./...
    ```

## Usage

### Configuration

Copy `.env.example` to `.env` and add your OpenAI API key.

```bash
cp .env.example .env
