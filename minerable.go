package diamond_client

import (
    . "github.com/CharLemAznable/gokits"
)

type Minerable interface {
    GetStone(group string, dataId string) (string, error)
    GetStoneOrDefault(group string, dataId string, defaultValue string) (string, error)

    GetProperties(group string, dataId string) (*Properties, error)
    GetPropertiesOfDefaultGroup(dataId string) (*Properties, error)

    GetMiner(group string, dataId string) (Minerable, error)
    GetMinerOfDefaultGroup(dataId string) (Minerable, error)

    GetString(key string) string
    GetStringOrDefault(key string, defaultValue string) string
}
