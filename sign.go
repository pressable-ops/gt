package gt

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"net/url"
	"strings"
)

// ErrInvalidPost signifies that one of the post strings did not conform to key=value
var ErrInvalidPost = fmt.Errorf("One or more POST pairs was invalid")

// ErrWritingHasher is an error that should not normally occur and, perhaps, should generate a panic.
// the hash.Hash returned an error from it's Write() function fall
var ErrWritingHasher = fmt.Errorf("An error occured attempting to generate a sha256 hmac")

// Sign allows you to sign the uri and post parameters with the given PSK for use in authenticating
// a Guided transfer HTTP request.  The post slice should be in the form of "key=value".  the Sign
// function will split the key, and value, and do any necessary url encoding, so do not pre urlencode
// your data for this function.
func Sign(psk, uri string, post []string) ([]byte, error) {
	var signString = ""
	var hasher = hmac.New(sha256.New, []byte(psk))

	if uri[:1] != "/" {
		uri = "/" + uri
	}
	if len(post) > 0 {
		for i, v := range post {
			vs := strings.SplitN(v, "=", 2)
			if len(vs) != 2 {
				log.Printf("%#v", v)
				log.Printf("%#v", vs)
				return nil, ErrInvalidPost
			}
			vs[0] = url.QueryEscape(vs[0])
			vs[1] = url.QueryEscape(vs[1])
			post[i] = strings.Join(vs, "=")
		}
		signString = url.QueryEscape(fmt.Sprintf("%s?%s", uri, strings.Join(post, "&")))
	} else {
		signString = url.QueryEscape(uri)
	}
	log.Println("signing:", signString)
	if _, err := hasher.Write([]byte(signString)); err != nil {
		return nil, ErrWritingHasher
	}
	return hasher.Sum(nil), nil
}
