package diamond_client

import (
    "errors"
    . "github.com/CharLemAznable/gokits"
    "github.com/radovskyb/watcher"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "time"
)

type Local struct {
    existFilesTimestamp map[string]int64
    mutex               sync.RWMutex
    running             bool
    rootPath            string
    monitor             *watcher.Watcher
}

func NewLocal() *Local {
    return &Local{existFilesTimestamp: make(map[string]int64), monitor: watcher.New()}
}

func (local *Local) Start(conf *ManagerConf) {
    local.mutex.Lock()
    defer local.mutex.Unlock()

    if local.running {
        return
    }
    local.running = true
    local.rootPath = conf.GetFilePath() + Separator + DataDir

    local.initDataDir()
    local.startCheckLocalDir()
}

func (local *Local) initDataDir() {
    err := os.MkdirAll(local.rootPath, os.ModePerm)
    if nil != err {
        panic("create Local dir fail " + local.rootPath)
    }
}

func (local *Local) startCheckLocalDir() {
    local.index()
    local.watchRoot()
}

func (local *Local) index() {
    rootAbs, err := filepath.Abs(local.rootPath)
    if nil != err {
        _ = LOG.Error("local-diamond index error: ", err)
        return
    }
    _ = filepath.Walk(local.rootPath,
        func(path string, f os.FileInfo, err error) error {
            if nil == f {
                return err
            }
            if f.IsDir() {
                return nil
            }
            abs, err := filepath.Abs(path)
            if nil != err || DiamondStoneExt != filepath.Ext(abs) {
                return nil
            }
            gpd := filepath.Dir(filepath.Dir(abs))
            if rootAbs != gpd {
                return nil
            }
            local.existFilesTimestamp[abs] = CurrentTimeMillis()
            LOG.Debug("%s file was added", abs)

            return nil
        })
}

func (local *Local) watchRoot() {
    local.monitor.FilterOps(
        watcher.Create, watcher.Write, watcher.Remove,
        watcher.Rename, watcher.Chmod, watcher.Move)

    rootAbs, err := filepath.Abs(local.rootPath)
    if nil != err {
        _ = LOG.Error("local-diamond watch root error: ", err)
        return
    }
    valid := regexp.MustCompile(rootAbs + ".+\\." + DiamondStoneExt[1:])
    local.monitor.AddFilterHook(watcher.RegexFilterHook(valid, true))
    go func() {
        for {
            select {
            case event := <-local.monitor.Event:
                switch event.Op {
                case watcher.Create, watcher.Write, watcher.Chmod:
                    local.existFilesTimestamp[event.Path] = CurrentTimeMillis()
                case watcher.Remove:
                    delete(local.existFilesTimestamp, event.Path)
                case watcher.Rename, watcher.Move:
                    paths := strings.Split(event.Path, " -> ")
                    delete(local.existFilesTimestamp, paths[0])
                    local.existFilesTimestamp[paths[1]] = CurrentTimeMillis()
                default:
                    LOG.Debug("unexpected Event: %s", event)
                }
            case err := <-local.monitor.Error:
                _ = LOG.Error("local-diamond monitor error: ", err)
            case <-local.monitor.Closed:
                LOG.Info("local-diamond monitor closed")
                return
            }
        }
    }()
    if err := local.monitor.AddRecursive(local.rootPath); err != nil {
        _ = LOG.Error("local-diamond monitor AddRecursive error: ", err)
        return
    }
    go func() {
        if err := local.monitor.Start(time.Millisecond * 100); err != nil {
            _ = LOG.Error("local-diamond monitor Start error: ", err)
        }
    }()
    local.monitor.Wait()
}

func (local *Local) ReadLocal(meta *Meta) (string, error) {
    localFilePath, err := local.getFilePath(meta.GetAxis())
    if nil != err {
        return "", err
    }

    _, ok := local.existFilesTimestamp[localFilePath]
    if !ok {
        if meta.IsUseLocal() {
            meta.Clear()
        }
        return "", errors.New("local-diamond un-exists")
    }

    return local.readFileContent(localFilePath)
}

func (local *Local) getFilePath(axis *Axis) (string, error) {
    path := local.rootPath + Separator +
        axis.GetGroup() + Separator +
        axis.GetDataId() + DiamondStoneExt
    return filepath.Abs(path)
}

func (local *Local) readFileContent(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if nil != err {
        return "", err
    }
    defer func() { _ = file.Close() }()

    bytes, err := ioutil.ReadAll(file)
    if nil != err {
        return "", err
    }
    return string(bytes), nil
}
