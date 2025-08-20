#!/usr/bin/env node

const fs = require('fs/promises');
const path = require('path');
const child_process = require('child_process');

const root = path.join(__dirname, '..');
const testdataDir = 'testdata/locks-e2e';

const OSV_DETECTOR_CMD = process.env.OSV_DETECTOR_CMD ?? 'osv-detector';

const runOsvDetector = async (...args) => {
  return new Promise((resolve, reject) => {
    const child = child_process.spawn(OSV_DETECTOR_CMD, args, {
      encoding: 'utf-8',
      cwd: root
    });

    let stdout = '';
    let stderr = '';

    child.stdout.on('data', data => {
      stdout += data;
    });

    child.stderr.on('data', data => {
      stderr += data;
    });

    child.on('error', reject);

    child.on('close', status => {
      if (status > 1) {
        reject(
          new Error(
            `osv-detector exited with unexpected code ${status}: ${stderr}`
          )
        );
      } else if (stderr.length) {
        console.warn('unexpected output to stderr', stderr);
      }

      resolve(stdout);
    });
  });
};

const wildcardDatabaseStats = output => {
  return output.replaceAll(
    /(\w+) \(\d+ vulnerabilities, including withdrawn - last updated \w{3}, \d\d \w{3} \d{4} [012]\d:\d\d:\d\d GMT\)/gu,
    '$1 (%% vulnerabilities, including withdrawn - last updated %%)'
  );
};

const regenerateFixture = async fileName => {
  const [, parseAs] = /\d+-(.*)/u.exec(fileName) ?? [];

  const p = path.join(testdataDir, fileName);
  if (!parseAs) {
    console.error('could not determine parser for', p);
  }

  const output = await runOsvDetector(`${parseAs}:${p}`);

  console.log('(re)generated', p, 'fixture', `(parsed as ${parseAs})`);

  await fs.writeFile(
    path.join(root, testdataDir, `${fileName}.out.txt`),
    wildcardDatabaseStats(output)
  );
};

(async () => {
  const files = (
    await fs.readdir(path.join(root, testdataDir), { withFileTypes: true })
  ).filter(dirent => dirent.isFile() && !dirent.name.endsWith('.out.txt'));

  await Promise.all(files.map(file => regenerateFixture(file.name)));
})();
