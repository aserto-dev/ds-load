package clients

import (
	"context"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/common/version"
	grpcClient "github.com/aserto-dev/go-aserto/client"
	"github.com/fullstorydev/grpcurl"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
)

const localhostDirectory = "localhost:9292"

type Config struct {
	Host     string
	APIKey   string
	Insecure bool `short:"i" help:"Disable TLS verification"`
	TenantID string
}

func NewDirectoryImportClient(c *cc.CommonCtx, cfg *Config) (dsi.ImporterClient, error) {
	if cfg.Host == "" {
		cfg.Host = localhostDirectory
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	opts := []grpcClient.ConnectionOption{
		grpcClient.WithAddr(cfg.Host),
		grpcClient.WithInsecure(cfg.Insecure),
	}

	if cfg.APIKey != "" {
		opts = append(opts, grpcClient.WithAPIKeyAuth(cfg.APIKey))
	}

	if cfg.TenantID != "" {
		opts = append(opts, grpcClient.WithTenantID(cfg.TenantID))
	}

	conn, err := grpcClient.NewConnection(c.Context, opts...)
	if err != nil {
		return nil, err
	}

	return dsi.NewImporterClient(conn.Conn), nil
}

func validate(cfg *Config) error {
	ctx := context.Background()

	tlsConf, err := grpcurl.ClientTLSConfig(cfg.Insecure, "", "", "")
	if err != nil {
		return errors.Wrap(err, "failed to create TLS config")
	}

	creds := credentials.NewTLS(tlsConf)

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