//Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package blockbookjson

import "time"

type Status struct {
	Blockbook struct {
		Coin            string    `json:"coin"`
		Host            string    `json:"host"`
		Version         string    `json:"version"`
		GitCommit       string    `json:"gitCommit"`
		BuildTime       time.Time `json:"buildTime"`
		SyncMode        bool      `json:"syncMode"`
		InitialSync     bool      `json:"initialSync"`
		InSync          bool      `json:"inSync"`
		BestHeight      int       `json:"bestHeight"`
		LastBlockTime   time.Time `json:"lastBlockTime"`
		InSyncMempool   bool      `json:"inSyncMempool"`
		LastMempoolTime time.Time `json:"lastMempoolTime"`
		MempoolSize     int       `json:"mempoolSize"`
		Decimals        int       `json:"decimals"`
		DbSize          int64     `json:"dbSize"`
		About           string    `json:"about"`
	} `json:"blockbook"`
	Backend struct {
		Chain           string `json:"chain"`
		Blocks          int    `json:"blocks"`
		Headers         int    `json:"headers"`
		BestBlockHash   string `json:"bestBlockHash"`
		Difficulty      string `json:"difficulty"`
		Version         string `json:"version"`
		Subversion      string `json:"subversion"`
		ProtocolVersion string `json:"protocolVersion"`
	} `json:"backend"`
}
