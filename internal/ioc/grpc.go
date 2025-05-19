package ioc

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	notificationv1 "github.com/JrMarcco/jotice-api/api/notification/v1"
	grpcapi "github.com/JrMarcco/jotice/internal/api/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var GrpcFxOpt = fx.Provide(
	fx.Provide(NewGrpcServer),
	fx.Invoke(RunGrpcServer),
)

func NewGrpcServer(server grpcapi.NotificationServer, etcdClient *clientv3.Client) *grpc.Server {
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

	svr := grpc.NewServer()
	notificationv1.RegisterNotificationServiceServer(svr, server)
	notificationv1.RegisterNotificationQueryServiceServer(svr, server)

	return svr
}

func RunGrpcServer(lc fx.Lifecycle, grpcSvr *grpc.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// TODO
			//lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
			//if err != nil {
			//	return err
			//}
			//
			//go func() {
			//	if err := grpcSvr.Serve(lis); err != nil {
			//		panic(err)
			//	}
			//}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// TODO logging
			grpcSvr.GracefulStop()
			return nil
		},
	})
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
