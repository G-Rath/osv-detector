#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const child_process = require('child_process');

const root = path.join(__dirname, '..');
const fixturesDir = 'fixtures/locks-e2e';

const OSV_DETECTOR_CMD = process.env.OSV_DETECTOR_CMD ?? 'osv-detector';

const files = fs
  .readdirSync(path.join(root, fixturesDir), { withFileTypes: true })
  .filter(dirent => dirent.isFile() && !dirent.name.endsWith('.out.txt'));

const runOsvDetector = (...args) => {
  const { stdout, stderr, status, error } = child_process.spawnSync(
    OSV_DETECTOR_CMD,
    args,
    { encoding: 'utf-8', cwd: root }
  );

  if (status > 1) {
    throw new Error(
      `osv-detector exited with unexpected code ${status}: ${stderr}`
    );
  }

  if (error) {
    throw error;
  }

  if (stderr.length) {
    console.warn('unexpected output to stderr', stderr);
  }

  return stdout;
};

const wildcardDatabaseStats = output => {
  return output.replaceAll(
    /(\w+) \(\d+ vulnerabilities, including withdrawn - last updated \w{3}, \d\d \w{3} \d{4} [012]\d:\d\d:\d\d GMT\)/gu,
    '$1 (%% vulnerabilities, including withdrawn - last updated %%)'
  );
};

for (const file of files) {
  const [, parseAs] = /\d+-(.*)/u.exec(file.name) ?? [];

  const p = path.join(fixturesDir, file.name);
  if (!parseAs) {
    console.error('could not determine parser for', p);
  }

  console.log('(re)generating', p, 'fixture', `(parsing as ${parseAs})`);
  const output = runOsvDetector(`${parseAs}:${p}`);

  fs.writeFileSync(
    path.join(root, fixturesDir, `${file.name}.out.txt`),
    wildcardDatabaseStats(output)
  );
}
