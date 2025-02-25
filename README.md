# audio-storage

audio-storage is a backend service for uploading, storing, converting, and retrieving audio files. It utilizes MinIO (S3-compatible storage) to store audio files and PostgreSQL to manage metadata. Files are converted to `.m4a` upon upload and can be retrieved in different formats.

## Features

- Upload audio files (supports `.wav` and `.m4a`)
- Automatic conversion to `.m4a` (AAC) on upload
- Store files in MinIO (S3-compatible storage)
- Retrieve files in multiple formats (supports `.wav` and `.m4a`)
- PostgreSQL integration to store file metadata
- Performance testing using K6

## Prerequisites

Ensure you have:
- Docker & Docker Compose installed
- Make installed
- K6 (for performance testing, or run via Docker)

## Running the Service

Copy the env sample and make adjustment if needed:
```
cp env.sample .env
```

Start the entire stack using Docker Compose:
```sh
make docker-compose.up
```

This will:
- Build Docker Image of audio-storage locally
- Start the API server (audio-storage) using newly built image
- Start PostgreSQL (postgres)
- Start MinIO (minio)

## Database and MinIO Access

To access PostgreSQL served by Docker Compose:

```
make postgres.login
```

To access MinIO web UI, just open http://localhost:9001 in browser. For login, use the value of `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD` defined in docker-compose.yml file (
defaults to `minioadmin:minioadmin`)

## Stopping the Service

To stop the running containers:

```sh
make docker-compose.down
```

## API Endpoints

### Upload an Audio File

#### Endpoint

```
POST /audio/user/{user_id}/phrase/{phrase_id}
```
#### Request Body

- multipart/form-data: Upload an audio file

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
- audio_format: audio format of downloaded file. Tested values are: m4a and wav

#### Example Request

```
curl -O http://localhost:8080/audio/user/123/phrase/456/wav
```

#### Example Response

- Returns the converted audio file

## Testing the Service

### Running End to End Test using shell script

To execute e2e test using predefined shell script:

```
make test
```

This will:
- Upload all audio file from e2e/testdata folder
- Download all files in a m4a and wav format in e2e/_downloaded folder


### Running End to End Test using K6

This repository also leverages [k6](https://k6.io/) to run e2e test which later can be converted to performance test script. To execute e2e test:

```
make k6-test
```

This will:
- Upload a random audio file from e2e/testdata folder
- Download the file in a random format
- Perform negative test cases (invalid formats, missing files)

## Future Improvements
### Cloud Storage (S3)
- Move from self-hosted MinIO to Amazon S3 / Google Cloud Storage for high availability.
- Implement S3 lifecycle policies to auto-delete old files.
### Asynchronous Processing
- Use a message queue (Kafka/RabbitMQ) for background processing. Upload -> Queue job -> Worker converts audio file -> Stores in MinIO or Amazon S3 / Google Cloud Storage.
### Auto-Scaling
- Deploy in Kubernetes (K8s) with Horizontal Pod Autoscaler (HPA) for high-traffic handling.
### Edge Caching for Faster Downloads
- Integrate Cloudflare / AWS CloudFront / Google's Cloud CDN as a CDN layer.
- Cache frequently downloaded files for faster delivery.
### Enable Resumable Uploads
- Implement multipart uploads (via TUS protocol or S3 multipart upload) to allow large files to resume uploading after failure.
- Parallel chunk uploads for speed improvements.
### Add Support for Chunked Downloads
- Implement supports for `Range` request header and making sure the storage supports partial object retrieval
- Send appropriate `Content-Range` and `Accept-Ranges` response headers so the client understands that partial downloads are supported.