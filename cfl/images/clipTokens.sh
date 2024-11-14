#!/bin/bash

function clip_token() {
  # Define variables for center coordinates (x, y) and radius
  COLOR="${1}"
  CENTER_X="${2}"     # Replace with the actual X coordinate of the token center
  CENTER_Y="${3}"     # Replace with the actual Y coordinate of the token center
  RADIUS=82        # Replace with the actual radius of the token
  DIAMETER=$(( 2 * RADIUS ))

  # Define the input and output paths
  INPUT_IMAGE="debug_card.png"
  TOKEN_IMAGE="token_${COLOR}.png"
  GREY_IMAGE="token_${COLOR}_grey.png"

  # Perform the crop, resize, and mask application using the calculated DIAMETER
  magick "$INPUT_IMAGE" \
    -crop "${DIAMETER}x${DIAMETER}+$((CENTER_X - RADIUS))+$((CENTER_Y - RADIUS))" +repage \
    -resize "${DIAMETER}x${DIAMETER}" \
    \( +clone -alpha transparent -fill white -draw "circle $RADIUS,$RADIUS $DIAMETER,$RADIUS" \) \
    -compose CopyOpacity -composite "$TOKEN_IMAGE"

  echo "created '${TOKEN_IMAGE}'"

  magick "${TOKEN_IMAGE}" -modulate 70,0,100 "${GREY_IMAGE}"

  echo "created '${GREY_IMAGE}'"
}

# 376 398

############################################################################
clip_token red    375 129

clip_token black  375 392

clip_token green  204 896

clip_token yellow 552 896
