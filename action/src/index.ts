import * as core from "@actions/core"
import * as tc from "@actions/tool-cache"
import * as io from "@actions/io"
import {chmodSync} from "fs"
import * as path from "path"
import * as os from "os"
import {Octokit} from "@octokit/rest"

const defaultVersion = "v0.7.2";

async function getLatestVersion(platform: string, arch: string) {
	try {
		const octokit = new Octokit();
		const release = await octokit.repos.getLatestRelease({owner: "philipparndt", repo: "gokcat"});
		for (const asset of release.data.assets) {
			if (asset.name.includes(`${platform}_${arch}`) && asset.name.endsWith(".tar.gz")) {
				return release.data.name || defaultVersion; // Use the release name as version, fallback to default if not available
			}
		}
	} catch (error) {
		core.warning(`Error fetching latest version, defaulting to ${defaultVersion} ${error}`)
	}

	return defaultVersion
}

async function run() {
	try {
		const arch = os.arch() === "arm64" ? "arm64" : os.arch() === "x64" ? "x86_64" : os.arch();
		const platform = os.platform() === "darwin" ? "linux" : os.platform();
		const toolName = "gokcat";
		const installDir = core.getInput("install-dir") || "/usr/local/bin";
		const version = await getLatestVersion(platform, arch);

		// Try to find in tool cache
		let toolPath = tc.find(toolName, version, arch);
		if (toolPath) {
			const destPath = path.join(installDir, toolName);
			await io.mkdirP(installDir);
			await io.cp(path.join(toolPath, toolName), destPath);
			chmodSync(destPath, 0o755);
			core.info(`Found gokcat in tool cache: ${destPath}`);
			core.setOutput("gokcat-path", destPath);
			return;
		}

		const assetUrl = `https://github.com/philipparndt/gokcat/releases/download/${version}/gokcat_${platform}_${arch}.tar.gz`;

		// Download and extract
		const downloadPath = await tc.downloadTool(assetUrl);
		const extractPath = await tc.extractTar(downloadPath);
		const gokcatPath = path.join(extractPath, toolName);

		// Cache the extracted binary
		toolPath = await tc.cacheFile(gokcatPath, toolName, toolName, version, arch);
		core.info(`Cached gokcat at ${toolPath}`);

		// Move binary to installDir
		const destPath = path.join(installDir, toolName);
		await io.mkdirP(installDir);
		await io.cp(gokcatPath, destPath);
		chmodSync(destPath, 0o755);
		core.info(`Installed gokcat to ${destPath}`);
		core.setOutput("gokcat-path", destPath);
	} catch (error) {
		if (error instanceof Error) {
			core.setFailed(error.message);
		} else {
			core.setFailed("Unknown error occurred");
		}
	}
}

run().then();
