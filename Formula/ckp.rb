# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ckp < Formula
  desc ""
  homepage "https://github.com/elhmn/ckp"
  version "0.2.34"
  bottle :unneeded

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.2.34/ckp_0.2.34_darwin_amd64.tar.gz"
      sha256 "1c3776fd08e14fb20ebe8264379bf88abc4007fbde15fd243972acce6dbb5614"
    end
    if Hardware::CPU.arm?
      url "https://github.com/elhmn/ckp/releases/download/v0.2.34/ckp_0.2.34_darwin_arm64.tar.gz"
      sha256 "658fb5f7b9b7980d45c11c9360f95c841b3a71d8ed5d9d28140bc476259bf624"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.2.34/ckp_0.2.34_linux_amd64.tar.gz"
      sha256 "9a25c816b849cb6b5d55d35e9e03d725101d2e60cdebd6148379e61cb5b013e5"
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/elhmn/ckp/releases/download/v0.2.34/ckp_0.2.34_linux_arm64.tar.gz"
      sha256 "11402682ebc7bafcb6579904ca276f0cb76bf4595417828975603a53c98a7b29"
    end
  end

  def install
    bin.install "ckp"
  end
end
