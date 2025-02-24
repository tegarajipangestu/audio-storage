#!/bin/bash

USER_ID="123"
BASE_URL="http://localhost:8080/audio/user/$USER_ID/phrase"
TESTDATA_DIR="e2e/testdata"
DOWNLOAD_DIR="e2e/_downloaded"

mkdir -p "$DOWNLOAD_DIR"

for FILE in "$TESTDATA_DIR"/*.wav; do
    FILENAME=$(basename -- "$FILE")
    PHRASE_ID="${FILENAME%.*}"

    echo "Uploading $FILE as phrase $PHRASE_ID..."

    curl -X POST "$BASE_URL/$PHRASE_ID" \
        -F "audio=@$FILE"

    echo -e "\nUploaded $FILE"
done

sleep 2

for FILE in "$TESTDATA_DIR"/*.wav; do
    FILENAME=$(basename -- "$FILE")
    PHRASE_ID="${FILENAME%.*}"

    echo "Downloading $PHRASE_ID in m4a format..."
    curl -o "$DOWNLOAD_DIR/$PHRASE_ID.m4a" "$BASE_URL/$PHRASE_ID/m4a"

    echo "Downloading $PHRASE_ID in wav format..."
    curl -o "$DOWNLOAD_DIR/$PHRASE_ID.wav" "$BASE_URL/$PHRASE_ID/wav"
done

echo "All files uploaded and downloaded successfully! Check at $DOWNLOAD_DIR folder"