package diamond_client

import (
    "bytes"
    "errors"
    . "github.com/CharLemAznable/gokits"
    "math/rand"
    "regexp"
    "strings"
    "time"
)

func ToBool(str string) bool {
    lower := strings.ToLower(str)
    return "true" == lower || "yes" == lower || "on" == lower || "y" == lower
}

func CurrentTimeMillis() int64 {
    return time.Now().UnixNano() / 1e6
}

var src = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(length int) string {
    result := make([]byte, 0)
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < length; i++ {
        result = append(result, src[r.Intn(len(src))])
    }
    return string(result)
}

func ParseStoneToProperties(stone string) *Properties {
    properties := NewProperties()
    if "" != stone {
        _ = properties.Load(strings.NewReader(stone))
    }
    return TryDecryptProperties(properties)
}

var encryptPattern = regexp.MustCompile(`{(.+)}`)
var PBE = NewPbeDefault()

func TryDecrypt(original string, dataId string) (string, error) {
    if "" == original {
        return "", nil
    }

    loc := encryptPattern.FindStringIndex(original)
    if nil == loc || 0 != loc[0] {
        return original, nil
    }

    encrypted := original[loc[1]:]
    algorithm := original[loc[0]+1 : loc[1]-1]
    if "PBE" == algorithm {
        return PBE.Decrypt(paddingBase64(encrypted), dataId)
    }

    return "", errors.New(algorithm + " is not supported now")
}

func TryDecryptProperties(properties *Properties) *Properties {
    newProperties := NewProperties()
    for _, key := range properties.StringPropertyNames() {
        property := properties.GetProperty(key)
        decrypt, _ := TryDecrypt(property, key)
        newProperties.Put(key, decrypt)
    }
    return newProperties
}

func paddingBase64(str string) string {
    return padding(str, '=', (4-len(str)%4)%4)
}

func padding(str string, letter byte, repeats int) string {
    buf := bytes.NewBufferString(str)
    for ; repeats > 0; repeats-- {
        buf.WriteByte(letter)
    }
    return buf.String()
}
