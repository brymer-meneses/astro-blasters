//go:build wasm

package main

import (
	"fmt"
	"log"
	"space-shooter/client"
	"space-shooter/client/config"
	"strings"
	"syscall/js"
)

func getServerUrl() (serverUrl string, isSecure bool) {
	href := js.Global().Get("window").Get("location").Get("href").String()
	href = strings.TrimSuffix(href, "/")

	isHttps := strings.HasPrefix(href, "https://")
	if isHttps {
		href = strings.TrimPrefix(href, "https://")
	} else {
		href = strings.TrimPrefix(href, "http://")
	}

	return href, isHttps
}

func main() {
	serverUrl, _ := getServerUrl()
	serverWebsocketUrl := fmt.Sprintf("ws://%s/play/ws", serverUrl)

	config := config.ClientConfig{
		ScreenWidth:        1080,
		ScreenHeight:       720,
		ServerWebsocketURL: serverWebsocketUrl,
	}

	app := client.NewApp(&config)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
