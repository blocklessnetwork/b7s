package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	ma "github.com/multiformats/go-multiaddr"
)

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	listenF := pflag.IntP("listen", "l", 0, "wait for incoming connections")
	insecureF := pflag.BoolP("insecure", "i", false, "use an unencrypted connection")
	seedF := pflag.Int64P("seed", "s", 0, "set random seed for id generation")
	allowedPeerF := pflag.StringP("allowed-peer", "a", "", "allowed peer ID")
	pflag.Parse()

	if *listenF == 0 {
		logger.Fatal().Msg("Please provide a port to bind on with -l")
	}

	ha, err := makeBasicHost(*listenF, *insecureF, *seedF)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create host")
	}

	startListener(ctx, ha, *listenF, *insecureF, *allowedPeerF, logger)
	<-ctx.Done()
}

func makeBasicHost(listenPort int, insecure bool, randseed int64) (host.Host, error) {
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}

	if insecure {
		opts = append(opts, libp2p.NoSecurity)
	}

	return libp2p.New(opts...)
}

func getHostAddress(ha host.Host) string {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID().Pretty()))
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func startListener(ctx context.Context, ha host.Host, listenPort int, insecure bool, allowedPeer string, logger zerolog.Logger) {
	fullAddr := getHostAddress(ha)
	logger.Info().Msgf("I am %s", fullAddr)

	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		if allowedPeer != "" && s.Conn().RemotePeer().Pretty() != allowedPeer {
			logger.Info().Msg("Connection from disallowed peer")
			s.Reset()
			return
		}
		logger.Info().Msg("listener received new stream")
		if err := doEcho(s); err != nil {
			logger.Info().Err(err).Msg("Error in doEcho")
			s.Reset()
		} else {
			s.Close()
		}
	})

	logger.Info().Msg("listening for connections")

	if insecure {
		logger.Info().Msgf("Now run \"./echo -l %d -d %s --insecure\" on a different terminal", listenPort+1, fullAddr)
	} else {
		logger.Info().Msgf("Now run \"./echo -l %d -d %s\" on a different terminal", listenPort+1, fullAddr)
	}
}

func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	var msg Message
	err = json.Unmarshal([]byte(str), &msg)
	if err != nil {
		return err
	}

	baseURL := "https://github.com/blocklessnetwork/cli/releases/download"
	version := "0.0.46"

	switch msg.Type {
	case "install_bls":
		go func() {
			installBlsCLI(baseURL, version)
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			binPath := filepath.Join(usr.HomeDir, ".b7s", "bin")
			createServiceAndStartB7s(binPath)
		}()
	}

	_, err = s.Write([]byte(str))
	return err
}

func installBlsCLI(baseURL string, version string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	binPath := filepath.Join(usr.HomeDir, ".b7s", "bin")
	os.MkdirAll(binPath, os.ModePerm)

	arch := runtime.GOARCH
	platform := runtime.GOOS

	// maybe change this in ci
	if platform == "darwin" {
		platform = "macOS"
	}

	url := fmt.Sprintf("%s/%s/bls-%s-%s-blockless-cli.tar.gz", baseURL, version, platform, arch)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	archiveData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(archiveData))
	if err != nil {
		log.Fatal(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if header.Typeflag == tar.TypeReg {
			target := filepath.Join(binPath, "b7s")
			outFile, err := os.Create(target)
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatal(err)
			}

			if err := os.Chmod(target, 0755); err != nil {
				log.Fatal(err)
			}

			log.Printf("b7s CLI installed in %s", binPath)
			break
		}
	}
}

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
