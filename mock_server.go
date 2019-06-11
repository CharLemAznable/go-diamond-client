package diamond_client

import (
    "github.com/pkg/errors"
    "sync"
)

type mockServer struct {
    mutex    sync.RWMutex
    testMode bool
    mocks    map[Axis]string
}

var Mock = &mockServer{mocks: make(map[Axis]string)}

func (mock *mockServer) SetUpMockServer() {
    mock.mutex.Lock()
    defer mock.mutex.Unlock()
    mock.testMode = true
}

func (mock *mockServer) TearDownMockServer() {
    mock.mutex.Lock()
    defer mock.mutex.Unlock()
    mock.mocks = make(map[Axis]string)
    mock.testMode = false
}

func (mock *mockServer) IsTestMode() bool {
    mock.mutex.RLock()
    defer mock.mutex.RUnlock()
    return mock.testMode
}

func (mock *mockServer) GetDiamond(axis *Axis) (string, error) {
    mock.mutex.RLock()
    defer mock.mutex.RUnlock()
    if diamond, ok := mock.mocks[*axis]; ok {
        return diamond, nil
    }
    return "", errors.New("Mock-diamond un-exists")
}

func (mock *mockServer) SetDiamond(axis *Axis, info string) {
    mock.mutex.Lock()
    defer mock.mutex.Unlock()
    mock.mocks[*axis] = info
}

func (mock *mockServer) SetConfigInfo(group string, dataId string, info string) {
    mock.SetDiamond(NewAxis(group, dataId), info)
}

func (mock *mockServer) SetConfigInfoOfDefaultGroup(dataId string, info string) {
    mock.SetDiamond(NewAxisOfDefaultGroup(dataId), info)
}

func (mock *mockServer) SetConfigInfos(configInfos map[string]string) {
    if nil == configInfos {
        return
    }

    for key, value := range configInfos {
        mock.SetConfigInfoOfDefaultGroup(key, value)
    }
}
