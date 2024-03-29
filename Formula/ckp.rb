# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ckp < Formula
  desc ""
  homepage "https://github.com/elhmn/ckp"
  version "0.19.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/elhmn/ckp/releases/download/v0.19.0/ckp_0.19.0_darwin_arm64.tar.gz"
      sha256 "c19a09be952b6f7ad3b6227ae1f7f5fab52308587914d0bf0bb1be35ff42cb26"

      def install
        bin.install "ckp"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.19.0/ckp_0.19.0_darwin_amd64.tar.gz"
      sha256 "506169508d593429054e9cd1a223bd57740bd9f7b1fecad54cde2d083aea76a2"

      def install
        bin.install "ckp"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/elhmn/ckp/releases/download/v0.19.0/ckp_0.19.0_linux_amd64.tar.gz"
      sha256 "d8b107ae6f132fae74f424e770b6271af8b7e655ca4a4a95cc2ee5bfe37a7ef0"

      def install
        bin.install "ckp"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/elhmn/ckp/releases/download/v0.19.0/ckp_0.19.0_linux_arm64.tar.gz"
      sha256 "c29a70c732ea401c7137d24b7cbdf3ff2fc39fb6be782040928f596cfde369b4"

      def install
        bin.install "ckp"
      end
    end
  end
end
