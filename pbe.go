package diamond_client

import (
    "crypto/cipher"
    "crypto/des"
    "crypto/md5"
    "encoding/base64"
    "strings"
)

type Pbe struct {
    salt           []byte
    iterationCount int
}

func NewPbe(salt []byte, iterationCount int) *Pbe {
    return &Pbe{salt: salt, iterationCount: iterationCount}
}

func NewPbeDefault() *Pbe {
    return NewPbe([]byte{0xde, 0x33, 0x10, 0x12, 0xde, 0x33, 0x10, 0x12}, 20)
}

func (pbe *Pbe) Encrypt(plainText string, password string) (string, error) {
    padNum := byte(8 - len(plainText)%8)
    for i := byte(0); i < padNum; i++ {
        plainText += string(padNum)
    }

    dk, iv := getDerivedKey(password, pbe.salt, pbe.iterationCount)

    block, err := des.NewCipher(dk)

    if err != nil {
        return "", err
    }

    encryptor := cipher.NewCBCEncrypter(block, iv)
    encrypted := make([]byte, len(plainText))
    encryptor.CryptBlocks(encrypted, []byte(plainText))

    return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (pbe *Pbe) Decrypt(cipherText string, password string) (string, error) {
    msgBytes, err := base64.StdEncoding.DecodeString(cipherText)
    if err != nil {
        return "", err
    }

    dk, iv := getDerivedKey(password, pbe.salt, pbe.iterationCount)
    block, err := des.NewCipher(dk)

    if err != nil {
        return "", err
    }

    decryptor := cipher.NewCBCDecrypter(block, iv)
    decrypted := make([]byte, len(msgBytes))
    decryptor.CryptBlocks(decrypted, msgBytes)

    decryptedString := strings.TrimRight(string(decrypted), "\x01\x02\x03\x04\x05\x06\x07\x08")

    return decryptedString, nil
}

func getDerivedKey(password string, salt []byte, count int) ([]byte, []byte) {
    key := md5.Sum([]byte(password + string(salt)))
    for i := 0; i < count-1; i++ {
        key = md5.Sum(key[:])
    }
    return key[:8], key[8:]
}
