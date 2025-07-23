
import * as core from '@actions/core';
import * as cache from '@actions/cache';
import * as io from '@actions/io';
import fetch from 'node-fetch';
import { createWriteStream, chmodSync } from 'fs';
import { pipeline } from 'stream';
import { promisify } from 'util';
import * as path from 'path';
import * as os from 'os';
import { dir as tmpDir } from 'tmp-promise';
import { execSync } from 'child_process';

const streamPipeline = promisify(pipeline);

async function run() {
  try {
    const repo = 'philipparndt/gokcat';
    const arch = os.arch() === 'arm64' ? 'arm64' : os.arch() === 'x64' ? 'x86_64' : os.arch();
    const platform = os.platform() === 'darwin' ? 'linux' : os.platform(); // gokcat only provides linux binaries
    const cacheKey = `gokcat-${platform}-${arch}`;
    const installDir = core.getInput('install-dir') || '/usr/local/bin';
    const cachePath = path.join(installDir, 'gokcat');


    // Try to restore from cache
    const restored = await cache.restoreCache([cachePath], cacheKey);
    if (restored) {
      core.info(`Restored gokcat from cache: ${restored}`);
      return;
    }

    // Get latest release info
    const apiUrl = `https://api.github.com/repos/${repo}/releases/latest`;
    const response = await fetch(apiUrl);
    const release = await response.json();
    let assetUrl = '';
    for (const asset of release.assets) {
      if (asset.name.includes(`${platform}_${arch}`) && asset.name.endsWith('.tar.gz')) {
        assetUrl = asset.browser_download_url;
        break;
      }
    }
    if (!assetUrl) {
      core.warning('Could not find latest release asset, falling back to v0.7.2');
      assetUrl = `https://github.com/${repo}/releases/download/v0.7.2/gokcat_${platform}_${arch}.tar.gz`;
    }

    // Download and extract

    // Create temporary directory
    const tmp = await tmpDir();
    const tarPath = path.join(tmp.path, 'gokcat.tar.gz');
    const res = await fetch(assetUrl);
    if (!res.ok) throw new Error(`Failed to download: ${assetUrl}`);
    await streamPipeline(res.body, createWriteStream(tarPath));

    // Extract tar.gz
    execSync(`tar -xzf ${tarPath} -C ${tmp.path}`);

    // Move binary to installDir
    await io.mkdirP(installDir);
    await io.mv(path.join(tmp.path, 'gokcat'), cachePath);
    chmodSync(cachePath, 0o755);
    core.info(`Installed gokcat to ${cachePath}`);

    // Save to cache
    try {
      await cache.saveCache([cachePath], cacheKey);
      core.info('Saved gokcat to cache.');
    } catch (e) {
      core.warning(`Could not save to cache: ${e}`);
    }

    core.setOutput('gokcat-path', cachePath);
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    } else {
      core.setFailed('Unknown error occurred');
    }
  }
}

run();
