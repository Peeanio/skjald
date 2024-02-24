
package display

import (
	"time"
	"strings"
	"github.com/gen2brain/beeep"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)
var period_work int = 25
var period_rest int = 5
var period_lunch int = 30
var start_time = time.Now()
var end_time = time.Now()
var remaining time.Duration = 0

func Main(starting_period string) {
	period := starting_period
	a := app.New()
	w := a.NewWindow("Skjald")

	clock := widget.NewLabel("Period Timer")
	startPeriod(clock, starting_period)
	updateTime(clock)


	w.SetContent(clock)
	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
			if remaining.Seconds() < 1 {
				notify("Times Up!")
				if period == "work" {
					period = "rest"
				} else if period == "rest" {
					period = "work"
				} else if period == "lunch" {
					period = "work"
				}
				startPeriod(clock, period)
			}
		}
	}()
	w.ShowAndRun()

}

func startPeriod(clock *widget.Label, period string){
	start_time = time.Now()
	if period == "work" {
		end_time = start_time.Add(time.Minute * time.Duration(period_work))
	} else if period == "rest" {
		end_time = start_time.Add(time.Minute * time.Duration(period_rest))
	} else if period == "lunch" {
		end_time = start_time.Add(time.Minute * time.Duration(period_lunch))
	}
	remaining = time.Until(end_time)

}

func updateTime(clock *widget.Label) {
	start_formatted := start_time.Format("Period Start: 03:04:05")
	end_formatted := end_time.Format("Period End: 03:04:05")
	remaining = time.Until(end_time)
	remaining_formatted := "Remaining: " + remaining.Truncate(time.Second).String()
	clock.SetText(strings.Join([]string{start_formatted, end_formatted, remaining_formatted}, "\n"))
}

func notify(text string) {
	err := beeep.Notify("Skjald says:", text, "assets/skjald.jpg")
	if err != nil {
		panic(err)
	}
}
