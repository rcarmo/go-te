#!/usr/bin/env python3
import ast
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]
VENDOR = ROOT / "vendor" / "esctest2" / "esctest" / "tests"
GO = ROOT / "pkg" / "te"

ignore = {"__init__.py"}

def parse_tests(py_path: pathlib.Path):
    tree = ast.parse(py_path.read_text())
    tests = []
    for node in tree.body:
        if isinstance(node, ast.ClassDef):
            for item in node.body:
                if isinstance(item, ast.FunctionDef) and item.name.startswith("test_"):
                    tests.append(f"{py_path.name}::{item.name}")
        if isinstance(node, ast.FunctionDef) and node.name.startswith("test_"):
            tests.append(f"{py_path.name}::{node.name}")
    return tests


def parse_go_tests():
    go_tests = {}
    for path in GO.glob("esctest2_*_test.go"):
        text = path.read_text()
        for line in text.splitlines():
            line = line.strip()
            if line.startswith("// From esctest2/esctest/tests/"):
                tag = line[len("// From "):]
                key = tag.split("::")[0].split("/")[-1]
                go_tests.setdefault(key, set()).add(tag)
    return go_tests


def main():
    go_tests = parse_go_tests()
    missing = []
    for py in sorted(VENDOR.glob("*.py")):
        if py.name in ignore:
            continue
        tests = parse_tests(py)
        missing_for_file = []
        for t in tests:
            tag = f"esctest2/esctest/tests/{t}"
            if tag not in go_tests.get(py.name, set()):
                missing_for_file.append(tag)
        if missing_for_file:
            missing.append((py.name, missing_for_file))
    if missing:
        print("Missing tests:")
        for fname, items in missing:
            print(fname)
            for item in items:
                print("  ", item)
        sys.exit(1)
    print("All esctest2 tests have From tags.")

if __name__ == "__main__":
    main()
