package diamond_client

type Miner struct {
    AbstractMiner
    defaultGroupName string
}

func NewMiner() *Miner {
    return NewMinerDefaultGroupName(DefaultGroupName)
}

func NewMinerDefaultGroupName(defaultGroupName string) *Miner {
    miner := &Miner{defaultGroupName: defaultGroupName}
    miner.AbstractMiner = *(newAbstractMiner(miner, miner))
    return miner
}

func (miner *Miner) GetStone(group string, dataId string) (string, error) {
    diamond, err := NewManager(group, dataId).GetDiamond()
    if nil != err {
        return "", err
    }
    return Substitute(diamond, true, group, dataId, nil)
}

func (miner *Miner) DefaultGroupName() string {
    return miner.defaultGroupName
}
