package qrclient

import (
	"time"

	knot "github.com/eaciit/knot/knot.v1"
	"github.com/ranggadablues/qrreq"
)

var (
	_qrConfig   Config
	_loginToken *tokenMap
)

type Config struct {
	App                *knot.App
	AuthCallback       func(*qrreq.LoginRequest) AuthResult
	LoginCallback      func(string, *knot.WebContext) interface{}
	ForgotPassCallback func(*qrreq.ForgotPassRequest) qrreq.ForgotPassResponse
	// ClientKey      string
	// clientEcdsaKey *ecdsa.PublicKey
}

type AuthResult struct {
	Success   bool
	Mobile    string
	Email     string
	Name      string
	LastLogin time.Time
}

// func initkey(cfg *Config) error {
// 	fullkey := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", cfg.ClientKey)
// 	block, _ := pem.Decode([]byte(fullkey))
// 	if block == nil {
// 		return errors.New("error decoding pem for qr client key")
// 	}
// 	genericPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		return errors.New("error parsing PKIX client key")
// 	}

// 	publicKey := genericPublicKey.(*ecdsa.PublicKey)

// 	cfg.clientEcdsaKey = publicKey

// 	return nil
// }

func Configure(config Config) error {
	_qrConfig = config
	_loginToken = newTokenMap()

	_loginToken.start()

	// err := initkey(&_qrConfig)
	// if err != nil {
	// 	return err
	// }

	err := _qrConfig.App.Register(new(QRClientController))
	if err != nil {
		return err
	}

	return nil
}
