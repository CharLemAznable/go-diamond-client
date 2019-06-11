package diamond_client

import (
    . "github.com/CharLemAznable/gokits"
    "os"
)

var Separator = string(os.PathSeparator)

var DataDir = "config-data"
var SnapshotDir = "snapshot"

var DefaultGroupName = "DEFAULT_GROUP"

var PollingInterval = 15              // seconds
var OnceTimeout = 3000                // milli seconds
var RecvWaitTimeout = OnceTimeout * 5 // milli seconds
var ConnTimeout = 3000                // milli seconds

var DiamondStoneExt = ".diamond"

func init() {
    LOG.AddFilter("stdout", DEBUG, NewConsoleLogWriter())
}
