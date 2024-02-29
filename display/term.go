
package display

import (
	"time"
	"fmt"
	"log"
	"strings"
	"os"
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
var period_length float64 = 0

func Main(starting_period string) {
	period := starting_period
	a := app.New()
	w := a.NewWindow("Skjald")

	clock := widget.NewLabel("Period Timer")
	startPeriod(clock, period)
	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = remaining.Seconds()
	period_length = float64(remaining.Seconds())

	notes := widget.NewMultiLineEntry()
	notes.SetPlaceHolder("Enter notes...")
	clock_container := container.NewVBox(clock, progress, widget.NewButton("Skip", func() {
		log.Println("skipped")
		end_time = time.Now()
		remaining = time.Until(time.Now())
	}))
	notesCont := container.NewVBox(widget.NewRichTextWithText("Notes for period"), notes)

	content := container.New(layout.NewHBoxLayout(), clock_container, layout.NewSpacer(), notesCont)
	updateTime(clock)


	w.SetContent(content)
	go func() {
		for range time.Tick(time.Second) {
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
				notify(fmt.Sprintf("%s period is over! Next %s cycle is %d minutes", strings.Title(last_period), strings.Title(period), length))
				log.Println("Notes were:", notes.Text)
				write_notes(notes.Text)
				notes.SetText("")
				startPeriod(clock, period)
			}
			updateTime(clock)
			fmt.Println(period_length - remaining.Seconds())
			progress.SetValue(period_length - remaining.Seconds())
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

func write_notes(content string) {
	f, err := os.OpenFile("notes.md", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	header := strings.Join([]string{"\n#", start_time.Format("03:04 PM"), " - ", end_time.Format("03:04 PM"), "\n"}, "")
	f.WriteString(strings.Join([]string{header, content}, ""))

	f.Close()
}
