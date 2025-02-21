#!/bin/bash

# User ID variable
USER_ID="123"

# API base URL
BASE_URL="http://localhost:8080/audio/user/$USER_ID/phrase"

# Directories
TESTDATA_DIR="testdata"
DOWNLOAD_DIR="_downloaded"

# Create the _downloaded directory if it doesn't exist
mkdir -p "$DOWNLOAD_DIR"

# Iterate over all .wav files in the testdata directory
for FILE in "$TESTDATA_DIR"/*.wav; do
    # Extract the filename without extension
    FILENAME=$(basename -- "$FILE")
    PHRASE_ID="${FILENAME%.*}"  # Removes the extension

    echo "Uploading $FILE as phrase $PHRASE_ID..."

    # Upload the file
    curl -X POST "$BASE_URL/$PHRASE_ID" \
        -F "audio=@$FILE"

    echo -e "\nUploaded $FILE"
done

# Wait for upload completion
sleep 2

# Download the files in m4a and wav format
for FILE in "$TESTDATA_DIR"/*.wav; do
    FILENAME=$(basename -- "$FILE")
    PHRASE_ID="${FILENAME%.*}"  # Removes the extension

    echo "Downloading $PHRASE_ID in m4a format..."
    curl -o "$DOWNLOAD_DIR/$PHRASE_ID.m4a" "$BASE_URL/$PHRASE_ID/m4a"

    echo "Downloading $PHRASE_ID in wav format..."
    curl -o "$DOWNLOAD_DIR/$PHRASE_ID.wav" "$BASE_URL/$PHRASE_ID/wav"
done

echo "All files uploaded and downloaded successfully!"
