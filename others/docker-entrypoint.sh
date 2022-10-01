#!/bin/sh

# Abort on any error (including if wait-for-it fails).
set -e

# Wait for the backend to be up, if we know where it is.
if [ -n "$NEO4J_HOST" ]; then
  ./wait-for-it.sh "$NEO4J_HOST:${NEO4J_PORT:-6000}"
fi

# Run the main container command.
exec "$@"