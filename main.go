package main

import (
	"fmt"
	"github.com/romanornr/viacoin-node-setup/binaries"
	"os/exec"
	"sync"
	"time"

	"github.com/btcsuite/btcd/rpcclient"

	log "github.com/Sirupsen/logrus"
	"github.com/romanornr/viacoin-node-setup/client"
)

func main() {

	// homepath := ""

	// if runtime.GOOS == "linux" {
	// 	homepath = "~/.viacoin2"
	// }

	// fmt.Println(homepath)
	//binaries.Download()
	binaries.Untar()
	syncNode()

}

func syncNode() {

	var wg sync.WaitGroup

	rpcclient := client.GetInstance()
	getBlockCount(rpcclient)

	blockcount, err := rpcclient.GetBlockCount()
	if err != nil {
		fmt.Errorf("getting blockcount failed: %s \n", err)
	}
	log.Infof("viacoin blockcount %d \n", blockcount)

	// blocks added in the sync progress. Close Viacoind and these blocks will be saved
	// without the need to resync
	blocksToAddInDisk := 500000 + blockcount
	tip := 6834361

	for {
		blockcount, err := getBlockCount(rpcclient)
		if err != nil {
			fmt.Errorf("getting blockcount failed: %s \n", err)
		}

		completion := float32(100) / float32(tip) * float32(blockcount)
		log.Infof("viacoin blockcount %d: synced %.2f %s", blockcount, completion, "%")
		time.Sleep(time.Second * 10)

		if SyncCompleted(blockcount) {
			return // return to block stop.sh from executing

		}

		// if enough blocks got synced, get out of the loop
		// and close viacoind
		if blockcount > blocksToAddInDisk {
			break
		}
	}

	log.Info("Stopping Viacoind & saving all blocks")
	//rpcclient.Shutdown()
	wg.Add(1)
	go func() {
		exec.Command("/bin/sh", "stop.sh").Run()
		wg.Done()
	}()
	time.Sleep(time.Second * 30) // gracefully shutdown
	syncNode()
}

// this function will get the blockcount but if it seems like
// the daemon is not running or crashed, it will restart the daemon and keep
// trying to make a connection and output the blockcount back
func getBlockCount(rpcclient *rpcclient.Client) (int64, error) {
	blockcount, err := rpcclient.GetBlockCount()
	if err != nil {
		fmt.Println("Viacoind was not running, starting it now")
		log.Errorf("%s\n", err)
		// viacoind could not have started yet and it's loadin block index
		// When this happens we need to make sure it started
		log.Warn("Viacoind will load block index & sync headers")
		log.Warn("This can take longer or shorter depending on the server hardware")
		go func() {
			exec.Command("/bin/sh", "start.sh").Run() // blocking. Needs fixtime.Sleep(time.Minute * 5)
		}()

		//keep checking every 10 seconds
		for {
			_, err := rpcclient.GetBlockCount()
			if err != nil {
				log.Warn(err)
				time.Sleep(time.Second * 10)
			}

			if err == nil {
				return blockcount, nil
			}
		}
	}
	return blockcount, nil
}

// imagine the tip is equal the blockcount
// We dont' want viacoind to stop running
// instead do a return to escape the function
func SyncCompleted(blockcount int64) bool {
	tip := 6834361
	//tip := blockcount
	if blockcount >= int64(tip) {
		log.Info("Chain fully synced")
		return true
	}
	return false
}
