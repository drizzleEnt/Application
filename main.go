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
	w     fyne.Window
	a     fyne.App
	id    int
	chats ChatInfo
)

func main() {
	a = app.New()
	w = a.NewWindow("IWeather")

	w.Resize(fyne.NewSize(800, 500))
	w.SetMaster()
	w.CenterOnScreen()

	lgnLog := widget.NewEntry()
	pswLog := widget.NewEntry()
	pswLog.Password = true
	chatsContent := container.NewVBox()
	chatsContent.Hide()

	loginContent := container.NewVBox()
	loginBtn := widget.NewButton("Enter", func() {
		lgn := lgnLog.Text
		if lgn == "" {
			pswLog.SetText("")
		}

		pswd := pswLog.Text
		if pswd == "" {
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
			chatsContent.Show()
			res, err := GetChats()
			if err != nil {
				dialog.ShowInformation("Error", "failed get information about chats", w)
				return
			}
			chats = *res
		}
	})

	chatLog := widget.NewMultiLineEntry()
	input := widget.NewEntry()

	loginContent.Add(lgnLog)
	loginContent.Add(pswLog)
	loginContent.Add(loginBtn)

	chatsContent.Add(chatLog)
	chatsContent.Add(input)

	content := container.NewVBox(
		loginContent,
		chatsContent,
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

func ReceiveTokens(lgn string, pwd string) error {
	url := "http://localhost:8080/Login"
	method := http.MethodPost

	str := fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, lgn, pwd)

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

	variables := make(map[string]interface{}, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&variables)
	if err != nil {
		return err
	}
	if !variables["success"].(bool) {
		return fmt.Errorf("Success false")
	}

	tkn := variables["refresh_token"]
	id = int(variables["user_id"].(float64))
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

func GetChats() (*ChatInfo, error) {
	url := "http://localhost:9091/chats"
	method := http.MethodPost
	fmt.Printf("id: %v\n", id)
	str := fmt.Sprintf(`{
		"user_id": %v
	}`, id)

	payload := strings.NewReader(str)

	cl := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := cl.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var ch ChatInfo
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ch)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	if !ch.Success {
		return nil, fmt.Errorf("failed get chats")
	}

	fmt.Printf("chats: %v\n", ch)
	return &ch, nil
}
