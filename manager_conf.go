package diamond_client

import (
    . "github.com/CharLemAznable/gokits"
    "math"
    "math/rand"
    "os"
    "sync"
    "sync/atomic"
)

type ManagerConf struct {
    mutex                          sync.RWMutex
    pollingInterval                int
    onceTimeout                    int
    receiveWaitTime                int
    domainNamePos                  uint32
    servers                        []string
    maxHostConnections             int
    connectionStaleCheckingEnabled bool
    maxTotalConnections            int
    connectionTimeout              int
    retrieveDataRetryTimes         int
    filePath                       string // local data dir root
}

func NewManagerConf() *ManagerConf {
    userHomeDir, err := os.UserHomeDir()
    if nil != err {
        panic(err)
    }
    filePath := userHomeDir + Separator + ".diamond-client"
    err = os.MkdirAll(filePath, os.ModePerm)
    if nil != err {
        panic("create diamond-miner dir fail " + filePath)
    }

    return &ManagerConf{
        pollingInterval:                PollingInterval,
        onceTimeout:                    OnceTimeout,
        receiveWaitTime:                RecvWaitTimeout,
        servers:                        make([]string, 0),
        maxHostConnections:             1,
        connectionStaleCheckingEnabled: true,
        maxTotalConnections:            20,
        connectionTimeout:              ConnTimeout,
        retrieveDataRetryTimes:         math.MaxInt32 / 10,
        filePath:                       filePath,
    }
}

func (conf *ManagerConf) GetPollingInterval() int {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.pollingInterval
}

func (conf *ManagerConf) SetPollingInterval(pollingInterval int) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    if conf.pollingInterval < PollingInterval && !Mock.IsTestMode() {
        return
    }
    conf.pollingInterval = pollingInterval
}

func (conf *ManagerConf) GetOnceTimeout() int {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.onceTimeout
}

func (conf *ManagerConf) SetOnceTimeout(onceTimeout int) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.onceTimeout = onceTimeout
}

func (conf *ManagerConf) GetReceiveWaitTime() int {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.receiveWaitTime
}

func (conf *ManagerConf) SetReceiveWaitTime(receiveWaitTime int) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.receiveWaitTime = receiveWaitTime
}

func (conf *ManagerConf) GetDomainName() string {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    if 0 == len(conf.servers) {
        panic("no name server available!")
    }
    return conf.servers[atomic.LoadUint32(&(conf.domainNamePos))]
}

func (conf *ManagerConf) randomDomainNamePos() {
    serverNum := int32(len(conf.servers))
    if serverNum > 1 {
        atomic.StoreUint32(&(conf.domainNamePos), uint32(rand.Int31n(serverNum)))
        LOG.Info("random DiamondServer to：" + conf.GetDomainName())
    }
}

func (conf *ManagerConf) RotateToNextDomain(httpClient *HttpClient) {
    serverNum := (int32)(len(conf.servers))
    if serverNum == 0 {
        _ = LOG.Error("diamond server list is empty, please contact administrator")
        return
    }

    if serverNum <= 1 {
        httpClient.ResetHostConfig(conf.GetDomainName())
        return
    }

    index := atomic.AddUint32(&(conf.domainNamePos), 1)
    atomic.StoreUint32(&(conf.domainNamePos), index%uint32(serverNum))
    httpClient.ResetHostConfig(conf.GetDomainName())
    _ = LOG.Warn("rotate diamond server to " + conf.GetDomainName())
}

func (conf *ManagerConf) GetServers() []string {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.servers
}

func (conf *ManagerConf) HasServers() bool {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return len(conf.servers) > 0
}

func (conf *ManagerConf) SetServers(servers *Hashset, httpClient *HttpClient) {
    newServers := make([]string, 0)
    for _, server := range servers.Items() {
        newServers = append(newServers, server.(string))
    }
    conf.mutex.Lock()
    conf.servers = newServers
    conf.mutex.Unlock()

    conf.randomDomainNamePos()
    httpClient.ResetHostConfig(conf.GetDomainName())
}

func (conf *ManagerConf) AddDomainName(domainName string) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.servers = append(conf.servers, domainName)
}

func (conf *ManagerConf) GetMaxHostConnections() int {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.maxHostConnections
}

func (conf *ManagerConf) SetMaxHostConnections(maxHostConnections int) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.maxHostConnections = maxHostConnections
}

func (conf *ManagerConf) IsConnectionStaleCheckingEnabled() bool {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.connectionStaleCheckingEnabled
}

func (conf *ManagerConf) SetConnectionStaleCheckingEnabled(connectionStaleCheckingEnabled bool) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.connectionStaleCheckingEnabled = connectionStaleCheckingEnabled
}

func (conf *ManagerConf) GetMaxTotalConnections() int {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.maxTotalConnections
}

func (conf *ManagerConf) SetMaxTotalConnections(maxTotalConnections int) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.maxTotalConnections = maxTotalConnections
}

func (conf *ManagerConf) GetConnectionTimeout() int {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.connectionTimeout
}

func (conf *ManagerConf) SetConnectionTimeout(connectionTimeout int) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.connectionTimeout = connectionTimeout
}

func (conf *ManagerConf) GetRetrieveDataRetryTimes() int {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.retrieveDataRetryTimes
}

func (conf *ManagerConf) SetRetrieveDataRetryTimes(retrieveDataRetryTimes int) {
    conf.mutex.Lock()
    defer conf.mutex.Unlock()
    conf.retrieveDataRetryTimes = retrieveDataRetryTimes
}

func (conf *ManagerConf) GetFilePath() string {
    conf.mutex.RLock()
    defer conf.mutex.RUnlock()
    return conf.filePath
}
