package findpackagesrc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleSrc = `go.opencensus.io v0.15.0/go.mod h1:UffZAU+4sDEINUGP/B7UfBBkq4fqLu9zXAX7ke6CHW0=
go.opencensus.io v0.21.0/go.mod h1:mSImk1erAIZhrmZN+AvHh14ztQfjbGwt4TtuofqLduU=
go.opencensus.io v0.22.0 h1:C9hSCOW830chIVkdja34wa6Ky+IzWllkUinR+BtRZd4=
go.opencensus.io v0.22.0/go.mod h1:+kGneAE2xo2IficOXnaByMWTGM9T73dGwxeWcUqIpI8=
gocloud.dev v0.17.0 h1:UuDiCphYsiNhRNLtgHVL/eZheQeCt00hL3XjDfbt820=
gocloud.dev v0.17.0/go.mod h1:tIHTRdR1V5dlD8sTkzYdTGizBJ314BDykJ8KmadEXwo=
github.com/GoogleCloudPlatform/cloudsql-proxy v0.0.0-20190605020000-c4ba1fdf4d36/go.mod h1:aJ4qN3TfrelA6NZ6AXsXRfmEVaYin3EDbSPJrKS8OXo=
github.com/OneOfOne/xxhash v1.2.5 h1:zl/OfRA6nftbBK9qTohYBJ5xvw6C/oNKizR7cZGl3cI=
github.com/OneOfOne/xxhash v1.2.5/go.mod h1:eZbhyaAYD41SGSSsnmcpxVoRiQ/MPUTjUdIIOT9Um7Q=
`

func Test_parseGoSum(t *testing.T) {
	result, err := parseGoSum(strings.NewReader(sampleSrc))
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, "go.opencensus.io@v0.22.0", result["go.opencensus.io"])
	assert.Equal(t, "github.com/!one!of!one/xxhash@v1.2.5", result["github.com/OneOfOne/xxhash"])
}
