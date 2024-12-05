package kraken

import (
	"fmt"
	"golang.org/x/exp/rand"
	"net"
	"path"
	"time"
)

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func getActivePort() (port int, err error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		err = fmt.Errorf("rrror starting listener: %w", err)
		return
	}
	defer func() {
		_ = listener.Close()
	}()
	addr := listener.Addr().(*net.TCPAddr)
	port = addr.Port
	return
}

func randInt(min, max int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}

func JustSleep() {
	t := time.Duration(randInt(1, 3)) * time.Millisecond
	log.Debugf("just sleep: %s", t)
	time.Sleep(t)
}

func SleepForSlowOperation() {
	t := time.Duration(randInt(5, 15)) * time.Millisecond
	log.Debugf("sleep: %s for slow operation", t)
	time.Sleep(t)
}

func SleepRandSeconds(min, max int) {
	t := time.Duration(randInt(min, max)) * time.Millisecond
	log.Debugf("sleep rand time: %s", t)
	time.Sleep(t)
}

func prepareEventScript(event string) string {
	return fmt.Sprintf("const event = new MouseEvent('%s', { bubbles: true }); arguments[0].dispatchEvent(event);", event)
}

func sumTimeout(timeout []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range timeout {
		total += d
	}
	return total
}
