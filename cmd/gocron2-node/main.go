// Command gocron2-node
package main

import (
	"flag"
	"os"
	"runtime"
	"strings"

	"github.com/ouqiang/goutil"
	log "github.com/sirupsen/logrus"
	"github.com/up9cloud/gocron2/internal/modules/rpc/auth"
	"github.com/up9cloud/gocron2/internal/modules/rpc/server"
	"github.com/up9cloud/gocron2/internal/modules/utils"
)

var (
	AppVersion, BuildDate, GitCommit string
)

func main() {
	var serverAddr string
	var allowRoot bool
	var version bool
	var CAFile string
	var certFile string
	var keyFile string
	var enableTLS bool
	var logLevel string
	flag.BoolVar(&allowRoot, "allow-root", false, "./gocron2-node -allow-root")
	flag.StringVar(&serverAddr, "s", "0.0.0.0:5921", "./gocron2-node -s ip:port")
	flag.BoolVar(&version, "v", false, "./gocron2-node -v")
	flag.BoolVar(&enableTLS, "enable-tls", false, "./gocron2-node -enable-tls")
	flag.StringVar(&CAFile, "ca-file", "", "./gocron2-node -ca-file path")
	flag.StringVar(&certFile, "cert-file", "", "./gocron2-node -cert-file path")
	flag.StringVar(&keyFile, "key-file", "", "./gocron2-node -key-file path")
	flag.StringVar(&logLevel, "log-level", "info", "-log-level error")
	flag.Parse()
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	if version {
		goutil.PrintAppVersion(AppVersion, GitCommit, BuildDate)
		return
	}

	if enableTLS {
		if !utils.FileExist(CAFile) {
			log.Fatalf("failed to read ca cert file: %s", CAFile)
		}
		if !utils.FileExist(certFile) {
			log.Fatalf("failed to read server cert file: %s", certFile)
			return
		}
		if !utils.FileExist(keyFile) {
			log.Fatalf("failed to read server key file: %s", keyFile)
			return
		}
	}

	certificate := auth.Certificate{
		CAFile:   strings.TrimSpace(CAFile),
		CertFile: strings.TrimSpace(certFile),
		KeyFile:  strings.TrimSpace(keyFile),
	}

	if runtime.GOOS != "windows" && os.Getuid() == 0 && !allowRoot {
		log.Fatal("Do not run gocron2-node as root user")
		return
	}

	server.Start(serverAddr, enableTLS, certificate)
}
