import { existsSync, readdirSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';

const distRoot = existsSync('dist/build/mp-weixin') ? 'dist/build/mp-weixin' : 'dist/mp-weixin';
const MB = 1024 * 1024;
const LIMITS = {
  mainPackage: 2 * MB,
  subPackage: 2 * MB,
  total: 8 * MB,
  warningMainPackage: 1.8 * MB,
  singleAssetWarning: 200 * 1024
};

function walk(dir) {
  if (!existsSync(dir)) return [];
  return readdirSync(dir).flatMap((entry) => {
    const fullPath = join(dir, entry);
    const stat = statSync(fullPath);
    if (stat.isDirectory()) return walk(fullPath);
    return [{ path: fullPath, size: stat.size }];
  });
}

function bytes(value) {
  return `${(value / 1024).toFixed(1)} KB`;
}

if (!existsSync(distRoot)) {
  console.error('mp-weixin dist does not exist. Run build:mp-weixin before check:size.');
  process.exit(1);
}

const files = walk(distRoot);
const total = files.reduce((sum, file) => sum + file.size, 0);
const subPackageFiles = files.filter((file) => relative(distRoot, file.path).startsWith('subpackages/driver/'));
const subPackage = subPackageFiles.reduce((sum, file) => sum + file.size, 0);
const mainPackage = total - subPackage;
const largest = files.reduce((max, file) => (file.size > max.size ? file : max), { path: '', size: 0 });
const largeAssets = files.filter((file) => file.size > LIMITS.singleAssetWarning);

console.log(JSON.stringify({
  distRoot,
  mainPackage: bytes(mainPackage),
  subPackages: {
    'subpackages/driver': bytes(subPackage)
  },
  total: bytes(total),
  largestFile: {
    path: largest.path ? relative(distRoot, largest.path) : '',
    size: bytes(largest.size)
  },
  largeAssets: largeAssets.map((file) => ({
    path: relative(distRoot, file.path),
    size: bytes(file.size)
  }))
}, null, 2));

const failures = [];
if (mainPackage >= LIMITS.mainPackage) failures.push(`main package ${bytes(mainPackage)} >= 2 MB`);
if (subPackage >= LIMITS.subPackage) failures.push(`driver subpackage ${bytes(subPackage)} >= 2 MB`);
if (total >= LIMITS.total) failures.push(`total ${bytes(total)} >= 8 MB`);
if (mainPackage >= LIMITS.warningMainPackage) {
  console.warn(`Warning: main package ${bytes(mainPackage)} reached 1.8 MB warning line.`);
}

if (failures.length > 0) {
  console.error('Size check failed:');
  for (const failure of failures) console.error(`- ${failure}`);
  process.exit(1);
}
