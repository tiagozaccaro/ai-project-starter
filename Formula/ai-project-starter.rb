class AiProjectStarter < Formula
  desc "Scaffolds standalone full-stack projects (Next.js, FastAPI, Flutter) from a validated project.config.yaml"
  homepage "https://github.com/tiagozaccaro/ai-project-starter"
  version "1.0.0"

  on_macos do
    on_arm do
      url "https://github.com/tiagozaccaro/ai-project-starter/releases/download/v#{version}/ai-project-starter-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_ARM64_SHA256"
    end

    on_intel do
      url "https://github.com/tiagozaccaro/ai-project-starter/releases/download/v#{version}/ai-project-starter-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_AMD64_SHA256"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/tiagozaccaro/ai-project-starter/releases/download/v#{version}/ai-project-starter-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_ARM64_SHA256"
    end

    on_intel do
      url "https://github.com/tiagozaccaro/ai-project-starter/releases/download/v#{version}/ai-project-starter-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_AMD64_SHA256"
    end
  end

  def install
    bin.install "ai-project-starter"
  end

  test do
    assert_match "ai-project-starter version #{version}", shell_output("#{bin}/ai-project-starter --version")
  end
end
