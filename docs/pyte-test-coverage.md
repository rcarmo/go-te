# Pyte Port Test Coverage

This document summarizes the behaviors covered by the Go port of the **pyte** test suite. It is meant as a high-level guide to what the port exercises and how it maps to emulator behavior. For exact test mappings, see `TEST_AUDIT_CHECKLIST.md` and `PYTE_PORTING_CONVENTION.md`.

## Screen & Cursor Behavior
- **Reset and defaults**: full reset, soft reset, default attributes, and cursor state.
- **Cursor movement**: CUU/CUD/CUF/CUB, CUP/HVP, CHA/HPA, VPA/VPR, CNL/CPL, HPR, CR/LF/NEL/IND/RI, tab movement.
- **Margins & scrolling regions**: DECSTBM, left/right margin handling, scrolling via index/reverse index and explicit scroll commands.
- **Insert/delete/erase operations**: insert/delete characters and lines, erase in line/display, erase characters, alignment display.
- **Autowrap & line wrapping**: wrap mode behavior, wrap state transitions, and boundary cases.
- **Tab stops**: setting, clearing, and forward/backward tabbing.

## Drawing & Character Handling
- **ASCII/UTF-8 rendering**: drawing behavior, wide characters, combining marks, and zero-width glyph handling.
- **Character sets**: G0/G1 selection, shift in/out behavior, and charset fallback.
- **Legacy mappings**: CP437 and VT100 line drawing translation.

## SGR & Attributes
- **Standard attributes**: bold, italics, underline, blink, reverse, conceal, strikethrough.
- **Color handling**: 8/16/256/truecolor foreground/background settings plus reset semantics.
- **Private SGR handling**: unsupported/private sequences are ignored.

## Modes & Device Reports
- **Mode set/reset**: SM/RM with private mode handling.
- **Device reports**: DA/DSR responses.

## Stream Parsing
- **Escape sequence dispatch**: basic, CSI, OSC, and charset sequence parsing.
- **Parameter handling**: missing parameters, overflow clamping, and control-character interruption.
- **Byte vs. rune stream behavior**: byte-stream UTF-8 decoding and C1 control handling.
- **Debug stream output**: JSON payloads matching pyte semantics.
- **Compatibility API**: screen attach/detach and stream compatibility behavior.

## Diff & History
- **Dirty tracking**: line dirty state marking, full screen updates, and updates on wrap.
- **History buffer**: scrollback accumulation, paging, bounds behavior, and resize interactions.

## Input/Output Semantics
- **Title/icon updates**: OSC updates for icon and window titles.
- **Selection/clipboard**: OSC 52 set/query semantics.

## Known Skips
- **`test_draw_width2_irm`**: marked xfail upstream in pyte and skipped in the Go port to preserve fidelity.
