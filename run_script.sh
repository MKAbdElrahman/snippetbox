#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
  set -o allexport
  source .env
  set +o allexport
  echo "Environment variables loaded from .env file."
else
  echo "Warning: .env file not found. Make sure it exists in the script's directory."
fi

echo "Environment variables exported to the shell."

# Execute the main binary
./tmp/main
