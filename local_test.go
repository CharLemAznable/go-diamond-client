package diamond_client

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func TestLocalMonitor(t *testing.T) {
    // create snapshot
    userHomeDir, err := os.UserHomeDir()
    if nil != err {
        panic(err)
    }
    snapshotPath := userHomeDir +
        Separator + ".diamond-client" +
        Separator + SnapshotDir +
        Separator + "Local" +
        Separator + "test" + DiamondStoneExt
    if err := os.MkdirAll(filepath.Dir(snapshotPath), os.ModePerm); nil != err {
        panic(err)
    }
    if err := ioutil.WriteFile(snapshotPath, []byte("abc"), 0644); nil != err {
        panic(err)
    }

    // get snapshot
    time.Sleep(time.Millisecond * 500)
    snapStone, _ := Diamond.GetStone("Local", "test")
    if "abc" != snapStone {
        t.Fail()
    }

    // create config-data
    time.Sleep(time.Millisecond * 500)
    localPath := userHomeDir +
        Separator + ".diamond-client" +
        Separator + DataDir +
        Separator + "Local" +
        Separator + "test" + DiamondStoneExt
    if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); nil != err {
        panic(err)
    }
    if err := ioutil.WriteFile(localPath, []byte("def"), 0644); nil != err {
        panic(err)
    }

    // get config-data
    time.Sleep(time.Millisecond * 500)
    localStone, _ := Diamond.GetStone("Local", "test")
    if "def" != localStone {
        t.Fail()
    }

    // remove config-data
    time.Sleep(time.Millisecond * 500)
    if err := os.RemoveAll(filepath.Dir(localPath)); nil != err {
        panic(err)
    }

    // get snapshot which has been re-written by config-data
    time.Sleep(time.Millisecond * 500)
    fallbackStone, _ := Diamond.GetStone("Local", "test")
    if "def" != fallbackStone {
        t.Fail()
    }

    // remove snapshot
    time.Sleep(time.Millisecond * 500)
    if err := os.RemoveAll(filepath.Dir(snapshotPath)); nil != err {
        panic(err)
    }

    // get empty string
    time.Sleep(time.Millisecond * 500)
    emptyStone, _ := Diamond.GetStone("Local", "test")
    if "" != emptyStone {
        t.Fail()
    }
}
