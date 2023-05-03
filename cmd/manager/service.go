package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

func createServiceAndStartB7s(binPath string) {
	platform := runtime.GOOS

	switch platform {
	case "linux":
		createLinuxService(binPath)
	case "darwin":
		createMacOSService(binPath)
	case "windows":
		createWindowsService(binPath)
	default:
		log.Fatalf("Unsupported platform: %s", platform)
	}
}

func createLinuxService(binPath string) {
	serviceContent := `[Unit]
Description=Blockless b7s CLI Service
After=network.target

[Service]
ExecStart=%s/bls
Restart=always
User=%s

[Install]
WantedBy=multi-user.target
`

	serviceFilePath := "/etc/systemd/system/b7s.service"
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	serviceContent = fmt.Sprintf(serviceContent, binPath, usr.Username)

	err = os.WriteFile(serviceFilePath, []byte(serviceContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("systemctl", "enable", "b7s")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("systemctl", "start", "b7s")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("b7s service created and started on Linux.")
}

func createMacOSService(binPath string) {
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.blockless.b7s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s/bls</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
</dict>
</plist>
`

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	launchAgentsPath := filepath.Join(usr.HomeDir, "Library", "LaunchAgents")
	os.MkdirAll(launchAgentsPath, os.ModePerm)

	plistFilePath := filepath.Join(launchAgentsPath, "com.blockless.b7s.plist")
	plistContent = fmt.Sprintf(plistContent, binPath)

	err = os.WriteFile(plistFilePath, []byte(plistContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("launchctl", "load", plistFilePath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("b7s service created and started on macOS.")
}

func createWindowsService(binPath string) {
	log.Fatal("Creating and starting a service on Windows is not supported in this code.")
}

func removeServiceAndB7s() {
	platform := runtime.GOOS

	switch platform {
	case "linux":
		removeLinuxService()
	case "darwin":
		removeMacOSService()
	case "windows":
		removeWindowsService()
	default:
		log.Fatalf("Unsupported platform: %s", platform)
	}

	removeB7s()
}

func removeLinuxService() {
	cmd := exec.Command("systemctl", "disable", "b7s")
	err := cmd.Run()
	if err != nil {
		log.Println("Error disabling b7s service:", err)
	}

	cmd = exec.Command("systemctl", "stop", "b7s")
	err = cmd.Run()
	if err != nil {
		log.Println("Error stopping b7s service:", err)
	}

	serviceFilePath := "/etc/systemd/system/b7s.service"
	err = os.Remove(serviceFilePath)
	if err != nil {
		log.Println("Error removing b7s service file:", err)
	}

	log.Println("b7s service removed on Linux.")
}

func removeMacOSService() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	plistFilePath := filepath.Join(usr.HomeDir, "Library", "LaunchAgents", "com.blockless.b7s.plist")

	cmd := exec.Command("launchctl", "unload", plistFilePath)
	err = cmd.Run()
	if err != nil {
		log.Println("Error unloading b7s service:", err)
	}

	err = os.Remove(plistFilePath)
	if err != nil {
		log.Println("Error removing b7s service file:", err)
	}

	log.Println("b7s service removed on macOS.")
}

func removeWindowsService() {
	log.Fatal("Removing a service on Windows is not supported in this code.")
}