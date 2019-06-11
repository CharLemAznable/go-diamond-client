package diamond_client

type Manager struct {
    axis          *Axis
    timeoutMillis int64
}

func NewManager(group string, dataId string) *Manager {
    axis := NewAxis(group, dataId)
    _, _ = Subscriber.GetCachedMeta(axis)
    return &Manager{axis: axis, timeoutMillis: 1000}
}

func NewManagerOfDefaultGroup(dataId string) *Manager {
    return NewManager(DefaultGroupName, dataId)
}

// func (manager *Manager) AddDiamondListener(diamondListener *DiamondListener) {
//     diamondSubscriber.AddDiamondListener(axis, diamondListener)
// }

// func (manager *Manager) RemoveDiamondListener(diamondListener *DiamondListener) {
//     diamondSubscriber.RemoveDiamondListener(axis, diamondListener)
// }

func (manager *Manager) GetDiamond() (string, error) {
    original, err := Subscriber.GetDiamond(manager.axis, manager.timeoutMillis)
    if nil != err {
        return "", err
    }
    return TryDecrypt(original, manager.axis.GetDataId())
}

// public Object getCache() {
// return diamondSubscriber.getCache(axis, timeoutMillis);
// }

// public Object getDynamicCache(Object... dynamics) {
// return diamondSubscriber.getCache(axis, timeoutMillis, dynamics);
// }

func (manager *Manager) SetTimeoutMillis(timeoutMillis int64) {
    manager.timeoutMillis = timeoutMillis
}
