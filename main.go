package main

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/romanornr/viacoin-node-setup/client"

	"github.com/cavaliercoder/grab"
)

func main() {

	homepath := ""

	if runtime.GOOS == "linux" {
		homepath = "~/.viacoin2"
	}

	fmt.Println(homepath)
	//DownloadBinaries()
	untar()
	sync()

}

// Download binaries from github
func DownloadBinaries() {

	url := "https://github.com/viacoin/viacoin/releases/download/v0.16.3/viacoin-0.16.3-x86_64-linux-gnu.tar.gz"

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", url)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("Downloading Viacoin binaries (%.2f%%)\n",
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	r, err := zip.OpenReader(resp.Filename)
	if err != nil {
		//return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
}

func untar() {
	exec.Command("/bin/sh", "untar.sh").Run()
}

func sync() {
	fmt.Println("Restarting viacoind")
	exec.Command("/bin/sh", "start.sh").Run() // blocking. Needs fix
	fmt.Println("givging 5 minutes to start up")
	time.Sleep(time.Minute * 5)

	rpcclient := client.GetInstance()

	blocktemplate, err := rpcclient.GetBlockChainInfo()
	if err != nil {
		log.Errorf("%s\n")
	}
	log.Info("chain: %s \n", blocktemplate.Chain)

	blockcount, err := rpcclient.GetBlockCount()
	if err != nil {
		fmt.Errorf("getting blockcount failed: %s \n", err)
	}
	log.Infof("viacoin blockcount %d \n", blockcount)

	time.Sleep(time.Minute * 5)

	log.Infof("viacoin blockcount %d \n", blockcount)

	exec.Command("/bin/sh", "stop.sh").Run()
}
