package sha256

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

/*
func ComputeHmac256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
*/

func Sha256() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	key := []byte("secret")
	hasher := hmac.New(sha256.New, key)
	//h.Write([]byte(message))

	//hash := md5.New()

	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s MD5 checksum is %x \n", f.Name(), hasher.Sum(nil))
	fmt.Printf(base64.StdEncoding.EncodeToString(hasher.Sum(nil)))

	//fmt.Println(ComputeHmac256("Message", "secret"))
}
