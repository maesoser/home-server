#!/bin/sh

HUGO_SOURCE="${HUGO_SOURCE:=/output}"
echo "HUGO_THEME:" $HUGO_THEME
echo "HUGO_BASEURL" $HUGO_BASEURL
echo "HUGO_SOURCE" $HUGO_SOURCE

# HUGO=/usr/local/sbin/hugo
# echo "Hugo path: $HUGO"

hugo server \
  --minify \
  --gc \
  --disableFastRender \
  --watch=true \
  --appendPort=false \
  --theme="$HUGO_THEME" \
  --source="$HUGO_SOURCE" \
  --baseURL="$HUGO_BASEURL" \
  --bind="0.0.0.0"
