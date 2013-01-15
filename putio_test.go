package putio

import (
	"testing"
)

var (
	putio_clientid    = "263"
	putio_appsecret   = "7vtklth37axmtlpooaxf"
	putio_callbackurl = "https://github.com/bryon/put.io.go"
	user_code         = "cbb13b625b5511e29620001018321b64"
)

func TestPutio(t *testing.T) {
	p, _ := NewPutio(putio_clientid, putio_appsecret, putio_callbackurl, user_code)
	exstr := "ABV9KDHN"
	if p.OauthToken != exstr {
		t.Errorf("OAuth token appears invalid.  Expected: %s got :%s", exstr, p.OauthToken)
	}
}
