package main

import (
    "os"
    "fmt"
    "golang.org/x/crypto/sha3"
    "time"
    "encoding/hex"
    "math/big"
    "hash"
    "encoding/binary"
)

func FormatHashrate(b int64) string {
    const unit = 1000
    if b < unit {
            return fmt.Sprintf("%d H/s", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
            div *= unit
            exp++
    }
    return fmt.Sprintf("%.1f %cH/S", float64(b)/float64(div), "kMGTPE"[exp])
}


type KeccakState interface {
    hash.Hash
    Read([]byte) (int, error)
}

func hashData(data ...[]byte) []byte {
    b := make([]byte, 32)
    d := sha3.NewLegacyKeccak256().(KeccakState)
    for _, b := range data {
        d.Write(b)
    }
    d.Read(b)
    return b
}


func main() {
    fmt.Println("Welcome to the SiriCoin GO Miner!")

	starttime := time.Now()
    lasti := 0
    hashrate := int64(0)
    target := os.Args[1]

    root, err := hex.DecodeString(os.Args[2])

  if err != nil { panic(err) }

    e := new(big.Int)
    e.SetString(target, 0)

    i := 1

    for {
        i++
            if time.Now().Sub(starttime) > time.Second * 3 {
                hashrate = int64((i - lasti) / 3)
                fmt.Println("Hashrate:", FormatHashrate(hashrate))
				starttime = time.Now()
                lasti = i
			}

            data := make([]byte, 32)
            binary.BigEndian.PutUint64(data[len(data)-8:], uint64(i))

            hashraw := append(root, data...)


            h := hashData(hashraw)

            if new(big.Int).SetBytes(h[:]).Cmp(e) == -1 {
                fmt.Println("you found a share")
                fmt.Println("Hash:", hex.EncodeToString(h), "Nonce:", i)
                break
        }
    }
}
