# audio-storage

audio-storage is a backend service for uploading, storing, converting, and retrieving audio files. It utilizes MinIO (S3-compatible storage) to store audio files and PostgreSQL to manage metadata. Files are converted to `.m4a` (AAC) upon upload and can be retrieved in different formats.

## Features

- Upload audio files (supports `.wav`, `.mp3`, `.flac`, etc.).
- Automatic conversion to `.m4a` (AAC) on upload.
- Store files in MinIO (S3-compatible storage).
- Retrieve files in multiple formats (`mp3`, `wav`, `flac`, etc.).
- Supports chunked uploads and downloads.
- Resumable downloads via HTTP Range Requests.
- PostgreSQL integration to store file metadata.
- Performance testing using K6.

## Prerequisites

Ensure you have:
- Docker & Docker Compose installed.
- Make installed.
- K6 (for performance testing, or run via Docker).

## Running the Service

Copy the env sample and make adjustment if needed:
```
cp env.sample .env
```

Start the entire stack using Docker Compose:
```sh
make docker.up
```

This will:
- Start the API server (audio-storage).
- Start PostgreSQL (postgres).
- Start MinIO (minio).

## Stopping the Service

To stop the running containers:

```sh
make docker.down
```

## API Endpoints

### Upload an Audio File

#### Endpoint

```
POST /audio/user/{user_id}/phrase/{phrase_id}
```
#### Request Body

- multipart/form-data: Upload an audio file.

#### Path Parameters

- user_id: user id of users owned the audio
- phrase_id: a phrase corresponding with the audio

#### Example Request

```
curl -X POST -F "audio=@sample.wav" http://localhost:8080/audio/user/123/phrase/456
```

#### Example Response

```
{
  "message": "Upload successful",
  "filename": "123_456.m4a"
}
```

### Download an Audio File in a Specific Format

#### Endpoint

```
GET /audio/user/{user_id}/phrase/{phrase_id}/{audio_format}
```

#### Path Parameters

- user_id: user id of users owned the audio
- phrase_id: a phrase corresponding with the audio
- audio_format: audio format of downloaded file. Tested values are: mp3, wav, and flac

#### Example Request

```
curl -O http://localhost:8080/audio/user/123/phrase/456/mp3
```

#### Example Response

- Returns the converted audio file.

## Testing the Service

### Running K6 Performance Tests

To run K6 performance tests:

```
make k6-test
```
This will:
- Upload a random audio file from testdata/.
- Download the file in a random format.
- Perform negative test cases (invalid formats, missing files).

## Database Access

To access PostgreSQL served by Docker Compose:

```
make postgres.login
```