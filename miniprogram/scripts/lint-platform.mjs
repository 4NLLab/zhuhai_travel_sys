import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';

const roots = [
  'src/pages',
  'src/components',
  'src/api',
  'src/adapters',
  'src/stores',
  'src/utils',
  'dist/build/mp-weixin',
  'dist/mp-weixin'
];
const sourceFileExts = new Set(['.vue', '.ts', '.js', '.json']);
const distFileExts = new Set(['.js', '.wxml']);
const bannedPatterns = [
  { name: 'window', regex: /\bwindow\b/ },
  { name: 'document', regex: /\bdocument\b/ },
  { name: 'localStorage', regex: /\blocalStorage\b/ },
  { name: 'fetch', regex: /\bfetch\s*\(/ },
  { name: 'history', regex: /\bhistory\b/ },
  { name: 'external-script', regex: /<script[^>]+src=["']https?:\/\//i }
];

function walk(dir, allowedExts) {
  if (!existsSync(dir)) return [];
  return readdirSync(dir).flatMap((entry) => {
    const fullPath = join(dir, entry);
    const stat = statSync(fullPath);
    if (stat.isDirectory()) return walk(fullPath, allowedExts);
    const ext = fullPath.slice(fullPath.lastIndexOf('.'));
    return allowedExts.has(ext) ? [fullPath] : [];
  });
}

const violations = [];
for (const root of roots) {
  const isDistRoot = root.startsWith('dist/');
  for (const filePath of walk(root, isDistRoot ? distFileExts : sourceFileExts)) {
    const relativePath = relative(process.cwd(), filePath);
    if (isDistRoot && relativePath.endsWith('common/vendor.js')) {
      continue;
    }
    const source = readFileSync(filePath, 'utf8');
    for (const pattern of bannedPatterns) {
      if (pattern.regex.test(source)) {
        violations.push(`${relativePath}: ${pattern.name}`);
      }
    }
  }
}

if (violations.length > 0) {
  console.error('Platform API violations found:');
  for (const violation of violations) console.error(`- ${violation}`);
  process.exit(1);
}

console.log(`Platform lint passed. Scanned roots: ${roots.join(', ')}`);
