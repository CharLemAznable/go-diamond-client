package diamond_client

import (
    "fmt"
    . "github.com/CharLemAznable/gokits"
)

type diamond struct {
    miner *Miner
}

var Diamond = &diamond{miner: NewMiner()}

func (d *diamond) GetStone(group string, dataId string) (string, error) {
    return d.miner.GetStone(group, dataId)
}

func (d *diamond) GetStoneOrDefault(group string, dataId string, defaultValue string) (string, error) {
    return d.miner.GetStoneOrDefault(group, dataId, defaultValue)
}

func (d *diamond) GetProperties(group string, dataId string) (*Properties, error) {
    return d.miner.GetProperties(group, dataId)
}

func (d *diamond) GetPropertiesOfDefaultGroup(dataId string) (*Properties, error) {
    return d.miner.GetPropertiesOfDefaultGroup(dataId)
}

func (d *diamond) GetString(key string) string {
    fmt.Printf("%+v", d.miner)
    return d.miner.GetString(key)
}

func (d *diamond) GetStringOrDefault(key string, defaultValue string) string {
    return d.miner.GetStringOrDefault(key, defaultValue)
}
