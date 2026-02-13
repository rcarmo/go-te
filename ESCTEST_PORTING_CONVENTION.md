# ESCTest Porting Guidelines

These guidelines mirror the pyte porting conventions, adapted for the esctest2 suite in this repository. They define how to translate the Python tests to Go and how to handle emulator fixes.

## Core Rules
1. **Faithful porting only.**
   - Go tests must be a 1:1 semantic translation of the Python tests.
   - Do not “fix” a test to make it pass; fix emulator behavior instead.
2. **Do not modify vendored tests.**
   - Keep `/workspace/vendor/esctest2` untouched. All changes belong in Go files.
3. **Follow naming conventions.**
   - Use a deterministic test name that encodes the Python file and test name.
   - Recommended format: `TestEsctest<PyFileBaseCamel><PyTestNameCamel>`.
   - Example: `tests/decstbm.py::test_decstbm` → `TestEsctestDecstbmTestDecstbm`.
4. **One test case per Go test.**
   - Avoid merging multiple Python tests into one Go test.
   - Only combine cases when a Python test itself uses `@pytest.mark.parametrize`.
5. **Retain original assertions.**
   - Assert exactly the same conditions and data as the Python test.
   - Preserve ordering and expected output strings.
6. **Coordinate semantics.**
   - Keep the same 1-based/0-based expectations as in the Python tests.
   - Any existing emulator API conventions must be preserved (do not shift semantics in tests).

## Mapping Structure
- **Python file → Go test file**
  - Create a Go test file alongside existing esctest2 Go ports, named `esctest2_<py_file_base>_test.go`.
  - Example: `tests/decstbm.py` → `esctest2_decstbm_test.go`.
- **Python test → Go test**
  - Each `def test_*` becomes one Go `Test...` function.
  - Keep the test order roughly matching the Python file for readability.

## Translation Guidelines
- **Inputs**
  - Use the same escape sequences and test input data from the Python test.
  - Prefer shared helpers when they already exist in `esctest2_helpers_test.go`.
- **Assertions**
  - Preserve explicit expectations (strings, cursor positions, flags).
  - If a Python test asserts intermediate state, do so in Go as well.
- **Expected screen output**
  - If the Python test compares screen contents, use the Go display helpers.
  - Keep cell-level assertions where needed (attributes, colors, etc.).

## Emulator Fix Policy
- **Tests are not the fix.**
  - If a Go port test fails, update emulator code to match the Python test semantics.
- **Avoid reinterpreting tests.**
  - If unsure about a port, re-check the Python source and related helper utilities.

## Pyte Test Policy
- **Pyte tests are not perfect**
  - If a pyte test fails after introducing an esctest test, you need to understand which is the correct behavior

## Skips and XFail
- Only skip tests when the original Python test is marked xfail/skip.
- Document any skipped test and the upstream reason.

## Verification Workflow
1. Port a batch of tests from one Python file.
2. Run `go test ./...`.
3. Fix emulator issues that surface.
4. Repeat until the Python file’s tests all pass.

## References
- Vendored tests: `/workspace/vendor/esctest2/esctest/tests`
- Helpers: `/workspace/esctest2_helpers_test.go`
