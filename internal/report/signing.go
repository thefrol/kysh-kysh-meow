package report

var signingKey []byte

func SetSigningKey(key string) {
	signingKey = []byte(key)
}
