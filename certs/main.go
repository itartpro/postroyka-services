package main

import (
	"os"
	"os/exec"
	"time"
)

func isOlderThanOneYear(t time.Time) bool {
	return time.Now().Sub(t) > 365*24*time.Hour
}

func main() {
	var valid = true

	files := []string{"ca.cert", "ca.key", "ca.srl", "jwtkey.pem", "jwtkey.pub", "service.csr", "service.key", "service.pem"}
	for _, v := range files {
		f, err := os.Stat(v)
		if os.IsNotExist(err) {
			valid = false
		} else {
			if f.Mode().IsRegular() {
				if isOlderThanOneYear(f.ModTime()) {
					valid = false
				}
			}
		}
	}

	if valid != true {
		cmd := exec.Command("sh", "newcerts.sh")
		err := cmd.Run()
		if err != nil {
			println("Certs microservice error:", err.Error())
		} else {
			println("Made new certs")
		}
	} else {
		println("Certs are in place")
	}
}
