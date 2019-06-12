package diamond_client

import (
    . "github.com/CharLemAznable/gokits"
    "io/ioutil"
    "os"
)

type Snapshot struct {
    dir string
}

func NewSnapshot(managerConf *ManagerConf) *Snapshot {
    dir := managerConf.GetFilePath() + Separator + SnapshotDir
    err := os.MkdirAll(dir, os.ModePerm)
    if nil != err {
        panic("create Snapshot dir fail " + dir)
    }

    return &Snapshot{dir: dir}
}

func (snapshot *Snapshot) GetSnapshot(axis *Axis) (string, error) {
    return snapshot.getFileContent(axis, DiamondStoneExt)
}

func (snapshot *Snapshot) getFileContent(axis *Axis, extension string) (string, error) {
    path := snapshot.dir + Separator +
        axis.GetGroup() + Separator +
        axis.GetDataId() + extension
    file, err := os.Open(path)
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

func (snapshot *Snapshot) SaveSnapshot(axis *Axis, content string) error {
    err := snapshot.saveFileContent(axis, content, DiamondStoneExt)
    if nil != err {
        _ = LOG.Error("save snapshot error %s by %s : %s", axis, content, err.Error())
    }
    return err
}

func (snapshot *Snapshot) saveFileContent(axis *Axis, content string, extension string) error {
    path := snapshot.dir + Separator + axis.GetGroup()
    err := os.MkdirAll(path, os.ModePerm)
    if nil != err {
        return err
    }

    filepath := path + Separator + axis.GetDataId() + extension
    return ioutil.WriteFile(filepath, []byte(content), 0644)
}
