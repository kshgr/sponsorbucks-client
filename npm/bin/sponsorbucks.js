#!/usr/bin/env node
const os = require("os");
const fs = require("fs");
const path = require("path");
const https = require("https");
const { spawnSync } = require("child_process");

const repo = process.env.SPONSORBUCKS_REPO || "kshgr/sponsorbucks-client";
const version = process.env.SPONSORBUCKS_VERSION || "latest";
const base = process.env.SPONSORBUCKS_RELEASE_BASE;
const apiBase = process.env.SPONSORBUCKS_API_BASE || "https://enbmzimfbtnfrkpevzzf.supabase.co/functions/v1";
const platform = os.platform();
const arch = os.arch();

function assetName() {
  if (platform === "win32" && arch === "x64") return "sponsorbucks-windows-amd64.exe";
  if (platform === "linux" && arch === "x64") return "sponsorbucks-linux-amd64";
  if (platform === "darwin" && arch === "x64") return "sponsorbucks-darwin-amd64";
  if (platform === "darwin" && arch === "arm64") return "sponsorbucks-darwin-arm64";
  throw new Error(`Unsupported platform: ${platform}/${arch}`);
}

function urlFor(asset) {
  if (base) return `${base.replace(/\/$/, "")}/${asset}`;
  if (version === "latest") return `https://github.com/${repo}/releases/latest/download/${asset}`;
  return `https://github.com/${repo}/releases/download/${version}/${asset}`;
}

function binDir() {
  return path.join(os.homedir(), ".sponsorbucks", "bin");
}

function binaryPath() {
  return path.join(binDir(), platform === "win32" ? "sponsorbucks.exe" : "sponsorbucks");
}

function download(url, target) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(target);
    const req = https.get(url, { headers: { "User-Agent": "SponsorBucks-NPX" } }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        file.close();
        fs.unlinkSync(target);
        download(res.headers.location, target).then(resolve, reject);
        return;
      }
      if (res.statusCode !== 200) {
        file.close();
        fs.unlinkSync(target);
        reject(new Error(`Download failed: ${res.statusCode} ${res.statusMessage}`));
        return;
      }
      res.pipe(file);
      file.on("finish", () => file.close(resolve));
    });
    req.on("error", (err) => {
      file.close();
      try { fs.unlinkSync(target); } catch {}
      reject(err);
    });
  });
}

function configureDefaultApi(bin) {
  const result = spawnSync(bin, ["config", "set-api", apiBase], { stdio: "ignore" });
  return result.status === 0;
}

async function ensureBinary() {
  const target = binaryPath();
  if (fs.existsSync(target)) return target;
  fs.mkdirSync(binDir(), { recursive: true });
  const asset = assetName();
  const url = urlFor(asset);
  const tmp = `${target}.download`;
  console.error(`Downloading SponsorBucks preview from ${url}`);
  await download(url, tmp);
  fs.renameSync(tmp, target);
  if (platform !== "win32") fs.chmodSync(target, 0o755);
  return target;
}

(async () => {
  try {
    const bin = await ensureBinary();
    configureDefaultApi(bin);
    const result = spawnSync(bin, process.argv.slice(2), { stdio: "inherit" });
    process.exit(result.status ?? 0);
  } catch (err) {
    console.error("SponsorBucks npx launcher failed:");
    console.error(err.message);
    console.error("Set SPONSORBUCKS_RELEASE_BASE to a directory containing the release binary, or publish a GitHub release first.");
    process.exit(1);
  }
})();
