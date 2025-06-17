#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

function getBinaryName() {
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
  
  return binaryName + ext;
}

function main() {
  const binaryName = getBinaryName();
  const binaryPath = path.join(__dirname, binaryName);
  
  // Pass all arguments to the Go binary
  const args = process.argv.slice(2);
  
  const child = spawn(binaryPath, args, {
    stdio: 'inherit',
    env: process.env
  });
  
  child.on('error', (error) => {
    if (error.code === 'ENOENT') {
      console.error(`Binary not found: ${binaryPath}`);
      console.error('Make sure the binary is built for your platform.');
      console.error(`Platform: ${os.platform()}, Architecture: ${os.arch()}`);
    } else {
      console.error('Failed to start kirha-mcp-gateway:', error.message);
    }
    process.exit(1);
  });
  
  child.on('close', (code) => {
    process.exit(code);
  });
  
  // Handle signals
  process.on('SIGINT', () => {
    child.kill('SIGINT');
  });
  
  process.on('SIGTERM', () => {
    child.kill('SIGTERM');
  });
}

if (require.main === module) {
  main();
}