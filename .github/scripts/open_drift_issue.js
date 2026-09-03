// Open the weekly drift issue, assigned to the repository owner.
// The assignment exists to trigger GitHub's e-mail notification: the issue is created
// by the actions bot, so the owner gets an "assigned" mail (own actions would not notify).
// Loaded by actions/github-script; reads the model's answer from /tmp/answer.txt.
module.exports = async ({ github, context, core }) => {
  const fs = require('fs');
  let answer = fs.readFileSync('/tmp/answer.txt', 'utf8');
  // Some models prefix planning narration; the real answer starts at the first
  // Markdown heading. Everything before it is noise, so drop it.
  const firstHeading = answer.search(/^#{1,2} /m);
  if (firstHeading > 0) answer = answer.slice(firstHeading);
  // Backstop for the prompt's length budget; keeps a runaway answer readable.
  const LIMIT = 6000;
  if (answer.length > LIMIT) {
    answer = answer.slice(0, LIMIT) + '\n\n*(truncated: the answer exceeded the length budget)*';
  }
  const model = process.env.OPENCODE_MODEL;
  const body = `Automated weekly comparison of this setup against the current Claude Code docs (\`${model}\` via opencode).\n\n${answer}\n\nSources fetched: \`code.claude.com/docs/llms.txt\`, \`/docs/en/memory\`, \`/docs/en/hooks\`.`;
  const date = new Date().toISOString().slice(0, 10);
  const latest = process.env.LATEST_VERSION || 'unknown';
  const issue = await github.rest.issues.create({
    owner: context.repo.owner, repo: context.repo.repo,
    title: `Claude Code weekly digest ${date} (through v${latest})`,
    body, assignees: [context.repo.owner],
  });
  core.info(`Opened #${issue.data.number}, assigned to ${context.repo.owner}`);
};
