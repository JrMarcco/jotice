package ioc

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	grpcapi "github.com/JrMarcco/jotice/internal/api/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitGrpc(server grpcapi.NotificationServer, etcdClient *clientv3.Client) {
	type Config struct {
		priPem string `yaml:"private"`
		pubPem string `yaml:"public"`
	}
	var cfg Config

	if err := viper.UnmarshalKey("jwt", cfg); err != nil {
		panic(err)
	}

	// TODO
	//priKey, pubKey := loadJwtKey(cfg.priPem, cfg.pubPem)
	//jwtInterceptor := jwt.NewJwtAuth(priKey, pubKey).Build()

}

func loadJwtKey(priPem, pubPem string) (ed25519.PrivateKey, ed25519.PublicKey) {
	priKeyBlock, _ := pem.Decode([]byte(priPem))
	if priKeyBlock == nil {
		panic("failed to decode private key PEM")
	}

	// the PEM block itself is labeled public key, not specifically ed25519 public key.
	// all standard public key formats need to be handled by the x509 package first.
	// using ParsePKCS8PrivateKey/ParsePKIXPublicKey and then type-asserted into ed25519 PrivateKey/PublicKey.
	priKey, err := x509.ParsePKCS8PrivateKey(priKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	pubKeyBlock, _ := pem.Decode([]byte(pubPem))
	if pubKeyBlock == nil {
		panic("failed to decode public key PEM")
	}
	publicKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	return priKey.(ed25519.PrivateKey), publicKey.(ed25519.PublicKey)
}
