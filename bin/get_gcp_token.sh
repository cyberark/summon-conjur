#!/bin/bash
set -xeo pipefail

HOST_ID="$1"
ACCOUNT="$2"
OUTPUT_DIR="$3"

BASE_URL="http://metadata.google.internal/computeMetadata/v1"
IDENTITY_URL="$BASE_URL/instance/service-accounts/default/identity"
PROJECT_ID_URL="$BASE_URL/project/project-id"
METADATA_FLAVOR_HEADER="Metadata-Flavor: Google"

# Check if account, hostId, and output file are provided
if [[ -z "$ACCOUNT" || -z "$HOST_ID" || -z "$OUTPUT_DIR" ]]; then
  echo "Usage: $0 <account> <hostId> <outputFile>"
  exit 1
fi

rm -rf "$OUTPUT_DIR" 2>/dev/null
mkdir -p "$OUTPUT_DIR"

# Build audience parameter
AUDIENCE="conjur/$ACCOUNT/$HOST_ID"

# Make the request to the metadata server
TOKEN=$(curl -s "$IDENTITY_URL?audience=$AUDIENCE&format=full" -H "$METADATA_FLAVOR_HEADER")

# Check if the request was successful
if [[ $? -ne 0 || -z "$TOKEN" ]]; then
  echo "Failed to fetch the token."
  exit 1
fi

# Store the token in a file
echo "$TOKEN" > "$OUTPUT_DIR/token"
echo "Token saved to $OUTPUT_DIR/token"

# Store the project ID in a file
GCP_PROJECT=$(curl -s "$PROJECT_ID_URL" -H "$METADATA_FLAVOR_HEADER")
echo "$GCP_PROJECT" > "$OUTPUT_DIR/project-id"
echo "Project ID saved to $OUTPUT_DIR/project-id"