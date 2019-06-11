package diamond_client

import (
    "errors"
    . "github.com/CharLemAznable/gokits"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
)

type Local struct {
    existFilesTimestamp map[string]int64
    mutex               sync.RWMutex
    running             bool
    rootPath            string
}

func NewLocal() *Local {
    return &Local{existFilesTimestamp: make(map[string]int64)}
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
    // TODO
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
