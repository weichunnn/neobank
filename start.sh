#!/bin/sh

set -e # exit immediately if non zero status returned

echo "start app"
exec "$@"
