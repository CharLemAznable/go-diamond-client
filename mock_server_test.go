package diamond_client

import (
    "testing"
)

func TestMockServer(t *testing.T) {
    if Mock.IsTestMode() {
        t.Fail()
    }
    Mock.SetUpMockServer()
    if !Mock.IsTestMode() {
        t.Fail()
    }

    Mock.SetConfigInfo("group", "test", "123")
    test, _ := Mock.GetDiamond(NewAxis("group", "test"))
    if "123" != test {
        t.Fail()
    }

    infos := map[string]string{
        "a": "aa",
        "b": "bb",
    }
    Mock.SetConfigInfos(infos)
    aa, _ := Mock.GetDiamond(NewAxisOfDefaultGroup("a"))
    if "aa" != aa {
        t.Fail()
    }
    bb, _ := Mock.GetDiamond(NewAxisOfDefaultGroup("b"))
    if "bb" != bb {
        t.Fail()
    }

    Mock.TearDownMockServer()
    if Mock.IsTestMode() {
        t.Fail()
    }
    test, _ = Mock.GetDiamond(NewAxis("group", "test"))
    if "" != test {
        t.Fail()
    }
    aa, _ = Mock.GetDiamond(NewAxisOfDefaultGroup("a"))
    if "" != aa {
        t.Fail()
    }
    bb, _ = Mock.GetDiamond(NewAxisOfDefaultGroup("b"))
    if "" != bb {
        t.Fail()
    }
}

func BenchmarkMockServerParallel(b *testing.B) {
    pbe := NewPbeDefault()
    password := "mypass"
    axis := NewAxisOfDefaultGroup("bench_key")
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            str := RandString(8)
            enc, _ := pbe.Encrypt(str, password)
            Mock.SetDiamond(axis, str+":"+enc)

            temp, _ := Mock.GetDiamond(axis)
            str2 := temp[0:8]
            enc2 := temp[9:]
            dec, _ := pbe.Decrypt(enc2, password)
            if str2 != dec {
                b.Fail()
            }
        }
    })
}
