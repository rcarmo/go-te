# Pyte Port Test Coverage

This document summarizes the features covered by the Go port of the **pyte** test suite and the related emulator behavior we validate. It is intended as a high-level map of what the tests exercise, aligned with the porting rules in `PYTE_PORTING_CONVENTION.md`.

## Screen Core Behavior
- **Screen reset and defaults**
  - Full reset, soft reset, default attributes, and mode defaults.
  - Cursor state restoration and reset.
- **Cursor movement and positioning**
  - Absolute and relative cursor moves (CUU/CUD/CUF/CUB, CUP/HVP, CHA/HPA, VPA, VPR, HPR, CNL/CPL).
  - Carriage return, line feed, next line, index, reverse index, and tab movement.
  - Save/restore cursor and cursor visibility.
- **Margins and scrolling regions**
  - Top/bottom margins (DECSTBM) and left/right margins behavior.
  - Scrolling via index/reverse index and explicit scroll commands.
- **Insert/delete and erase operations**
  - Insert/delete characters and lines.
  - Erase in line/display, erase characters, and alignment display.
- **Line wrapping and autowrap**
  - Wrapping behavior for long runs, auto-wrap mode handling, and wrap state transitions.
- **Tab stops**
  - Setting and clearing tab stops, forward/backward tabbing.

## Drawing & Character Handling
- **ASCII and UTF-8 rendering**
  - Basic glyph drawing, wide characters, combining marks, and zero-width behavior.
- **Character sets and shifts**
  - G0/G1 character set selection and shift in/out handling.
  - Charset definition and fallback behavior.
- **CP437 and VT100 line drawing**
  - Character translation for legacy mappings.
- **Emoji and multi-byte glyphs**
  - Width-aware rendering and cursor placement semantics.

## SGR / Attributes
- **Standard SGR attributes**
  - Bold, italics, underline, blink, reverse, conceal, and strikethrough.
- **Color handling**
  - 8/16/256/truecolor foreground and background settings.
  - Reset semantics for foreground/background and attribute reset.
- **Private SGR handling**
  - Ignoring unsupported or private SGR sequences.

## Modes & Device Reports
- **Mode set/reset**
  - CSI mode set/reset (SM/RM) with private mode handling.
- **Device reports**
  - Device attributes (DA) and device status (DSR) responses.

## Stream Parsing
- **Escape sequence dispatch**
  - Basic, CSI, OSC, and charset escape sequence parsing.
  - Unknown sequence dispatch to debug hooks.
- **Parameter handling**
  - Missing parameters, overflow capping, and control character interruption.
- **Byte vs. rune stream behavior**
  - Byte-stream UTF-8 decoding and C1 control handling.
- **Debug stream output**
  - JSON payloads for debug streams matching pyte semantics.
- **Compatibility API**
  - Screen attach/detach and stream compatibility behavior.

## Diff Screen
- **Dirty tracking**
  - Line dirty state marking, full screen updates, and updates on wrap.
- **Index/reverse index effects**
  - Dirty tracking through scrolling operations.

## History Screen
- **Scrollback buffer**
  - History accumulation, page navigation, and bounds behavior.
- **Width changes**
  - History handling on resize and width changes.
- **Drawing with history**
  - Interaction between drawing and history preservation.

## Input/Output Semantics
- **Title/icon updates**
  - OSC updates for icon and window titles.
- **Selection and clipboard sequences**
  - Selection data set/query sequences via OSC 52.

## Known Skips
- **`test_draw_width2_irm`**
  - Marked xfail upstream in pyte; skipped in the Go port to preserve fidelity.
