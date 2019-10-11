//Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package blockbook

import (
	"encoding/json"
	"fmt"
	"github.com/romanornr/viacoin-node-setup/blockbookjson"
	"log"
	"net/http"
	"time"
)

func GetStatus() *blockbookjson.Status {
	url := fmt.Sprintf("https://blockbook.viacoin.org/api/status")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("NewRequest error: %s\n", err)
	}

	client := &http.Client{
		Timeout:       time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error client.DO: %s\n", err)
	}

	defer resp.Body.Close()

	var status blockbookjson.Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Fatalf("error decoding blockbook status: %s\n", err)
	}
	return &status
}
