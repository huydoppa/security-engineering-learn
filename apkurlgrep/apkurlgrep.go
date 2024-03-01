/*
 made with love by @ndelphit 5/2020
*/

package main

import (
	"example.com/cc/command/apktool"
	"example.com/cc/directory"
	"example.com/cc/extractor"
	"example.com/cc/util"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"os/exec"
	"time"
)

func main() {

	parser := argparse.NewParser("apkurlgrep", "ApkUrlGrep - Extract endpoints from APK files")
	apk := parser.String("a", "apk", &argparse.Options{Required: true, Help: "Input a path to APK file."})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}

	var baseApk = *apk

	var tempDir = directory.CreateTempDir()
	fmt.Println(tempDir)
	apktool.RunApktool(baseApk, tempDir)
	time1 := time.Now()
	fmt.Println("Thời gian bắt đầu là : ", time1)
	if util.IsReactNative(tempDir) {
		fmt.Println(tempDir)
		cmd := exec.Command("hbc-decompiler", tempDir+"\\assets\\index.android.bundle", tempDir+"\\assets\\index_decompile.js")
		output, err := cmd.CombinedOutput()
		if err == nil {
			fmt.Println(output)
		}
		extractor.Extract(tempDir + "\\assets")
	} else {
		extractor.Extract(tempDir)
	}
	directory.RemoveTempDir(tempDir)
	fmt.Println("\n\n Thời gian kết thúc là : ", time.Now())
}
