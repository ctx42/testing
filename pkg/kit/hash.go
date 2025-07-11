// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package kit

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// SHA1Reader returns SHA1 hash off everything in the reader. Panics on error.
func SHA1Reader(r io.Reader) string {
	hash := sha1.New()
	if _, err := io.Copy(hash, r); err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA1File returns SHA1 hash of the file. Panics on error.
func SHA1File(pth string) string {
	fil, err := os.Open(pth)
	if err != nil {
		panic(err)
	}
	defer func() { _ = fil.Close() }()
	return SHA1Reader(fil)
}
