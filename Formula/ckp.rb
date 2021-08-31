# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ckp < Formula
  desc ""
  homepage "https://github.com/elhmn/ckp"
  version "0.3.0"
  bottle :unneeded

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.3.0/ckp_0.3.0_darwin_amd64.tar.gz"
      sha256 "943f40f108a122c943afbea5c706e14a36fe801c8631b9255287ffda98a739dd"
    end
    if Hardware::CPU.arm?
      url "https://github.com/elhmn/ckp/releases/download/v0.3.0/ckp_0.3.0_darwin_arm64.tar.gz"
      sha256 "eefce98b67d692808d39f2b2b808bf202217235f2f0904e66bc16b19cb2be409"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.3.0/ckp_0.3.0_linux_amd64.tar.gz"
      sha256 "42def17f63c1ddcc878ebfea642fe49e1dd23eac295f512e9a35e0d8313a76a9"
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/elhmn/ckp/releases/download/v0.3.0/ckp_0.3.0_linux_arm64.tar.gz"
      sha256 "d7c0582ac1910e21a1e7665e146d1cb9714d846b78834246649671a0bb9f2b5c"
    end
  end

  def install
    bin.install "ckp"
  end
end
