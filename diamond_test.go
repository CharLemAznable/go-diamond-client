package diamond_client

import (
    "testing"
)

func TestDiamondBasic(t *testing.T) {
    Mock.SetUpMockServer()

    Mock.SetConfigInfoOfDefaultGroup("SOLR_URL", "abc")
    solrUrl1 := Diamond.GetString("SOLR_URL")
    if "abc" != solrUrl1 {
        t.Fail()
    }
    solrUrl2, _ := Diamond.GetStone("DEFAULT_GROUP", "SOLR_URL")
    if "abc" != solrUrl2 {
        t.Fail()
    }

    Mock.TearDownMockServer()
}

func TestDiamondDecrypt(t *testing.T) {
    Mock.SetUpMockServer()

    Mock.SetConfigInfo("EqlConfig", "DEFAULT", "mypass={PBE}mzb7VnJcM/c=")
    properties, _ := Diamond.GetProperties("EqlConfig", "DEFAULT")
    if "secret" != properties.GetProperty("mypass") {
        t.Fail()
    }

    Mock.SetConfigInfo("EqlConfig", "mypass", "{PBE}mzb7VnJcM/c=")
    stone, _ := Diamond.GetStone("EqlConfig", "mypass")
    if "secret" != stone {
        t.Fail()
    }

    Mock.TearDownMockServer()
}
