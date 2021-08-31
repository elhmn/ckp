# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ckp < Formula
  desc ""
  homepage "https://github.com/elhmn/ckp"
  version "0.7.0"
  bottle :unneeded

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.7.0/ckp_0.7.0_darwin_amd64.tar.gz"
      sha256 "a6cdf6ce3e2efaea1ec265befe40af215d7a1d962fd3dc6f08d2bafb658f91e9"
    end
    if Hardware::CPU.arm?
      url "https://github.com/elhmn/ckp/releases/download/v0.7.0/ckp_0.7.0_darwin_arm64.tar.gz"
      sha256 "ecedd09b03edaccb839256d19d1bb5fd448f3baee41308c8ed2053ef3d801f5a"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.7.0/ckp_0.7.0_linux_amd64.tar.gz"
      sha256 "33faaab3221b9f460e9477b0b31a01c73c0fb56157ca4607e78ab3b91709615b"
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/elhmn/ckp/releases/download/v0.7.0/ckp_0.7.0_linux_arm64.tar.gz"
      sha256 "aa7f51f017d57dd76795928bfa32d293cc4247c09beb2433f344ec1f120237a0"
    end
  end

  def install
    bin.install "ckp"
  end
end
