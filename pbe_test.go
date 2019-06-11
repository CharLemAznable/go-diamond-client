package diamond_client

import (
    "testing"
)

func TestPbe(t *testing.T) {
    pbe := NewPbeDefault()
    originalPassword := "secret"

    encryptedPassword, err := pbe.Encrypt(originalPassword, "mypass")
    if nil != err || "mzb7VnJcM/c=" != encryptedPassword {
        t.Fail()
    }
    decryptedPassword, err := pbe.Decrypt(encryptedPassword, "mypass")
    if nil != err || originalPassword != decryptedPassword {
        t.Fail()
    }
}

func BenchmarkPbeParallel(b *testing.B) {
    pbe := NewPbeDefault()
    password := "mypass"
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            str := RandString(8)
            enc, _ := pbe.Encrypt(str, password)
            dec, _ := pbe.Decrypt(enc, password)
            if str != dec {
                b.Fail()
            }
        }
    })
}
