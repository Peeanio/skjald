
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
	"github.com/spf13/viper"
)
var start_time = time.Now()
var end_time = time.Now()
var remaining time.Duration = 0
var period_length float64 = 0

func Main(starting_period string) {
	period := starting_period
	a := app.New()
	w := a.NewWindow("Skjald")

	clock := widget.NewLabel("Period Timer")
	progress := widget.NewProgressBar()
	startPeriod(clock, period, progress)

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
					length = viper.GetInt("period_rest")
				} else if period == "rest" {
					period = "work"
					length = viper.GetInt("period_work")
				} else if period == "lunch" {
					period = "work"
					length = viper.GetInt("period_work")
				}
				notify(fmt.Sprintf("%s period is over! Next %s cycle is %d minutes", strings.Title(last_period), strings.Title(period), length))
				log.Println("Notes were:", notes.Text)
				write_notes(notes.Text)
				notes.SetText("")
				startPeriod(clock, period, progress)
			}
			updateTime(clock)
			progress.SetValue(period_length - remaining.Seconds())
		}
	}()
	w.ShowAndRun()

}

func startPeriod(clock *widget.Label, period string, progress *widget.ProgressBar){
	start_time = time.Now()
	if period == "work" {
		end_time = start_time.Add(time.Minute * time.Duration(viper.GetInt("period_work")))
	} else if period == "rest" {
		end_time = start_time.Add(time.Minute * time.Duration(viper.GetInt("period_rest")))
	} else if period == "lunch" {
		end_time = start_time.Add(time.Minute * time.Duration(viper.GetInt("period_lunch")))
	}
	remaining = time.Until(end_time)

	progress.Min = 0
	progress.Max = remaining.Seconds()
	period_length = float64(remaining.Seconds())
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
	f, err := os.OpenFile(time.Now().Format("2006_01_02_notes.md"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	header := strings.Join([]string{"\n#", start_time.Format("03:04 PM"), " - ", end_time.Format("03:04 PM"), "\n"}, "")
	f.WriteString(strings.Join([]string{header, content}, ""))

	f.Close()
}
