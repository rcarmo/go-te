#!/usr/bin/env python3
import ast
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]
VENDOR = ROOT / "vendor" / "pyte" / "tests"
GO = ROOT / "pkg" / "te"


def parse_tests(py_path: pathlib.Path):
    tree = ast.parse(py_path.read_text())
    tests = []
    for node in tree.body:
        if isinstance(node, ast.FunctionDef) and node.name.startswith("test_"):
            tests.append(f"{py_path.name}::{node.name}")
    return tests


def parse_go_tests():
    go_tests = set()
    for path in GO.glob("*_test.go"):
        for line in path.read_text().splitlines():
            line = line.strip()
            if line.startswith("// From pyte/tests/"):
                go_tests.add(line[len("// From "):])
    return go_tests


def main():
    go_tests = parse_go_tests()
    missing = []
    for py in sorted(VENDOR.glob("test_*.py")):
        tests = parse_tests(py)
        for t in tests:
            tag = f"pyte/tests/{t}"
            if tag not in go_tests:
                missing.append(tag)
    if missing:
        print("Missing tests:")
        for tag in missing:
            print("  ", tag)
        sys.exit(1)
    print("All pyte tests have From tags.")

if __name__ == "__main__":
    main()
