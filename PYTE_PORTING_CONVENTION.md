# Pyte Test Porting Convention (Go)

Pyte tests MUST be ported to Go with 100% fidelity. Any discrepencies between that logic and terminal emulator behavior must be fixed in the emulator.

## File naming
- Use `pyte_<pyfile>_test.go` where `<pyfile>` is the base name of the Python test file.
  - Examples:
    - `test_screen.py` → `pyte_screen_test.go`
    - `test_stream.py` → `pyte_stream_test.go`
    - `test_history.py` → `pyte_history_test.go`
    - `test_diff.py` → `pyte_diff_test.go`

## Test function naming
- Use `TestPyte<PyFileBaseCamel><PyTestNameInCamelCase>` to avoid name collisions across modules.
  - Convert the Python file base (e.g., `test_screen`) to CamelCase and append the CamelCase form of the test name after `test_`.
  - Examples:
    - `test_screen.py::test_draw_width2_line_end` → `TestPyteTestScreenDrawWidth2LineEnd`
    - `test_diff.py::test_mark_whole_screen` → `TestPyteTestDiffMarkWholeScreen`

## Required comment annotation
- Add a comment immediately above each Go test:
  - `// From pyte/tests/<pyfile>::<py_test_name>`
  - Example:
    - `// From pyte/tests/test_screen.py::test_draw_width2_line_end`

## Fidelity rule
- The Go test must be a **100% faithful port** of the Python test.
- Once a Go test matches the Python semantics, it must not be modified (only the emulator may change).
