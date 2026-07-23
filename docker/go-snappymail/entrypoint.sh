#!/bin/sh
set -e
./go-snappymail migrate
exec ./go-snappymail serve
