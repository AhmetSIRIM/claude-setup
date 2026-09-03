#!/usr/bin/env bash
# Session-start health check for this setup. Prints one warning line per problem
# and stays completely silent when healthy, so a healthy session adds no context.
# Catches configuration that breaks silently: an enabled code-intelligence plugin
# whose language server binary is missing, a hook referenced by settings but absent
# on disk, and a broken ~/.claude symlink.
# CLAUDE_DIR overrides the config directory (used by tests).
# The interpreter is resolved by test-running each candidate, because on Windows a
# python3 shim can exist on PATH and still refuse to run (the Store alias trap).
PYTHON=""
for candidate in python3 python "py -3"; do
  if $candidate -c 'pass' >/dev/null 2>&1; then PYTHON="$candidate"; break; fi
done
if [ -z "$PYTHON" ]; then
  echo "setup-doctor: no working python interpreter found; checks skipped"
  exit 0
fi
$PYTHON - <<'PY'
import json, os, shutil, subprocess

claude_dir = os.environ.get('CLAUDE_DIR', os.path.expanduser('~/.claude'))
warns = []

LSP_BINARIES = {
    'kotlin-lsp': 'kotlin-language-server', 'swift-lsp': 'sourcekit-lsp',
    'gopls-lsp': 'gopls', 'typescript-lsp': 'typescript-language-server',
    'jdtls-lsp': 'jdtls', 'pyright-lsp': 'pyright-langserver',
    'rust-analyzer-lsp': 'rust-analyzer', 'clangd-lsp': 'clangd',
    'csharp-lsp': 'csharp-ls', 'lua-lsp': 'lua-language-server', 'php-lsp': 'intelephense',
}
# Hooks may run without the login shell's PATH; search the usual install dirs too.
EXTRA_DIRS = [os.path.expanduser('~/go/bin'), '/opt/homebrew/bin', '/usr/local/bin']

def binary_exists(name):
    if shutil.which(name):
        return True
    if any(os.access(os.path.join(d, name), os.X_OK) for d in EXTRA_DIRS):
        return True
    if name == 'sourcekit-lsp':
        try:
            return subprocess.run(['xcrun', '--find', 'sourcekit-lsp'],
                                  capture_output=True).returncode == 0
        except OSError:
            return False
    return False

try:
    settings = json.load(open(os.path.join(claude_dir, 'settings.json')))
except Exception as e:
    print(f'setup-doctor: settings.json unreadable: {e}')
    raise SystemExit(0)

for plugin_id, enabled in (settings.get('enabledPlugins') or {}).items():
    name = plugin_id.split('@')[0]
    if enabled and name in LSP_BINARIES and not binary_exists(LSP_BINARIES[name]):
        warns.append(f"plugin '{name}' is enabled but its language server "
                     f"'{LSP_BINARIES[name]}' is not installed")

# A plugin enabled in settings but never installed silently does nothing.
try:
    installed = json.load(open(os.path.join(claude_dir, 'plugins', 'installed_plugins.json')))
    installed_ids = set(installed.get('plugins') or {})
    for plugin_id, enabled in (settings.get('enabledPlugins') or {}).items():
        if enabled and plugin_id not in installed_ids:
            warns.append(f"plugin '{plugin_id}' is enabled in settings but not installed; "
                         f"run: claude plugin install {plugin_id}")
except FileNotFoundError:
    pass
except Exception as e:
    warns.append(f'could not read installed_plugins.json: {e}')

for entry in ('rules', 'skills', 'hooks', 'CLAUDE.md'):
    path = os.path.join(claude_dir, entry)
    if os.path.islink(path) and not os.path.exists(path):
        warns.append(f'~/.claude/{entry} is a broken symlink')

claude_md = os.path.join(claude_dir, 'CLAUDE.md')
if not os.path.exists(claude_md):
    warns.append('~/.claude/CLAUDE.md is missing; copy CLAUDE.template.md from the setup repo and fill it in')
else:
    # Content-based check: drop HTML comments and the template's scaffold lines;
    # whatever remains is the owner's own writing. A kept template comment above
    # real content therefore does not trigger the warning.
    try:
        import re
        body = re.sub(r'<!--.*?-->', '', open(claude_md).read(), flags=re.S)
        scaffold = re.compile(r'^\s*$'
                              r'|^[A-Z][A-Za-z &]+:\s*$'
                              r'|^- Detailed working rules ')
        own_lines = [l for l in body.splitlines() if not scaffold.match(l)]
        if not own_lines:
            warns.append('~/.claude/CLAUDE.md is still the unfilled template; make it yours')
    except OSError:
        pass

for event_entries in (settings.get('hooks') or {}).values():
    for entry in event_entries:
        for hook in entry.get('hooks', []):
            command = hook.get('command', '')
            first = os.path.expanduser(command.split(' ')[0]) if command else ''
            if first.startswith('/') or command.startswith('~'):
                if not os.path.exists(first):
                    warns.append(f'settings references a missing hook: {command}')

status_line = settings.get('statusLine') or {}
if status_line.get('type') == 'command':
    # The command may be a bare program with path arguments ('bash "$HOME/x.sh"'),
    # so check every token that expands to an absolute path, not the first word.
    for token in status_line.get('command', '').split():
        expanded = os.path.expanduser(os.path.expandvars(token.strip('"\'')))
        if expanded.startswith('/') and not os.path.exists(expanded):
            warns.append(f'statusLine references a missing file: {token}')

for w in warns:
    print('setup-doctor:', w)
PY
