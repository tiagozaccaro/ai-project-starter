// Postinstall script: downloads the precompiled ai-project-starter
// binary matching the current platform from GitHub Releases and
// places it in bin/ so the package.json "bin" entry resolves.
const { createWriteStream, mkdirSync, chmodSync, existsSync } = require("node:fs");
const { platform, arch } = require("node:os");
const { join } = require("node:path");
const { pipeline } = require("node:stream/promises");
const { createGunzip } = require("node:zlib");

const VERSION = require("./package.json").version;
const OWNER = "tiagozaccaro";
const REPO = "ai-project-starter";

const GOOS = platform() === "darwin" ? "darwin" : "linux";
const NODE_TO_GOARCH = { x64: "amd64", arm64: "arm64" };
const GOARCH = NODE_TO_GOARCH[arch()];

if (!GOARCH) {
  console.error(`@ai-project-starter/cli: unsupported architecture ${arch()}`);
  process.exit(1);
}

const asset = `ai-project-starter-${GOOS}-${GOARCH}.tar.gz`;
const url = `https://github.com/${OWNER}/${REPO}/releases/download/v${VERSION}/${asset}`;

async function main() {
  const binDir = join(__dirname, "bin");
  if (!existsSync(binDir)) {
    mkdirSync(binDir, { recursive: true });
  }
  const target = join(binDir, "ai-project-starter");

  const res = await fetch(url, { redirect: "follow" });
  if (!res.ok || !res.body) {
    throw new Error(`failed to download ${url}: ${res.status} ${res.statusText}`);
  }

  // The tar.gz contains a single binary named "ai-project-starter";
  // extract it by streaming the gzip payload and skipping the tar
  // header (512-byte blocks) up to the file data.
  const gunzip = createGunzip();
  await pipeline(res.body, gunzip);

  // tar format: 512-byte header + padded content. Find the first
  // non-header block that holds the binary.
  const { Readable } = require("node:stream");
  const chunks = [];
  for await (const chunk of gunzip) {
    chunks.push(chunk);
  }
  const buf = Buffer.concat(chunks);
  const nameEnd = buf.indexOf(0x00, 0);
  const headerName = buf.subarray(0, nameEnd).toString("utf8").replace(/^\.\//, "");
  const contentStart = 512;
  const sizeOctal = buf.subarray(124, 136).toString("utf8").replace(/\0.*$/, "").trim();
  const size = parseInt(sizeOctal, 8);
  if (!headerName.includes("ai-project-starter")) {
    throw new Error(`unexpected tar layout in ${asset}`);
  }
  const data = buf.subarray(contentStart, contentStart + size);

  const out = createWriteStream(target, { mode: 0o755 });
  out.write(data);
  out.end();
  await new Promise((resolve, reject) => {
    out.on("finish", resolve);
    out.on("error", reject);
  });
  chmodSync(target, 0o755);
  console.log(`@ai-project-starter/cli: installed ${GOOS}-${GOARCH} binary`);
}

main().catch((err) => {
  console.error(`@ai-project-starter/cli: ${err.message}`);
  process.exit(1);
});
