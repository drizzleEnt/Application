package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

var (
	w fyne.Window
	a fyne.App
)

func main() {
	a = app.New()
	w = a.NewWindow("IWeather")

	w.Resize(fyne.NewSize(800, 500))
	w.SetMaster()
	w.CenterOnScreen()

	chatLog := widget.NewMultiLineEntry()
	input := widget.NewEntry()

	sendBtn := widget.NewButton("send", func() {
		msg := input.Text
		if msg != "" {
			//send msg
			input.SetText("")
		} else {
			dialog.ShowInformation("Error", "Please select a user and enter a message", w)
		}

	})

	content := container.NewVBox(
		chatLog,
		input,
		sendBtn,
	)
	w.SetContent(content)

	go systray.Run(onReady, onExit)

	//w.ShowAndRun()
	a.Run()
}

func onReady() {
	//systray.SetTitle("Chat")

	mShow := systray.AddMenuItemCheckbox("Show", "Show chat window", false) //AddMenuItem("Show", "Show chat window")
	mQuit := systray.AddMenuItemCheckbox("Quit", "Qiut app", false)         //AddMenuItem("Quit", "Qiut app")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				w.Show()
			case <-mQuit.ClickedCh:
				systray.Quit()
				a.Quit()
				return
			}
		}
	}()
}

func onExit() {

}

func GetTemperature() ([]string, []string, error) {
	s := weatherService{}
	weather, err := s.GetWeather()
	if err != nil {
		return nil, nil, err
	}
	count := len(weather.Hourly.Temperatures[24:74])
	temps := make([]string, 0, count)
	hours := make([]string, 0, count)
	for i := 0; i < count; i++ {
		temps = append(temps, fmt.Sprintf("%.1f", weather.Hourly.Temperatures[i]))
		hours = append(hours, weather.Hourly.Time[i])
	}

	return temps, hours, nil
}
