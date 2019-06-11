package diamond_client

import (
    "os"
    "testing"
)

func TestNewManagerConf(t *testing.T) {
    conf := NewManagerConf()
    info, _ := os.Stat(conf.GetFilePath())
    if !info.IsDir() {
        t.Fail()
    }
}
