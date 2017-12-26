package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/spf13/cobra"
)

var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "Show the live player panel",
	RunE:  player,
}

func player(cmd *cobra.Command, args []string) error {
	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	songList := ui.NewList()
	songList.Border = false
	songList.X = 0
	songList.Y = 0
	songList.Height = 2
	songList.Width = 40

	ctlList := ui.NewList()
	ctlList.Border = false
	ctlList.X = 0
	ctlList.Y = 3
	ctlList.Height = 2
	ctlList.Width = 40

	currPosLabel := ui.NewPar("")
	currPosLabel.X = 0
	currPosLabel.Y = 6
	currPosLabel.Width = 6
	currPosLabel.Border = false

	progressGauge := ui.NewGauge()
	progressGauge.LabelAlign = ui.AlignCenter
	progressGauge.Height = 2
	progressGauge.Y = 6
	progressGauge.X = 6
	progressGauge.Width = 30
	progressGauge.Border = false
	progressGauge.Label = ""
	progressGauge.Percent = 0
	progressGauge.PaddingBottom = 1

	totalSecLabel := ui.NewPar("")
	totalSecLabel.X = 38
	totalSecLabel.Y = 6
	totalSecLabel.Width = 6
	totalSecLabel.Border = false

	volGauge := ui.NewGauge()
	volGauge.LabelAlign = ui.AlignCenter
	volGauge.Height = 2
	volGauge.Y = 8
	volGauge.X = 0
	volGauge.Width = 44
	volGauge.Border = false
	volGauge.BarColor = ui.ColorBlue
	volGauge.PaddingBottom = 1

	helpLabel := ui.NewPar("Press q - quit, p - play/pause, l/h - next/previous track, j/k - vol up/down, s - shuffle, r - repeat.")
	helpLabel.X = 0
	helpLabel.Y = 10
	helpLabel.Width = 40
	helpLabel.Height = 5
	helpLabel.Border = false
	helpLabel.WrapLength = 40

	draw := func() {
		ui.Render(
			songList,
			ctlList,
			currPosLabel,
			progressGauge,
			totalSecLabel,
			volGauge,
			helpLabel,
		)
	}

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/p", func(ui.Event) {
		state, err := client.PlayerState()
		if err != nil {
			quitAndFatal(err)
		}
		if state.Playing {
			if err := pause(pauseCmd, []string{}); err != nil {
				quitAndFatal(err)
			}
		} else {
			if err := play(playCmd, []string{}); err != nil {
				quitAndFatal(err)
			}
		}
	})

	ui.Handle("/sys/kbd/l", func(ui.Event) {
		if err := next(nextCmd, []string{}); err != nil {
			quitAndFatal(err)
		}
	})

	ui.Handle("/sys/kbd/h", func(ui.Event) {
		if err := prev(prevCmd, []string{}); err != nil {
			quitAndFatal(err)
		}
	})

	ui.Handle("/sys/kbd/j", func(ui.Event) {
		if err := vol(volCmd, []string{"up"}); err != nil {
			quitAndFatal(err)
		}
	})

	ui.Handle("/sys/kbd/k", func(ui.Event) {
		if err := vol(volCmd, []string{"down"}); err != nil {
			quitAndFatal(err)
		}
	})

	ui.Handle("/sys/kbd/s", func(ui.Event) {
		if err := shuffle(shuffleCmd, []string{}); err != nil {
			quitAndFatal(err)
		}
	})

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		if err := repeat(repeatCmd, []string{}); err != nil {
			quitAndFatal(err)
		}
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		state, err := client.PlayerState()
		if err != nil {
			quitAndFatal(err)
		}

		volGauge.Percent = state.Device.Volume
		volGauge.Label = "Volume {{percent}}%"

		shuffleState := "off"
		if state.ShuffleState {
			shuffleState = "on"
		}

		ctlList.Items = []string{
			"Shuffle " + shuffleState,
			"Repeat " + state.RepeatState,
		}

		if state.Playing {
			progressGauge.Label = "Playing"
			progressGauge.BarColor = ui.ColorGreen
		} else {
			progressGauge.Label = "Paused"
			progressGauge.BarColor = ui.ColorRed
		}

		if state.Item != nil {
			var artists []string
			for _, a := range state.Item.Artists {
				artists = append(artists, a.Name)
			}
			songList.Items = []string{
				state.Item.Name,
				fmt.Sprintf("%s - %s", strings.Join(artists, ", "), state.Item.Album.Name),
			}

			currPosLabel.Text = durationToStr(state.Progress)
			totalSecLabel.Text = durationToStr(state.Item.Duration)

			progressGauge.Percent = progressPercent(state.Progress, state.Item.Duration)
		}

		draw()
	})

	ui.Loop()

	return nil
}

func quitAndFatal(err error) {
	ui.StopLoop()
	ui.Close()
	log.Fatal(err)
}

func durationToStr(d int) string {
	sec := roundToSec(d)
	minInt := int(sec / 60)
	secInt := int(sec - float64(minInt*60))
	return fmt.Sprintf("%d:%02d", minInt, secInt)
}

func roundToSec(d int) float64 {
	dur := (time.Duration(d) * time.Millisecond).Round(time.Second)
	return dur.Seconds()
}

func progressPercent(progress, total int) int {
	return int(roundToSec(progress) / roundToSec(total) * 100.0)
}
