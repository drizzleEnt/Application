package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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

	lgnLog := widget.NewEntry()
	pswLog := widget.NewEntry()

	//chatLog := widget.NewMultiLineEntry()
	//input := widget.NewEntry()
	loginContent := container.NewVBox()
	loginBtn := widget.NewButton("Enter", func() {
		lgn := lgnLog.Text
		if lgn != "" {
			//send msg
			//input.SetText("")
		} else {
			pswLog.SetText("")
			dialog.ShowInformation("Error", "Enter login", w)
		}
		pswd := pswLog.Text
		if pswd != "" {
			//send msg
			//input.SetText("")
		} else {
			pswLog.SetText("")
			dialog.ShowInformation("Error", "Enter password", w)
		}

		err := ReceiveTokens(lgn, pswd)
		if err != nil {
			dialog.ShowInformation("Error", "Wrong login or password", w)
			fmt.Println(err.Error())
			pswLog.SetText("")
		} else {
			lgnLog.SetText("")
			lgnLog.Hide()
			pswLog.SetText("")
			pswLog.Hide()
			loginContent.Hide()
		}
	})

	loginContent.Add(lgnLog)
	loginContent.Add(pswLog)
	loginContent.Add(loginBtn)

	content := container.NewVBox(
		loginContent,
		//chatLog,
		//input,
		//lgnLog,
		//pswLog,
		//loginBtn,
	)
	w.SetContent(content)

	w.ShowAndRun()
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

func ReceiveTokens(lgn string, pswd string) error {
	url := "http://localhost:8080/Login"
	method := http.MethodPost

	str := fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, lgn, pswd)

	payload := strings.NewReader(str)

	cl := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := cl.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	vars := make(map[string]interface{}, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&vars)
	if err != nil {
		return err
	}
	if !vars["success"].(bool) {
		return fmt.Errorf("Success false")
	}

	tkn := vars["refresh_token"]
	err = CreateRefresh(tkn.(string))
	if err != nil {
		return err
	}
	return nil
}

func CreateRefresh(refreshToken string) error {
	file, err := os.Create("bin/refresh_token.txt")
	if err != nil {
		return fmt.Errorf("failed to create or open file: %v", err)
	}
	defer file.Close()
	_, err = file.Write([]byte(refreshToken))
	if err != nil {
		return fmt.Errorf("failed to write in file: %v", err)
	}

	return nil
}

func CreateAccess(accessToken string) error {
	file, err := os.Create("bin/access_token.txt")
	if err != nil {
		return fmt.Errorf("failed to create or open file: %v", err)
	}
	defer file.Close()
	_, err = file.Write([]byte(accessToken))
	if err != nil {
		return fmt.Errorf("failed to write in file: %v", err)
	}

	return nil
}
