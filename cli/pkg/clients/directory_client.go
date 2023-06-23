package clients

import (
	"context"
	"io"

	"github.com/aserto-dev/ds/cli/pkg/cc"
	"github.com/aserto-dev/ds/sdk/common/js"
	"github.com/aserto-dev/ds/sdk/common/msg"
	"github.com/aserto-dev/ds/sdk/common/version"
	grpcClient "github.com/aserto-dev/go-aserto/client"
	"github.com/fullstorydev/grpcurl"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
)

const localhostDirectory = "localhost:9292"

type Config struct {
	Host     string `short:"s" env:"DIRECTORY_HOST" help:"Directory host address"`
	APIKey   string `short:"k" env:"DIRECTORY_API_KEY" help:"Directory API Key"`
	Insecure bool   `short:"i" help:"Disable TLS verification"`
	TenantID string `short:"t" env:"DIRECTORY_TENANT_ID" help:"Directory Tenant ID"`
}

type DirectoryClient interface {
	HandleMessages(stdout io.Reader) error
}

type directoryClient struct {
	commonCtx *cc.CommonCtx
	dirClient dsi.ImporterClient `kong:"-"`
}

func NewDirectoryImportClient(c *cc.CommonCtx, cfg *Config) (DirectoryClient, error) {
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

	return &directoryClient{
		dirClient: dsi.NewImporterClient(conn.Conn),
		commonCtx: c,
	}, nil
}

func validate(cfg *Config) error {
	ctx := context.Background()

	tlsConf, err := grpcurl.ClientTLSConfig(cfg.Insecure, "", "", "")
	if err != nil {
		return errors.Wrap(err, "failed to create TLS config")
	}

	creds := credentials.NewTLS(tlsConf)

	opts := []grpc.DialOption{
		grpc.WithUserAgent("ds " + version.GetInfo().Version),
	}
	if cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if _, err := grpcurl.BlockingDial(ctx, "tcp", cfg.Host, creds, opts...); err != nil {
		return err
	}
	return nil
}

func (d *directoryClient) HandleMessages(stdout io.Reader) error {
	reader, err := js.NewJSONArrayReader(stdout)
	if err != nil {
		return err
	}

	for {
		var message msg.Transform
		err := reader.ReadProtoMessage(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = d.importToDirectory(d.commonCtx, &message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *directoryClient) importToDirectory(c *cc.CommonCtx, message *msg.Transform) error {
	var sErr error
	errGroup, iCtx := errgroup.WithContext(c.Context)
	stream, err := d.dirClient.Import(iCtx)
	if err != nil {
		return err
	}
	errGroup.Go(receiver(stream))

	// import objects
	for _, object := range message.Objects {
		c.UI.Note().Msgf("object: [%s] type [%s]", object.Key, object.Type)
		sErr = stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Object{
				Object: object,
			},
		})
	}

	for _, relation := range message.Relations {
		c.UI.Note().Msgf("relation: [%s] obj: [%s] subj [%s]", relation.Relation, *relation.Object.Key, *relation.Subject.Key)
		sErr = stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Relation{
				Relation: relation,
			},
		})
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	err = errGroup.Wait()
	if err != nil {
		return err
	}

	// TODO handle stream errors
	if sErr != nil {
		c.Log.Err(sErr)
	}

	return nil
}

func receiver(stream dsi.Importer_ImportClient) func() error {
	return func() error {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}
		}
	}
}
