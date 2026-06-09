const { execSync, spawn } = require('child_process');
const path = require('path');

const root = __dirname;
const bin = path.resolve(root, 'tmp', 'main.exe');

try {
  console.log('[build] Building...');
  execSync('go build -o ./tmp/main.exe .', { stdio: 'inherit', cwd: root });
  console.log('[build] OK, starting...');
} catch (e) {
  console.error('[build] FAILED');
  process.exit(1);
}

const child = spawn(bin, [], { stdio: 'inherit', cwd: root });
child.on('exit', (code) => process.exit(code ?? 0));
process.on('SIGTERM', () => child.kill());
process.on('SIGINT', () => child.kill());
