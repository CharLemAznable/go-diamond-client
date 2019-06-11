package diamond_client

import (
    . "github.com/CharLemAznable/gokits"
)

type DefaultGroup interface {
    DefaultGroupName() string
}

type AbstractMiner struct {
    Minerable
    DefaultGroup
}

func newAbstractMiner(minerable Minerable, defaultGroup DefaultGroup) *AbstractMiner {
    return &AbstractMiner{Minerable: minerable, DefaultGroup: defaultGroup}
}

func (miner *AbstractMiner) GetStoneOrDefault(group string, dataId string, defaultValue string) (string, error) {
    stone, err := miner.GetStone(group, dataId)
    return Condition(nil == err, stone, defaultValue).(string), nil
}

func (miner *AbstractMiner) GetProperties(group string, dataId string) (*Properties, error) {
    stone, err := miner.GetStone(group, dataId)
    if nil != err {
        return nil, err
    }
    return ParseStoneToProperties(stone), nil
}

func (miner *AbstractMiner) GetPropertiesOfDefaultGroup(dataId string) (*Properties, error) {
    return miner.GetProperties(miner.DefaultGroupName(), dataId)
}

func (miner *AbstractMiner) GetMiner(group string, dataId string) (Minerable, error) {
    properties, err := miner.GetProperties(group, dataId)
    if nil != err {
        return nil, err
    }
    return NewPropertiesBasedMiner(properties), nil
}

func (miner *AbstractMiner) GetMinerOfDefaultGroup(dataId string) (Minerable, error) {
    return miner.GetMiner(miner.DefaultGroupName(), dataId)
}

func (miner *AbstractMiner) GetString(key string) string {
    stone, _ := miner.GetStone(miner.DefaultGroupName(), key)
    return stone
}

func (miner *AbstractMiner) GetStringOrDefault(key string, defaultValue string) string {
    stone, _ := miner.GetStoneOrDefault(miner.DefaultGroupName(), key, defaultValue)
    return stone
}
