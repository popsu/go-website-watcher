package producer

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type WebsiteConfig struct {
	URL       string
	RePattern string
}

func WebsiteConfigFromFile(filename string) ([]WebsiteConfig, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)
	scanner := bufio.NewScanner(buf)

	var c []WebsiteConfig

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// comment line
			continue
		}

		wc := WebsiteConfig{}

		splitted := strings.SplitN(line, " ", 2)

		wc.URL = splitted[0]

		if len(splitted) == 2 {
			wc.RePattern = splitted[1]
		}

		c = append(c, wc)
	}

	return c, nil
}
