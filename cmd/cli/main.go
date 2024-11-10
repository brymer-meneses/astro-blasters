package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"space-shooter/client"
	"space-shooter/client/config"
	"space-shooter/server"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "cli"}

	// Server command
	{
		var port int
		serverCmd := &cobra.Command{
			Use:   "server",
			Short: "Run the server",
			Run: func(cmd *cobra.Command, args []string) {

				var stderr bytes.Buffer

				build := exec.Command("go", "build", "-o", "server/static/game.wasm", "client/wasm/main.go")
				build.Stderr = &stderr
				build.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")

				if err := build.Run(); err != nil {
					fmt.Println(stderr.String())
					os.Exit(1)
				}

				if _, err := os.Stat("server/static/wasm_exec.js"); errors.Is(err, os.ErrNotExist) {
					output, err := exec.Command("go", "env", "GOROOT").Output()
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					goroot := strings.TrimSuffix(string(output), "\n")

					wasmExecPath := path.Join(goroot, "misc", "wasm", "wasm_exec.js")
					os.Link(wasmExecPath, "server/static/wasm_exec.js")

					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}

				server := server.NewServer()
				if err := server.Start(port); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			},
		}
		serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")

		rootCmd.AddCommand(serverCmd)
	}

	// Client command
	{
		var port int
		var address string
		var secure bool
		clientCmd := &cobra.Command{
			Use:   "client",
			Short: "Run the native client",
			Run: func(cmd *cobra.Command, args []string) {
				protocol := "ws"
				if secure {
					protocol = "wss"
				}

				url := fmt.Sprintf("%s://%s:%d/play/ws", protocol, address, port)
				config := config.ClientConfig{
					ScreenWidth:        1080,
					ScreenHeight:       720,
					ServerWebsocketURL: url,
				}

				app := client.NewApp(&config)
				if err := app.Run(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			},
		}

		clientCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port of the server")
		clientCmd.Flags().StringVarP(&address, "address", "a", "", "Address of the server")
		clientCmd.Flags().BoolVarP(&secure, "secure", "s", false, "Whether to use WSS")
		clientCmd.MarkFlagRequired("port")
		clientCmd.MarkFlagRequired("address")

		rootCmd.AddCommand(clientCmd)
	}

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
