package main

type ChatInfo struct {
	Chat    []Chat `json:"chats"`
	Success bool   `json:"success"`
}

type Chat struct {
	ChatId   string `json:"UID"`
	ChatName string `json:"ChatName"`
}
