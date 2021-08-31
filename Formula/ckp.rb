# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ckp < Formula
  desc ""
  homepage "https://github.com/elhmn/ckp"
  version "0.4.0"
  bottle :unneeded

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.4.0/ckp_0.4.0_darwin_amd64.tar.gz"
      sha256 "b8083b1a567a779b1720c2e86d6722d5971d0e9245583e35bb1e3fd16328327c"
    end
    if Hardware::CPU.arm?
      url "https://github.com/elhmn/ckp/releases/download/v0.4.0/ckp_0.4.0_darwin_arm64.tar.gz"
      sha256 "0a3c1559c82bb40338bea6ede5b9fa410b5a269751fcd8b51a65eff153024631"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.4.0/ckp_0.4.0_linux_amd64.tar.gz"
      sha256 "aa47973cea82353182dc39663fd0c216a82bad9c7616263dd91c86340f93dffe"
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/elhmn/ckp/releases/download/v0.4.0/ckp_0.4.0_linux_arm64.tar.gz"
      sha256 "ebc4938a0a9003ec9fdd1c3ad7930ed8ef006ff71ef24cf6f14ae224196ad2db"
    end
  end

  def install
    bin.install "ckp"
  end
end
