package main

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
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
	syncNode()

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

func syncNode() {

	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("Restarting viacoind")

	go func() {
		defer wg.Done()
		//exec.Command("/bin/sh", "start.sh").Run() // blocking. Needs fix
		fmt.Println("givging 5 minutes to start up")
	}()

	time.Sleep(time.Second * 2)

	rpcclient := client.GetInstance()

	_, err := rpcclient.GetBlockChainInfo()
	if err != nil {
		log.Errorf("%s\n", err) // viacoind could not have started yet and it's loadin block index
		// When this happens we need to make sure it started
	}

	blockcount, err := rpcclient.GetBlockCount()
	if err != nil {
		fmt.Errorf("getting blockcount failed: %s \n", err)
	}
	log.Infof("viacoin blockcount %d \n", blockcount)

	// blocks added in the sync progress. Close Viacoind and these blocks will be saved
	// without the need to resync
	blocksToAddInDisk := 100000 + blockcount
	tip := 6834361

	for {
		blockcount, err := rpcclient.GetBlockCount()
		if err != nil {
			fmt.Errorf("getting blockcount failed: %s \n", err)
		}

		completion := float32(100) / float32(tip) * float32(blockcount)
		log.Infof("viacoin blockcount %d: synced %.2f %s", blockcount, completion, "%")
		time.Sleep(time.Second * 10)

		// if enough blocks got synced, close viacoind
		if blockcount > blocksToAddInDisk {
			break
		}
	}

	if SyncCompleted(blockcount) {
		return // return to block stop.sh from executing
	}

	log.Info("Stopping Viacoind")
	exec.Command("/bin/sh", "stop.sh").Run()
}

// imagine the tip is equal the blockcount
// We dont' want viacoind to stop running
// instead do a return to escape the function
func SyncCompleted(blockcount int64) bool {
	//tip := 6834361
	tip := blockcount
	if blockcount >= int64(tip) {
		log.Info("Chain fully synced")
		return true
	}
	return false
}
