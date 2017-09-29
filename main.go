package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verify [PATH]")
		os.Exit(2)
	}

	algos := []algo{
		algo{name: "MD5       ", hasher: md5.New()},
		algo{name: "SHA1      ", hasher: sha1.New()},
		algo{name: "SHA256    ", hasher: sha256.New()},
		algo{name: "SHA256/224", hasher: sha256.New224()},
		algo{name: "SHA512    ", hasher: sha512.New()},
		algo{name: "SHA512/224", hasher: sha512.New512_224()},
		algo{name: "SHA512/256", hasher: sha512.New512_256()},
		algo{name: "SHA512/384", hasher: sha512.New384()},
	}

	if err := hashFile(os.Args[1], algos...); err != nil {
		log.Fatal(err)
	}
}

type algo struct {
	name   string
	hasher hash.Hash
}

func hashFile(fullPath string, algos ...algo) (err error) {
	file, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer file.Close()

	for _, a := range algos {
		if _, err = io.Copy(a.hasher, file); err != nil {
			return
		}
		digest := a.hasher.Sum(nil)
		first, last := digest[0:3], digest[len(digest)-3:]
		fmt.Printf("%s (%s) = %x...%x\n", a.name, filepath.Clean(fullPath), first, last)

		if _, err = file.Seek(0, 0); err != nil {
			return
		}
	}
	return
}

/*
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x", h.Sum(nil))
*/
