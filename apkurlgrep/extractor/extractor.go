package extractor

import (
	"example.com/cc/util"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var founds []string
var paths []string

const regexUrlsString = `(?:"|'|href=|src=)((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{2,})|([a-zA-Z0-9_\-/.]{1,}/[a-zA-Z0-9_\-/.]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/.]{1,}/[a-zA-Z0-9_\-/.]{0,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}?\.[a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|))|((?:"|')*([A-Za-z]{10,30})(?:\?[A-Za-z=]{3,})(?:"|')*)|(?:[\?|#][^"|']{0,}|)[a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)(?:"|')`
const regexUrlPath = `\".*/.*\"`

var blacklist_urls_path = []string{"layout", "activity", "text/plain", "application/json", "/proc/", "/vendor/", "/data/local", "content://", "kotlin/", "multipart/form-data", "M/d/yy", "fonts/", ".tft", ".js", ".jpg", ".jpeg", ".png", "application/octet-stream", "kotlinx/", " ", "image/jpeg", "message/rfc822", "multipart/", "native-libs/", "video/x-matroska", "image/heic", "cmdline",
	"/system/app/Superuser.apk", "/system/xbin/su", "application/x-www-form-urlencoded", "image/heif", "model/gltf-binary", ".txt", ".js", "text/",
	"google/", "application/javascript", "image/", "market://", "audio/", "text/", "/system/bin", "META-INF/MANIFEST.MF",
	"java/util", "java/lang", "/collections/", "/Nothing", "/Throwable", "/Number", "/Byte", "/Double", "/Float", "/Long", "/Short", "/Boolean", "/Char", "/CharSequence", "/String", "/Comparable", "/Enum", "/Array", "/ByteArray", "/DoubleArray", "/FloatArray", "/IntArray", "/LongArray", "/ShortArray", "/BooleanArray", "/CharArray", "/Cloneable",
	".properties", "application/xml", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", "okio/Okio",
	"/jvm/", "video/", "react-dom/server", "/sbin/su", "/system/xbin", "/su/bin", "http/",
	"application/x-sentry-envelope", "__", "application/x-protobuffer", "dexopt/baseline.prof", "react-navigation-shared-element", "sentry.io", "u000", "color/", "style/", "style", "dimen/", "animator/", "drawable/", "id/",
	"META-INF/", "font/", "okhttp3/", "string/", "button", "�", "integer", "image", "const", "float", "double", "long", "int", "u00", "androidx", "Model", "Impl", "HTTP", "this.", "Activity", "okhttp/", "org/", "//", "dd", "MM", "invoke-super/range", "yyyy", "move/", "move-wide/", "invoke-direct", "move-object",
	"invoke-virtual", "anim/", "array", "assets", "xml/", "sha", "^", "%", "$", "!", "[", "]", "@", "~", "bool", "font", "manifest", "dimen", "index.html", "string", "xbin", "sdk", "output", "012345678", "attr", "mipmap", "application", "<", ">", "dev/", "mnt/", "data/data/", "system/lib", "sys/"}

var regexpUrls = regexp.MustCompile(regexUrlsString)
var regexpPath = regexp.MustCompile(regexUrlPath)

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func extractTextFromFile(path string) error {

	var textBytes, er = ioutil.ReadFile(path)
	if er != nil {
		panic(er)
	}

	var indexes = regexpUrls.FindAllIndex(textBytes, -1)
	var path_index = regexpPath.FindAllIndex(textBytes, -1)
	if len(indexes) != 0 {
		for _, k := range indexes {
			var foundStart = k[0]
			var foundEnd = k[1]
			var link = string(textBytes[foundStart:foundEnd])
			founds = append(founds, link)
		}
	}
	if len(path_index) != 0 {
		for _, k := range path_index {
			var foundStart = k[0]
			var foundEnd = k[1]
			var link = string(textBytes[foundStart:foundEnd])
			paths = append(paths, link)
		}
	}
	return nil
}

func doHashWalk(dirPath string) error {
	fullPath, err := filepath.Abs(dirPath)

	if err != nil {
		return err
	}

	callback := func(path string, fi os.FileInfo, err error) error {
		return hashFile(path, fi, err)
	}

	return filepath.Walk(fullPath, callback)
}

func hashFile(path string, fileInfo os.FileInfo, err error) error {

	if util.SkipFolder(path) {
		return nil
	}

	var fileName = fileInfo.Name()

	if fileInfo.IsDir() {
		return nil
	}

	if SkipExtension(fileName) {
		return nil
	}

	if err != nil {
		return err
	}

	extractTextFromFile(path)
	return nil
}

func sortPaths(paths []string) []string {

	paths = unique(paths)
	var sortedPaths []string
	for i := 1; i < len(paths); i++ {
		if !util.ContainsNonASCII(paths[i]) && !util.HasUpper(paths[i]) && !util.Contains(blacklist_urls_path, paths[i]) {
			sortedPaths = append(sortedPaths, paths[i])
		}
	}

	return sortedPaths

}

func sortUrls(urls []string) []string {
	blacklist_urls := []string{"google.com", "crashlytics.com", "googlesyndication.com", "android.com", "facebook.com", "app-measurement.com", "googleadservices.com", "googleapis.com", "twitter.com",
		"w3.org", "mozilla.org", "adobe.com", "pinterest.com", "linkedin.com", "paypal.com", "yahoo.com", "recaptcha.net", "live.com", "vimeo.com", "apache.org"}

	urls = unique(urls)

	var sortedUrls []string

	for i := 1; i < len(urls); i++ {

		urls[i] = strings.ReplaceAll(urls[i], "'", "")
		urls[i] = strings.ReplaceAll(urls[i], "\"", "")

		if len(urls[i]) < 5 {
			continue
		}

		if urls[i][:4] == "http" || urls[i][:5] == "https" {
			process_url, err := url.Parse(urls[i])

			if err == nil && !util.Contains(blacklist_urls, process_url.Host) {
				sortedUrls = append(sortedUrls, urls[i])
				continue
			}
		}

	}

	return sortedUrls
}

func Extract(tempDir string) {

	doHashWalk(tempDir)

	sortedUrls := sortUrls(founds)
	sortedPaths := sortPaths(paths)

	if len(sortedUrls) > 0 {
		fmt.Println("Result of URLs:")
		fmt.Printf("\n \n")
		fmt.Printf(strings.Join(sortedUrls, "\n"))
	}

	fmt.Printf("\n \n \n")

	if len(sortedPaths) > 0 {
		fmt.Println("Result of URLs Paths:")
		fmt.Printf("\n \n")
		fmt.Printf(strings.Join(sortedPaths, "\n"))
	}

}
