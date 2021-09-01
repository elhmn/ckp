# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ckp < Formula
  desc ""
  homepage "https://github.com/elhmn/ckp"
  version "0.16.0"
  bottle :unneeded

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.16.0/ckp_0.16.0_darwin_amd64.tar.gz"
      sha256 "20f9658fa2b8dd792e24df8448740cd78617ed2b1bfaf896c532a7b417037a25"
    end
    if Hardware::CPU.arm?
      url "https://github.com/elhmn/ckp/releases/download/v0.16.0/ckp_0.16.0_darwin_arm64.tar.gz"
      sha256 "2c606362da56404300bd76294f70d9b56e1371299e939cc9b557381bbaaee1e5"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.16.0/ckp_0.16.0_linux_amd64.tar.gz"
      sha256 "4436cea7ff42a15c0fa72139201059f1fff4c0e88a2b9808ae455a4054ab8893"
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/elhmn/ckp/releases/download/v0.16.0/ckp_0.16.0_linux_arm64.tar.gz"
      sha256 "32fe304afb6f18ae74c3c727ec3e768a25332538ee923cc05e04e1a5d74a7c21"
    end
  end

  def install
    bin.install "ckp"
  end
end
