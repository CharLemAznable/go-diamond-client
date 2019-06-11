package diamond_client

import (
    . "github.com/CharLemAznable/gokits"
    "sync"
)

type subscriberServer struct {
    metaCache *CacheTable

    mutex       sync.RWMutex
    managerConf *ManagerConf
    local       *Local
    snapshot    *Snapshot
}

var Subscriber *subscriberServer

func init() {
    Subscriber = &subscriberServer{}
    Subscriber.metaCache = CacheExpireAfterWrite("SubscriberMetaCache")
    Subscriber.metaCache.SetDataLoader(func(key interface{}, args ...interface{}) (*CacheItem, error) {
        axis := key.(*Axis)
        return NewCacheItem(axis, 0, NewMeta(axis)), nil
    })
    Subscriber.managerConf = NewManagerConf()
    Subscriber.local = NewLocal()
    Subscriber.local.Start(Subscriber.managerConf)
    Subscriber.snapshot = NewSnapshot(Subscriber.managerConf)
}

func (subscriber *subscriberServer) GetDiamond(axis *Axis, timeout int64) (string, error) {
    if Mock.IsTestMode() {
        return Mock.GetDiamond(axis)
    }

    result, err := subscriber.RetrieveDiamondLocalAndRemote(axis, timeout)
    if nil == err && "" != result {
        return result, nil
    }

    if Mock.IsTestMode() {
        return "", nil
    }

    return subscriber.GetSnapshot(axis)
}

func (subscriber *subscriberServer) RetrieveDiamondLocalAndRemote(axis *Axis, timeout int64) (string, error) {
    meta, err := subscriber.GetCachedMeta(axis)
    if nil != err {
        return "", err
    }

    localConfig, err := subscriber.local.ReadLocal(meta)
    if nil == err {
        meta.IncSuccessCounterAndGet()
        subscriber.saveSnapshot(axis, localConfig)
        return localConfig, nil
    }

    // TODO retrieveRemote

    return "", nil
}

func (subscriber *subscriberServer) GetSnapshot(axis *Axis) (string, error) {
    meta, err := subscriber.GetCachedMeta(axis)
    if nil != err {
        return "", err
    }
    content, err := subscriber.snapshot.GetSnapshot(axis)
    if nil != err {
        _ = LOG.Error("GetSnapshot diamondAxis %s error %s",
            axis.ToString(), err.Error())
    } else {
        meta.IncSuccessCounterAndGet()
    }
    return content, err
}

func (subscriber *subscriberServer) GetCachedMeta(axis *Axis) (*Meta, error) {
    item, err := subscriber.metaCache.Value(axis)
    if nil != err {
        return nil, err
    }
    return item.Data().(*Meta), nil
}

func (subscriber *subscriberServer) saveSnapshot(axis *Axis, content string) {
    _ = subscriber.snapshot.SaveSnapshot(axis, content)
}
