"""Assemble the weekly digest prompt.

Two narrow tasks for the model:
1. Newsletter: summarize ONLY the changelog sections newer than the last covered version
   (parsed deterministically here from prior issue titles, "through vX.Y.Z"). The model
   never decides what is new and never relates items to the setup.
2. Breakage check: report a problem ONLY for something the setup actually uses that the
   current docs no longer support. Suggesting unused features is forbidden.

Prints the prompt to stdout and writes /tmp/digest-meta.env with LATEST_VERSION and
NEW_SECTION_COUNT for the workflow.
"""
import pathlib
import re

ROOT = pathlib.Path(".")
TMP = pathlib.Path("/tmp")


def read(path: pathlib.Path, cap: int = 25_000) -> str:
    text = path.read_text(encoding="utf-8", errors="replace")
    if len(text) > cap:
        text = text[:cap] + "\n[... truncated ...]"
    return text


def version_key(v: str):
    return tuple(int(x) for x in v.split("."))


# --- setup files ---
setup = [f"### CLAUDE.template.md\n{read(ROOT / 'CLAUDE.template.md')}"]
setup += [f"### rules/{p.name}\n{read(p)}" for p in sorted((ROOT / "rules").glob("*.md"))]
setup += [f"### hooks/{p.name}\n{read(p)}" for p in sorted((ROOT / "hooks").glob("*.sh"))]
setup.append(f"### settings.template.json\n{read(ROOT / 'settings.template.json')}")
setup.append(f"### README.md\n{read(ROOT / 'README.md')}")
setup_text = "\n\n".join(setup)

# --- prior issues (memory) ---
prior_path = TMP / "prior-issues.md"
prior = prior_path.read_text(encoding="utf-8", errors="replace") if prior_path.exists() else ""

# --- changelog, cropped deterministically to sections newer than the last covered version ---
covered = [m for m in re.findall(r"through v(\d+\.\d+\.\d+)", prior)]
last_covered = max(covered, key=version_key) if covered else None

changelog = (TMP / "changelog.md").read_text(encoding="utf-8", errors="replace") if (TMP / "changelog.md").exists() else ""
sections = re.split(r"(?m)^## ", changelog)
new_sections, latest = [], None
for sec in sections[1:]:
    m = re.match(r"(\d+\.\d+\.\d+)", sec)
    if not m:
        continue
    v = m.group(1)
    if latest is None:
        latest = v
    if last_covered and version_key(v) <= version_key(last_covered):
        break
    new_sections.append("## " + sec.strip())
if last_covered is None:
    new_sections = new_sections[:10]  # first ever run: last 10 releases
news_text = "\n\n".join(new_sections)[:35_000]

(TMP / "digest-meta.env").write_text(
    f"LATEST_VERSION={latest or 'unknown'}\nNEW_SECTION_COUNT={len(new_sections)}\n"
)

# --- docs for the breakage check ---
docs = "\n\n".join(
    f"### {name}\n{read(TMP / name)}"
    for name in ("llms.txt", "memory.md", "hooks.md")
    if (TMP / name).exists()
)

prior_section = (
    "\n## Previously reported\n"
    "Do not repeat anything already covered by these earlier digests, open or closed.\n"
    "When a closed issue's last comment says the setup stays as it is, treat that as an\n"
    "accepted deviation and stay silent about it permanently.\n\n" + prior + "\n"
    if prior.strip()
    else ""
)

newsletter_task = (
    f"""## Task 1: newsletter
The sections below are the Claude Code changelog entries released since the last digest.
Summarize them as a short newsletter for the user: group related items, keep only what a
daily Claude Code user would care about, drop internal fixes. Plain reporting only; do
NOT relate items to the setup and do NOT recommend adopting anything. Every bullet is one line and
starts with the version it landed in (for example `v2.1.240:`). At most 8 bullets in
total; merge or drop the rest.

{news_text}"""
    if new_sections
    else "## Task 1: newsletter\nNo new releases since the last digest. Write nothing for this task."
)

print(f"""You produce a weekly Claude Code digest for the owner of this setup. Answer in English.
Write ONLY the final answer: no reasoning narration ("Let me ..."), no notes about
what you are doing. Hard budget: the whole answer stays under 60 lines; shorter is
better.

{newsletter_task}

## Task 2: breakage check
Compare the setup below against the current docs. Report ONLY two kinds of finding:
(a) a mechanism, setting key, hook event, or file location that the setup ACTUALLY USES
and that the docs no longer support or describe differently in a way that would break
it; (b) a factual claim these files state about how Claude Code works that the current
docs contradict, even when nothing breaks. Do not mention features the
setup does not use; suggesting them is forbidden. Style, tone and preferences are out of
scope. Each finding is one `##` heading naming the file and setting, plus at most
3 lines of explanation, one of which quotes the exact sentence from the fetched
documents that proves the breakage. No quote from the documents above means the
finding does not exist; do not report it.
{prior_section}
If there is nothing to report for EITHER task, answer with exactly this single line:
NOTHING_TO_REPORT

Otherwise answer in Markdown with exactly this skeleton, dropping a section only when
it is empty:

# TL;DR
Two or three bullets: the week at a glance, and whether anything breaks this setup.

# Newsletter
The version-prefixed one-line bullets from Task 1.

# Breakage
The findings from Task 2.

## Current docs
{docs}

## Setup under review
{setup_text}
""")
