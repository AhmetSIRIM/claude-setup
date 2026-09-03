"""Write the opencode credentials file from the OPENCODE_GO_KEY environment variable."""
import json
import os
import sys

key = os.environ.get("OPENCODE_GO_KEY", "")
if not key:
    sys.exit("OPENCODE_GO_KEY secret is empty")
path = os.path.expanduser("~/.local/share/opencode/auth.json")
os.makedirs(os.path.dirname(path), exist_ok=True)
with open(path, "w") as fh:
    json.dump({"opencode-go": {"type": "api", "key": key}}, fh)
os.chmod(path, 0o600)
print("credentials written")
