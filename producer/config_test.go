package producer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		filename string
		want     []WebsiteConfig
	}{
		{
			desc:     "testdata",
			filename: "testdata/testconfig.txt",
			want: []WebsiteConfig{
				{URL: "https://www.google.com", RePattern: ""},
				{URL: "https://golang.org", RePattern: `put"><noscript>Hello, 世界</`},
				{URL: "https://www.hs.fi", RePattern: "HS Digillä"},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := WebsiteConfigFromFile(tt.filename)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
