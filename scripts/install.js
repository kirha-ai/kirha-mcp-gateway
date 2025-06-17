#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');
const https = require('https');
const crypto = require('crypto');

const GITHUB_REPO = 'kirha-ai/kirha-mcp-gateway';
const GITHUB_API_URL = `https://api.github.com/repos/${GITHUB_REPO}/releases/latest`;

function getBinaryInfo() {
  const platform = os.platform();
  const arch = os.arch();
  
  let binaryName = 'kirha-mcp-gateway';
  let ext = '';
  
  // Map Node.js arch to Go arch
  let goArch = arch;
  if (arch === 'x64') goArch = 'amd64';
  
  if (platform === 'win32') {
    binaryName += `-windows-${goArch}`;
    ext = '.exe';
  } else if (platform === 'darwin') {
    binaryName += `-darwin-${goArch}`;
  } else if (platform === 'linux') {
    binaryName += `-linux-${goArch}`;
  } else {
    console.error(`Unsupported platform: ${platform}`);
    process.exit(1);
  }
  
  return {
    name: binaryName + ext,
    checksumName: binaryName + ext + '.sha256'
  };
}

function downloadFile(url, destination) {
  return new Promise((resolve, reject) => {
    console.log(`Downloading ${url}`);
    const file = fs.createWriteStream(destination);
    
    https.get(url, {
      headers: {
        'User-Agent': 'kirha-mcp-gateway-npm-installer'
      }
    }, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        return downloadFile(response.headers.location, destination)
          .then(resolve)
          .catch(reject);
      }
      
      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode} ${response.statusMessage}`));
        return;
      }
      
      response.pipe(file);
      
      file.on('finish', () => {
        file.close();
        resolve();
      });
      
      file.on('error', (error) => {
        fs.unlink(destination, () => {}); // Delete partial file
        reject(error);
      });
    }).on('error', reject);
  });
}

function getLatestRelease() {
  return new Promise((resolve, reject) => {
    https.get(GITHUB_API_URL, {
      headers: {
        'User-Agent': 'kirha-mcp-gateway-npm-installer'
      }
    }, (response) => {
      let data = '';
      
      response.on('data', (chunk) => {
        data += chunk;
      });
      
      response.on('end', () => {
        if (response.statusCode !== 200) {
          reject(new Error(`GitHub API request failed: ${response.statusCode}`));
          return;
        }
        
        try {
          const release = JSON.parse(data);
          resolve(release);
        } catch (error) {
          reject(new Error(`Failed to parse GitHub API response: ${error.message}`));
        }
      });
    }).on('error', reject);
  });
}

function verifyChecksum(filePath, checksumPath) {
  if (!fs.existsSync(checksumPath)) {
    console.warn('Checksum file not found, skipping verification');
    return true;
  }
  
  const expectedChecksum = fs.readFileSync(checksumPath, 'utf8').trim().split(' ')[0];
  const fileBuffer = fs.readFileSync(filePath);
  const actualChecksum = crypto.createHash('sha256').update(fileBuffer).digest('hex');
  
  if (actualChecksum !== expectedChecksum) {
    console.error(`Checksum verification failed!`);
    console.error(`Expected: ${expectedChecksum}`);
    console.error(`Actual: ${actualChecksum}`);
    return false;
  }
  
  console.log('✅ Checksum verification passed');
  return true;
}

async function main() {
  const binaryInfo = getBinaryInfo();
  const binDir = path.join(__dirname, '..', 'bin');
  const binaryPath = path.join(binDir, binaryInfo.name);
  const checksumPath = path.join(binDir, binaryInfo.checksumName);
  const finalBinaryPath = path.join(binDir, 'kirha-mcp-gateway' + (os.platform() === 'win32' ? '.exe' : ''));
  
  console.log(`Installing kirha-mcp-gateway for ${os.platform()} ${os.arch()}`);
  console.log(`Binary: ${binaryInfo.name}`);
  
  // Create bin directory if it doesn't exist
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }
  
  // Check if binary already exists
  if (fs.existsSync(finalBinaryPath)) {
    console.log('Binary already exists, skipping download');
    
    // Make binary executable on Unix systems
    if (os.platform() !== 'win32') {
      try {
        fs.chmodSync(finalBinaryPath, '755');
      } catch (error) {
        console.warn(`Failed to make binary executable: ${error.message}`);
      }
    }
    
    console.log(`✅ kirha-mcp-gateway is ready for ${os.platform()} ${os.arch()}`);
    return;
  }
  
  try {
    // Get latest release info
    console.log('Fetching latest release information...');
    const release = await getLatestRelease();
    
    // Find the binary asset
    const binaryAsset = release.assets.find(asset => asset.name === binaryInfo.name);
    const checksumAsset = release.assets.find(asset => asset.name === binaryInfo.checksumName);
    
    if (!binaryAsset) {
      console.error(`Binary not found for your platform: ${binaryInfo.name}`);
      console.error(`Platform: ${os.platform()}, Architecture: ${os.arch()}`);
      console.error('Available binaries:');
      release.assets.forEach(asset => {
        if (asset.name.startsWith('kirha-mcp-gateway-')) {
          console.error(`  - ${asset.name}`);
        }
      });
      process.exit(1);
    }
    
    // Download binary
    await downloadFile(binaryAsset.browser_download_url, binaryPath);
    console.log(`✅ Downloaded ${binaryInfo.name}`);
    
    // Download and verify checksum if available
    if (checksumAsset) {
      await downloadFile(checksumAsset.browser_download_url, checksumPath);
      
      if (!verifyChecksum(binaryPath, checksumPath)) {
        fs.unlinkSync(binaryPath);
        fs.unlinkSync(checksumPath);
        console.error('Installation failed due to checksum mismatch');
        process.exit(1);
      }
      
      // Clean up checksum file
      fs.unlinkSync(checksumPath);
    }
    
    // Make binary executable on Unix systems
    if (os.platform() !== 'win32') {
      try {
        fs.chmodSync(binaryPath, '755');
        console.log(`Made ${binaryInfo.name} executable`);
      } catch (error) {
        console.warn(`Failed to make binary executable: ${error.message}`);
      }
    }
    
    // Rename binary to remove platform suffix
    try {
      fs.renameSync(binaryPath, finalBinaryPath);
      console.log(`Renamed ${binaryInfo.name} to ${path.basename(finalBinaryPath)}`);
    } catch (error) {
      console.error(`Failed to rename binary: ${error.message}`);
      fs.unlinkSync(binaryPath);
      process.exit(1);
    }
    
    console.log(`✅ kirha-mcp-gateway ${release.tag_name} is ready for ${os.platform()} ${os.arch()}`);
    
  } catch (error) {
    console.error(`Installation failed: ${error.message}`);
    
    // Clean up partial downloads
    if (fs.existsSync(binaryPath)) {
      fs.unlinkSync(binaryPath);
    }
    if (fs.existsSync(checksumPath)) {
      fs.unlinkSync(checksumPath);
    }
    
    console.error('\nTroubleshooting:');
    console.error('1. Check your internet connection');
    console.error('2. Verify the GitHub repository exists and has releases');
    console.error('3. Try installing again later');
    console.error('4. Download manually from: https://go.kirha.ai/kirha-mcp-gateway/releases');
    console.error(`5. Place the binary manually in: ${finalBinaryPath}`);
    
    process.exit(1);
  }
}

if (require.main === module) {
  main().catch(console.error);
}