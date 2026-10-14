class ProjectForge < Formula
  desc "Scaffolds standalone full-stack projects (Next.js, FastAPI, Flutter) from a validated project.config.yaml"
  homepage "https://github.com/tiagozaccaro/project-forge"
  version "1.0.0"

  on_macos do
    on_arm do
      url "https://github.com/tiagozaccaro/project-forge/releases/download/v#{version}/project-forge-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_ARM64_SHA256"
    end

    on_intel do
      url "https://github.com/tiagozaccaro/project-forge/releases/download/v#{version}/project-forge-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_AMD64_SHA256"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/tiagozaccaro/project-forge/releases/download/v#{version}/project-forge-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_ARM64_SHA256"
    end

    on_intel do
      url "https://github.com/tiagozaccaro/project-forge/releases/download/v#{version}/project-forge-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_AMD64_SHA256"
    end
  end

  def install
    bin.install "project-forge"
  end

  test do
    assert_match "project-forge version #{version}", shell_output("#{bin}/project-forge --version")
  end
end
