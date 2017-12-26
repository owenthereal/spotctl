class Spotctl < Formula
  desc "A CLI to Spotify"
  homepage "https://github.com/jingweno/spotctl"
  version "1.0.0"
  sha256 "89e926f5c1e66f466ea5607e30bc6b120c6363ff0df407e497aa0b7af8151127"
  url "https://github.com/jingweno/spotctl/releases/download/v1.0.0/darwin-amd64-1.0.0.tar.gz"

  def install
    bin.install "spotctl"
  end
end
