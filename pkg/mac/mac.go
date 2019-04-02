package mac

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/AiflooAB/prometheus-mac-proxy/pkg/proxy"
)

type macTransformer struct {
	ipmac map[string]string
}

func NewFileTransformer() proxy.Transformer {
	return &macTransformer{
		ipmac: readIPMacMap(),
	}
}

func (trans *macTransformer) Transform(line string) string {
	mac := getMac(trans.ipmac, line)

	return strings.Replace(line, "ip=", fmt.Sprintf("mac=\"%s\",ip=", mac), 1)
}

var ipRegex = regexp.MustCompile(`ip="(.+?)"`)

func getMac(ipmac map[string]string, line string) string {
	matches := ipRegex.FindStringSubmatch(line)
	if len(matches) >= 1 {
		ip := matches[1]
		mac, found := ipmac[ip]
		if found {
			return mac
		}
	}
	return "UNKNOWN"
}

func readIPMacMap() map[string]string {
	ipmac := make(map[string]string)
	file, err := os.Open("./all") // TODO
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ip, mac := parsePromscanLine(line)
		ipmac[ip] = mac
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Got error while reading file: %v\n", err)
	}

	return ipmac
}

func parsePromscanLine(line string) (string, string) {
	var datetime string
	var ip string
	var mac string
	fmt.Sscanf(line, "%s %s %s", &datetime, &ip, &mac)
	return ip, mac
}
