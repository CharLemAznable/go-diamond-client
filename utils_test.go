package diamond_client

import (
    "testing"
)

func TestToBool(t *testing.T) {
    if !ToBool("TRUE") {
        t.Fail()
    }
    if !ToBool("yes") {
        t.Fail()
    }
    if !ToBool("On") {
        t.Fail()
    }
    if !ToBool("Y") {
        t.Fail()
    }
    if ToBool("false") {
        t.Fail()
    }
    if ToBool("none") {
        t.Fail()
    }
}

func TestParseStoneToProperties(t *testing.T) {
    str := "#redis配置\n" +
        "redis.hosts=127.0.0.1:6379\n" +
        "redis.host=127.0.0.1\n" +
        "redis.port=6379\n" +
        "\n" +
        "#静态资源配置\n" +
        "static.prefix=http://test.fshow.easy-hi.com/fshow-res\n" +
        "\n" +
        "#用户上传图片访问前缀\n" +
        "attach.imgPrefix=http://test.res.fshow.easy-hi.com/images/" +
        "\n" +
        "#访问前缀\n" +
        "root=http://test.fshow.easy-hi.com/fshow/\n" +
        "static.version=1\n" +
        "\n" +
        "#活动分享地址前缀\n" +
        "#activity.prefix=http://go.weyoung.com/mobile/mobile/app/\n" +
        "#website.prefix=http://test.co.easy-hi.com/\n" +
        "mode=dev\n" +
        "#requirejs.optimze=on\n" +
        "template.path=http://test.fshow.easy-hi.com:8000/fshow-res/dev/modules/\n" +
        "innerResPath.prefix = http://test.fshow.easy-hi.com:8000/fshow-res\n" +
        "music.prefix=http://test.res.fshow.easy-hi.com/musics/\n" +
        "\n" +
        "#boss校验地址\n" +
        "origin.boss=http://127.0.0.1:8017/boss-biz/authorize/check-token\n" +
        "\n" +
        "#外接系统wxjsconfig调用地址\n" +
        "wxconfig=http://test.go.easy-hi.com/admin/scene/show/center/1508666666/initWxJs\n"
    properties := ParseStoneToProperties(str)

    if properties.GetProperty("redis.hosts") != "127.0.0.1:6379" {
        t.Fail()
    }
    if properties.GetProperty("redis.host") != "127.0.0.1" {
        t.Fail()
    }
    if properties.GetProperty("redis.port") != "6379" {
        t.Fail()
    }
    if properties.GetProperty("static.prefix") != "http://test.fshow.easy-hi.com/fshow-res" {
        t.Fail()
    }
    if properties.GetProperty("attach.imgPrefix") != "http://test.res.fshow.easy-hi.com/images/" {
        t.Fail()
    }
    if properties.GetProperty("root") != "http://test.fshow.easy-hi.com/fshow/" {
        t.Fail()
    }
    if properties.GetProperty("static.version") != "1" {
        t.Fail()
    }
    if properties.GetProperty("mode") != "dev" {
        t.Fail()
    }
    if properties.GetProperty("template.path") != "http://test.fshow.easy-hi.com:8000/fshow-res/dev/modules/" {
        t.Fail()
    }
    if properties.GetProperty("innerResPath.prefix") != "http://test.fshow.easy-hi.com:8000/fshow-res" {
        t.Fail()
    }
    if properties.GetProperty("music.prefix") != "http://test.res.fshow.easy-hi.com/musics/" {
        t.Fail()
    }
    if properties.GetProperty("origin.boss") != "http://127.0.0.1:8017/boss-biz/authorize/check-token" {
        t.Fail()
    }
    if properties.GetProperty("wxconfig") != "http://test.go.easy-hi.com/admin/scene/show/center/1508666666/initWxJs" {
        t.Fail()
    }
}

func TestTryDecrypt(t *testing.T) {
    decrypt, err := TryDecrypt("{PBE}mzb7VnJcM/c", "mypass")
    if nil != err || "secret" != decrypt {
        t.Fail()
    }
}
