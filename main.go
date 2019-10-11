//Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package main

import (
	"fmt"
	"github.com/romanornr/viacoin-node-setup/binaries"
	"github.com/romanornr/viacoin-node-setup/blockbook"
	"os/exec"
	"sync"
	"time"

	"github.com/btcsuite/btcd/rpcclient"

	log "github.com/Sirupsen/logrus"
	"github.com/romanornr/viacoin-node-setup/client"
)

// latest block
var tip int

func main() {
	binaries.Download()
	binaries.Untar()
	syncNode()

}

func syncNode() {

	var wg sync.WaitGroup

	rpcclient := client.GetInstance()
	blockcount, _ := getBlockCount(rpcclient)
	tip = blockbook.GetStatus().Backend.Blocks // get the latest block with the blockbook api

	// blocks added in the sync progress. Close Viacoind and these blocks will be saved
	// without the need to resync
	blocksToAddInDisk := 500000 + blockcount

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
	log.Warnf("blocks to add to disk: %d\n", blocksToAddInDisk)
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
	var wg sync.WaitGroup
	wg.Add(1)
	blockcount, err := rpcclient.GetBlockCount()
	if err != nil {
		log.Errorf("%s\n", err)
		// viacoind could not have started yet and it's loadin block index
		// When this happens we need to make sure it started
		log.Warn("Viacoind is loading....")
		go func() {
			exec.Command("/bin/sh", "start.sh").Run()
			wg.Done()
		}()
		wg.Wait()
		time.Sleep(time.Second * 15)
	}

	// wait for Daemon to come online, otherwise keep blocking
	blockcount, _ = waitForDaemon(rpcclient)

	return blockcount, nil
}

// imagine the tip is equal the blockcount
// We dont' want viacoind to stop running
// instead do a return to escape the function
func SyncCompleted(blockcount int64) bool {
	//tip := blockcount
	if blockcount >= int64(tip) {
		log.Info("Chain fully synced")
		return true
	}
	return false
}

func waitForDaemon(rpcclient *rpcclient.Client) (int64, error) {
	//keep checking every 10 seconds
	for {
		blockcount, err := rpcclient.GetBlockCount()
		if err != nil {
			log.Warn(err)
			time.Sleep(time.Second * 10)
		}

		if err == nil && blockcount > 1 {
			break
		}
	}
	blockcount, _ := rpcclient.GetBlockCount()
	return blockcount, nil
}
