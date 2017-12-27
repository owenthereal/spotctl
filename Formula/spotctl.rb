class Spotctl < Formula
  desc "A CLI to Spotify"
  homepage "https://github.com/jingweno/spotctl"
  version "1.0.1"
  sha256 "a0276bd0c0fb65e7b24885f17d3547276dcf9b059b9507abe51edfa5f5388f89"
  url "https://github.com/jingweno/spotctl/releases/download/v1.0.1/darwin-amd64-1.0.1.tar.gz"

  def install
    bin.install "spotctl"
  end
end
