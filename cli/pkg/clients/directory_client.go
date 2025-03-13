package clients

import (
	"context"

	"github.com/aserto-dev/ds-load/sdk/common/version"
	grpcClient "github.com/aserto-dev/go-aserto"
	"github.com/fullstorydev/grpcurl"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	dsiv3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
)

const localhostDirectory = "localhost:9292"

type Config struct {
	Host      string `short:"s" env:"DIRECTORY_HOST" help:"Directory host address"`
	APIKey    string `short:"k" env:"DIRECTORY_API_KEY" help:"Directory API Key"`
	Insecure  bool   `short:"i" help:"Disable TLS verification" xor:"tls"`
	Plaintext bool   `help:"Use plaintext connection" xor:"tls"`
	TenantID  string `short:"t" env:"DIRECTORY_TENANT_ID" help:"Directory Tenant ID"`
}

func NewDirectoryV3ImportClient(ctx context.Context, cfg *Config) (dsiv3.ImporterClient, error) {
	if cfg.Host == "" {
		cfg.Host = localhostDirectory
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	opts := []grpcClient.ConnectionOption{
		grpcClient.WithAddr(cfg.Host),
		grpcClient.WithInsecure(cfg.Insecure),
		grpcClient.WithNoTLS(cfg.Plaintext),
	}

	if cfg.APIKey != "" {
		opts = append(opts, grpcClient.WithAPIKeyAuth(cfg.APIKey))
	}

	if cfg.TenantID != "" {
		opts = append(opts, grpcClient.WithTenantID(cfg.TenantID))
	}

	conn, err := grpcClient.NewConnection(opts...)
	if err != nil {
		return nil, err
	}

	return dsiv3.NewImporterClient(conn), nil
}

func validate(cfg *Config) error {
	ctx := context.Background()

	tlsConf, err := grpcurl.ClientTLSConfig(cfg.Insecure, "", "", "")
	if err != nil {
		return errors.Wrap(err, "failed to create TLS config")
	}

	var creds credentials.TransportCredentials
	if !cfg.Plaintext {
		creds = credentials.NewTLS(tlsConf)
	}

	opts := []grpc.DialOption{
		grpc.WithUserAgent("ds-load " + version.GetInfo().Version),
	}

	if cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if _, err := grpcurl.BlockingDial(ctx, "tcp", cfg.Host, creds, opts...); err != nil {
		return err
	}

	return nil
}
