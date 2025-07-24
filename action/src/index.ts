import * as os from 'os'
import * as path from 'path'
import * as util from 'util'
import * as fs from 'fs'
import semver from 'semver'

import * as toolCache from '@actions/tool-cache'
import * as core from '@actions/core'

const toolName = 'gokcat'
const stableVersion = 'v0.7.2'
const allReleasesUrl = 'https://api.github.com/repos/philipparndt/gokcat/releases'

const getDownloadURL = (version: string): string => {
	const arch = os.arch() === "arm64" ? "arm64" : os.arch() === "x64" ? "x86_64" : os.arch()
	const platform = os.platform() === "darwin" ? "linux" : os.platform()

	return `https://github.com/philipparndt/gokcat/releases/download/${version}/gokcat_${platform}_${arch}.tar.gz`
}

const getstableVersion = async (): Promise<string> => {
	try {
		const downloadPath = await toolCache.downloadTool(allReleasesUrl)
		const responseArray = JSON.parse(fs.readFileSync(downloadPath, 'utf8').toString().trim())
		let latestVersion = semver.clean(stableVersion) || stableVersion
		responseArray.forEach((response: any) => {
			if (response && response.tag_name) {
				let currentVerison = semver.clean(response.tag_name.toString())
				if (currentVerison) {
					if (currentVerison.toString().indexOf('rc') == -1 && semver.gt(currentVerison, latestVersion)) {
						latestVersion = currentVerison
					}
				}
			}
		})
		latestVersion = "v" + latestVersion.toString()
		return latestVersion
	} catch (error) {
		core.warning(util.format("Cannot get the latest gokcat info from %s. Error %s. Using default gokcat version %s.", allReleasesUrl, error, stableVersion))
	}

	return stableVersion
}


const walkSync = (dir: string, filelist: string[], fileToFind: string) => {
	const files = fs.readdirSync(dir)
	filelist = filelist || []
	files.forEach(function (file) {
		if (fs.statSync(path.join(dir, file)).isDirectory()) {
			filelist = walkSync(path.join(dir, file), filelist, fileToFind)
		}
		else {
			core.debug(file)
			if (file == fileToFind) {
				filelist.push(path.join(dir, file))
			}
		}
	})
	return filelist
}

const downloadTool = async (version: string): Promise<string> => {
	if (!version) { version = await getstableVersion() }
	let cachedToolPath = toolCache.find(toolName, version)
	if (!cachedToolPath) {
	   let downloadPath
	   try {
		   downloadPath = await toolCache.downloadTool(getDownloadURL(version))
	   } catch (exception) {
		   throw new Error(util.format("Failed to download gokcat from location ", getDownloadURL(version)))
	   }

	   // Extract the tar.gz archive
	   let extractFolder
	   try {
		   extractFolder = await toolCache.extractTar(downloadPath)
	   } catch (exception) {
		   throw new Error(util.format("Failed to extract gokcat archive at ", downloadPath))
	   }

	   // Find the gokcat binary in the extracted folder
	   const gokcatPath = findTool(extractFolder)
	   if (!gokcatPath) {
		   throw new Error(util.format("gokcat executable not found after extraction in ", extractFolder))
	   }

	   // Cache the gokcat binary
	   cachedToolPath = await toolCache.cacheFile(gokcatPath, toolName, toolName, version)
	}

	const toolPath = findTool(cachedToolPath)
	if (!toolPath) {
		throw new Error(util.format("gokcat executable not found in path ", cachedToolPath))
	}

	fs.chmodSync(toolPath, '777')
	return toolPath
}

const findTool = (rootFolder: string): string => {
	fs.chmodSync(rootFolder, '777')
	var filelist: string[] = []
	walkSync(rootFolder, filelist, toolName)
	if (!filelist) {
		throw new Error(util.format("gokcat executable not found in path ", rootFolder))
	}
	else {
		return filelist[0]
	}
}

const run = async () => {
	let version = core.getInput('version', { 'required': true })
	if (version.toLocaleLowerCase() === 'latest') {
		version = await getstableVersion()
	} else if (!version.toLocaleLowerCase().startsWith('v')) {
		version = 'v' + version
	}

	let cachedPath = await downloadTool(version)

	try {
		const envPath = process.env['PATH'] || ''
		if (!envPath.startsWith(path.dirname(cachedPath))) {
			core.addPath(path.dirname(cachedPath))
		}
	}
	catch {
		//do nothing, set as output variable
	}

	console.log(`gokcat tool version: '${version}' has been cached at ${cachedPath}`)
	core.setOutput('gokcat-path', cachedPath)
}

run().catch(core.setFailed)
