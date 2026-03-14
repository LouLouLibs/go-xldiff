# VHS tape for recording go-xldiff demo GIF.
# Usage: vhs demo/demo.tape
#
# Prerequisites:
#   - go-xldiff binary built and in PATH
#   - Demo xlsx files created: uv run demo/gen.py

Output demo/demo.gif

Set FontSize 14
Set Width 1100
Set Height 550
Set Padding 15
Set Theme "GruvboxDarkHard"

Set TypingSpeed 80ms

Set Shell "bash"

Hide
    Type 'export PS1="> "'
    Enter
    Type "clear"
    Enter
Show

# ─── 1. SHOW THE INPUT FILES ───
Type "ls demo/*.xlsx"
Enter
Sleep 2s

# ─── 2. KEY-BASED DIFF ───
Type "go-xldiff demo/old.xlsx demo/new.xlsx --key ID"
Enter
Sleep 4s

# ─── 3. JSON OUTPUT ───
Type "go-xldiff demo/old.xlsx demo/new.xlsx --key ID --format json"
Enter
Sleep 4s
