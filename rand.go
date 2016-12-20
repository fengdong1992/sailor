package sailor

import (
        "time"
        "fmt"
        "io"
        "crypto/md5"
)

func RandString() (randString string) {
    randString = "null"
    t := time.Now()
    h := md5.New()
    io.WriteString(h, "www.airdb.com")
    io.WriteString(h, t.String())
    randString = fmt.Sprintf("%x", h.Sum(nil))
    return randString
}
