
package display

import (
	"time"
	"fmt"
	"log"
	"strings"
	"github.com/gen2brain/beeep"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)
var period_work int = 1
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
	startPeriod(clock, period)

	notes := widget.NewMultiLineEntry()
	notes.SetPlaceHolder("Enter notes...")
	clock_container := container.NewVBox(clock, widget.NewButton("Skip", func() {
		remaining = time.Until(time.Now())
	}))
	notesCont := container.NewVBox(widget.NewRichTextWithText("Notes for period"), notes)/*widget.NewTextSegment("Save", func() {
		log.Println("Notes were:", notes.Text)
		notes.SetText("")
	}))*/

	content := container.New(layout.NewHBoxLayout(), clock_container, layout.NewSpacer(), notesCont)
	updateTime(clock)


	w.SetContent(content)
	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
			if remaining.Seconds() < 1 {
				last_period := period
				length := 0

				if period == "work" {
					period = "rest"
					length = period_rest
				} else if period == "rest" {
					period = "work"
					length = period_work
				} else if period == "lunch" {
					period = "work"
					length = period_work
				}
				notify(fmt.Sprintf("%s is up! Next cycle is %d minutes", strings.Title(last_period), length))
				log.Println("Notes were:", notes.Text)
				notes.SetText("")
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
