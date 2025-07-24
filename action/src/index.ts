import * as core from "@actions/core"
import * as tc from "@actions/tool-cache"
import * as io from "@actions/io"
import {chmodSync} from "fs"
import * as path from "path"
import * as os from "os"
import {Octokit} from "@octokit/rest"

async function run() {
	try {
		const arch = os.arch() === "arm64" ? "arm64" : os.arch() === "x64" ? "x86_64" : os.arch();
		const platform = os.platform() === "darwin" ? "linux" : os.platform();
		const version = "latest"; // You can make this configurable if needed
		const toolName = "gokcat";
		const installDir = core.getInput("install-dir") || "/usr/local/bin";

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

		// Get latest release info
		let assetUrl = "";
		const octokit = new Octokit();
		const release = await octokit.repos.getLatestRelease({owner: "philipparndt", repo: "gokcat"});
		if (!release.data.assets || release.data.assets.length === 0) {
			core.warning("Could not find latest release asset, falling back to v0.7.2");
			assetUrl = `https://github.com/philipparndt/gokcat/releases/download/v0.7.2/gokcat_${platform}_${arch}.tar.gz`;
		} else {
			for (const asset of release.data.assets) {
				if (asset.name.includes(`${platform}_${arch}`) && asset.name.endsWith(".tar.gz")) {
					assetUrl = asset.browser_download_url;
					break;
				}
			}
		}

		if (!assetUrl) {
			throw new Error("Could not determine asset URL for gokcat");
		}

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
