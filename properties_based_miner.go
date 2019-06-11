package diamond_client

import (
    . "github.com/CharLemAznable/gokits"
)

type PropertiesBasedMiner struct {
    AbstractMiner
    properties *Properties
}

func NewPropertiesBasedMiner(properties *Properties) *PropertiesBasedMiner {
    miner := &PropertiesBasedMiner{properties: properties}
    miner.AbstractMiner = *(newAbstractMiner(miner, miner))
    return miner
}

func (miner *PropertiesBasedMiner) GetStone(group string, dataId string) (string, error) {
    key := Condition("" == group, dataId, group+"."+dataId).(string)
    property := miner.properties.GetProperty(key)
    return Substitute(property, true, group, dataId, nil)
}

func (miner *PropertiesBasedMiner) DefaultGroupName() string {
    return ""
}
