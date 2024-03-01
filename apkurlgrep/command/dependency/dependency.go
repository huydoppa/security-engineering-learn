package command

import (
	"os/exec"
)

func AreAllReady() bool {
	var areAllReady = true

	areAllReady = isApktoolInstalled()
	areAllReady = isHermescInstalled()
	if areAllReady != true {
		return false
	}

	return true
}

func isApktoolInstalled() bool {
	_, err := exec.LookPath("apktool")
	if err != nil {
		panic("Didn't find 'apktool' executable.")
		return false
	}

	return true
}

func isHermescInstalled() bool {
	_, err := exec.LookPath("hbc-decompiler")
	if err != nil {
		panic("Didn't find 'hbc-decompiler' executable.")
		return false
	}

	return true
}
