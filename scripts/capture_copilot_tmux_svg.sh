#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SESSION_NAME=${SESSION_NAME:-go-te-copilot}
SOCKET_NAME=${SOCKET_NAME:-go-te-copilot}
TERMINAL_WIDTH=${TERMINAL_WIDTH:-120}
TERMINAL_HEIGHT=${TERMINAL_HEIGHT:-40}
CAPTURE_DELAY=${CAPTURE_DELAY:-3}
OUTPUT_SVG=${OUTPUT_SVG:-"${ROOT}/copilot.svg"}
CAPTURE_FILE=${CAPTURE_FILE:-"${ROOT}/copilot.capture"}
COPILOT_CMD=${COPILOT_CMD:-"copilot"}
ATTACH_SESSION=${ATTACH_SESSION:-0}
INCLUDE_STATUS=${INCLUDE_STATUS:-1}
STATUS_BG=${STATUS_BG:-green}
STATUS_FG=${STATUS_FG:-black}
STATUS_TEXT=${STATUS_TEXT:-" ${SESSION_NAME} "}

color_to_sgr() {
  local color="$1"
  local kind="$2"
  local base=30
  if [[ "$kind" == "bg" ]]; then
    base=40
  fi
  case "$color" in
    black)
      echo "$base"
      ;;
    red)
      echo "$((base + 1))"
      ;;
    green)
      echo "$((base + 2))"
      ;;
    yellow|brown)
      echo "$((base + 3))"
      ;;
    blue)
      echo "$((base + 4))"
      ;;
    magenta)
      echo "$((base + 5))"
      ;;
    cyan)
      echo "$((base + 6))"
      ;;
    white)
      echo "$((base + 7))"
      ;;
    default|*)
      echo "$((base + 9))"
      ;;
  esac
}

fg_code=$(color_to_sgr "$STATUS_FG" "fg")
bg_code=$(color_to_sgr "$STATUS_BG" "bg")

for cmd in tmux go; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "error: missing required command '$cmd'" >&2
    exit 1
  fi
done

if ! command -v gh >/dev/null 2>&1 && ! command -v copilot >/dev/null 2>&1; then
  echo "warning: Copilot CLI not found. Set COPILOT_CMD to a valid command." >&2
fi

tmux -L "$SOCKET_NAME" new-session -d -x "$TERMINAL_WIDTH" -y "$TERMINAL_HEIGHT" -s "$SESSION_NAME" "TERM=xterm-256color ${COPILOT_CMD}"
tmux -L "$SOCKET_NAME" set-option -t "$SESSION_NAME" status on
tmux -L "$SOCKET_NAME" set-option -t "$SESSION_NAME" status-position bottom
tmux -L "$SOCKET_NAME" set-option -t "$SESSION_NAME" status-style "bg=${STATUS_BG},fg=${STATUS_FG}"
tmux -L "$SOCKET_NAME" set-option -t "$SESSION_NAME" status-format[0] "${STATUS_TEXT}"

if [[ "$ATTACH_SESSION" == "1" ]]; then
  echo "Attach to the tmux session and run Copilot. Detach (Ctrl-b d) when ready to capture." >&2
  tmux -L "$SOCKET_NAME" attach -t "$SESSION_NAME"
else
  sleep "$CAPTURE_DELAY"
fi

PANE_CAPTURE="${CAPTURE_FILE}.pane"

tmux -L "$SOCKET_NAME" capture-pane -p -e -t "${SESSION_NAME}:0.0" > "$PANE_CAPTURE"

pane_width=$(tmux -L "$SOCKET_NAME" display-message -p -t "${SESSION_NAME}:0.0" "#{pane_width}")
pane_height=$(tmux -L "$SOCKET_NAME" display-message -p -t "${SESSION_NAME}:0.0" "#{pane_height}")
render_height=$pane_height

if [[ "$INCLUDE_STATUS" == "1" ]]; then
  status_text=$(tmux -L "$SOCKET_NAME" display-message -p -F "$STATUS_TEXT")
  status_line=$(printf "%-*s" "$pane_width" "$status_text")
  status_line=${status_line:0:$pane_width}
  cp "$PANE_CAPTURE" "$CAPTURE_FILE"
  if [[ -s "$CAPTURE_FILE" ]]; then
    last_char=$(tail -c 1 "$CAPTURE_FILE" || true)
    if [[ "$last_char" != $'\n' ]]; then
      printf "\n" >> "$CAPTURE_FILE"
    fi
  fi
  printf "\033[0m\033[%s;%sm%s\033[0m" "$fg_code" "$bg_code" "$status_line" >> "$CAPTURE_FILE"
  render_height=$((pane_height + 1))
else
  cp "$PANE_CAPTURE" "$CAPTURE_FILE"
fi

rm -f "$PANE_CAPTURE"

tmux -L "$SOCKET_NAME" kill-session -t "$SESSION_NAME"

go run "${ROOT}/scripts/render_svg.go" "$CAPTURE_FILE" "$OUTPUT_SVG" "$pane_width" "$render_height"

echo "Wrote SVG to ${OUTPUT_SVG}" >&2
