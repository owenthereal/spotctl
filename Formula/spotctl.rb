class Spotctl < Formula
  desc "A CLI to Spotify"
  homepage "https://github.com/jingweno/spotctl"
  version "1.0.1"
  sha256 "372a610e31516d179e208f8dc60c97004be6d7c8831583ad99fd83d5c0daefe0"
  url "https://github.com/jingweno/spotctl/releases/download/v1.0.1/darwin-amd64-1.0.1.tar.gz"

  def install
    bin.install "spotctl"
  end
end
