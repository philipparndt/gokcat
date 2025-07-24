"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const core = __importStar(require("@actions/core"));
const cache = __importStar(require("@actions/cache"));
const io = __importStar(require("@actions/io"));
const node_fetch_1 = __importDefault(require("node-fetch"));
const fs_1 = require("fs");
const stream_1 = require("stream");
const util_1 = require("util");
const path = __importStar(require("path"));
const os = __importStar(require("os"));
const tmp_promise_1 = require("tmp-promise");
const child_process_1 = require("child_process");
const rest_1 = require("@octokit/rest");
const streamPipeline = (0, util_1.promisify)(stream_1.pipeline);
async function run() {
    try {
        const arch = os.arch() === "arm64" ? "arm64" : os.arch() === "x64" ? "x86_64" : os.arch();
        const platform = os.platform() === "darwin" ? "linux" : os.platform();
        const cacheKey = `gokcat-${platform}-${arch}`;
        const installDir = core.getInput("install-dir") || "/usr/local/bin";
        const cachePath = path.join(installDir, "gokcat");
        // Try to restore from cache
        const restored = await cache.restoreCache([cachePath], cacheKey);
        if (restored) {
            core.info(`Restored gokcat from cache: ${restored}`);
            return;
        }
        // Get latest release info
        let assetUrl = "";
        const octokit = new rest_1.Octokit();
        const release = await octokit.repos.getLatestRelease({ owner: "philipparndt", repo: "gokcat" });
        if (!release.data.assets || release.data.assets.length === 0) {
            core.warning("Could not find latest release asset, falling back to v0.7.2");
            assetUrl = `https://github.com/philipparndt/gokcat/releases/download/v0.7.2/gokcat_${platform}_${arch}.tar.gz`;
        }
        else {
            for (const asset of release.data.assets) {
                if (asset.name.includes(`${platform}_${arch}`) && asset.name.endsWith(".tar.gz")) {
                    assetUrl = asset.browser_download_url;
                    break;
                }
            }
        }
        // Create temporary directory
        const tmp = await (0, tmp_promise_1.dir)();
        const tarPath = path.join(tmp.path, "gokcat.tar.gz");
        const res = await (0, node_fetch_1.default)(assetUrl);
        if (!res.ok)
            throw new Error(`Failed to download: ${assetUrl}`);
        await streamPipeline(res.body, (0, fs_1.createWriteStream)(tarPath));
        // Extract tar.gz
        (0, child_process_1.execSync)(`tar -xzf ${tarPath} -C ${tmp.path}`);
        // Move binary to installDir
        await io.mkdirP(installDir);
        await io.mv(path.join(tmp.path, "gokcat"), cachePath);
        (0, fs_1.chmodSync)(cachePath, 0o755);
        core.info(`Installed gokcat to ${cachePath}`);
        // Save to cache
        try {
            await cache.saveCache([cachePath], cacheKey);
            core.info("Saved gokcat to cache.");
        }
        catch (e) {
            core.warning(`Could not save to cache: ${e}`);
        }
        core.setOutput("gokcat-path", cachePath);
    }
    catch (error) {
        if (error instanceof Error) {
            core.setFailed(error.message);
        }
        else {
            core.setFailed("Unknown error occurred");
        }
    }
}
run().then();
