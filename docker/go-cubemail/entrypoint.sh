#!/bin/sh
set -e
./go-cubemail migrate
exec ./go-cubemail serve
