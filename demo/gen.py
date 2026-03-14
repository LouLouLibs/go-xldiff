# /// script
# requires-python = ">=3.11"
# dependencies = ["openpyxl"]
# ///
"""Generate sample Excel files for the VHS demo."""

import openpyxl


from pathlib import Path

DIR = Path(__file__).resolve().parent


def write_file(name, rows):
    wb = openpyxl.Workbook()
    ws = wb.active
    ws.title = "Sheet1"
    for row in rows:
        ws.append(row)
    wb.save(DIR / name)


write_file("old.xlsx", [
    ["ID", "Name", "Department", "Salary"],
    [1, "Alice", "Engineering", 95000],
    [2, "Bob", "Marketing", 72000],
    [3, "Charlie", "Engineering", 88000],
    [4, "Dana", "Sales", 65000],
])

write_file("new.xlsx", [
    ["ID", "Name", "Department", "Salary"],
    [1, "Alice", "Engineering", 98000],
    [2, "Bob", "Design", 75000],
    [4, "Dana", "Sales", 65000],
    [5, "Eve", "Marketing", 70000],
])

print(f"Created {DIR / 'old.xlsx'} and {DIR / 'new.xlsx'}")
