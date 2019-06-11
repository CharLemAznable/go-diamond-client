package diamond_client

import (
    "sync"
    "sync/atomic"
)

type Meta struct {
    axis               *Axis
    mutex              sync.RWMutex
    lastModifiedHeader string
    md5                string
    localFile          string
    localVersion       int64
    useLocal           bool
    successCounter     int64
}

func NewMeta(axis *Axis) *Meta {
    return &Meta{axis: axis}
}

func (meta *Meta) GetAxis() *Axis {
    return meta.axis
}

func (meta *Meta) GetLastModifiedHeader() string {
    meta.mutex.RLock()
    defer meta.mutex.RUnlock()
    return meta.lastModifiedHeader
}

func (meta *Meta) SetLastModifiedHeader(lastModifiedHeader string) {
    meta.mutex.Lock()
    defer meta.mutex.Unlock()
    meta.lastModifiedHeader = lastModifiedHeader
}

func (meta *Meta) GetMd5() string {
    meta.mutex.RLock()
    defer meta.mutex.RUnlock()
    return meta.md5
}

func (meta *Meta) SetMd5(md5 string) {
    meta.mutex.Lock()
    defer meta.mutex.Unlock()
    meta.md5 = md5
}

func (meta *Meta) GetLocalFile() string {
    meta.mutex.RLock()
    defer meta.mutex.RUnlock()
    return meta.localFile
}

func (meta *Meta) SetLocalFile(localFile string) {
    meta.mutex.Lock()
    defer meta.mutex.Unlock()
    meta.localFile = localFile
}

func (meta *Meta) GetLocalVersion() int64 {
    meta.mutex.RLock()
    defer meta.mutex.RUnlock()
    return meta.localVersion
}

func (meta *Meta) SetLocalVersion(localVersion int64) {
    meta.mutex.Lock()
    defer meta.mutex.Unlock()
    meta.localVersion = localVersion
}

func (meta *Meta) IsUseLocal() bool {
    meta.mutex.RLock()
    defer meta.mutex.RUnlock()
    return meta.useLocal
}

func (meta *Meta) SetUseLocal(useLocal bool) {
    meta.mutex.Lock()
    defer meta.mutex.Unlock()
    meta.useLocal = useLocal
}

func (meta *Meta) FetchCount() int64 {
    return atomic.LoadInt64(&(meta.successCounter))
}

func (meta *Meta) IncSuccessCounterAndGet() int64 {
    return atomic.AddInt64(&(meta.successCounter), 1)
}

func (meta *Meta) Clear() {
    meta.SetLastModifiedHeader("")
    meta.SetMd5("")
    meta.SetLocalFile("")
    meta.SetLocalVersion(0)
    meta.SetUseLocal(false)
}
