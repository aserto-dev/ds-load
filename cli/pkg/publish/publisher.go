package publish

import (
	"context"
	"io"

	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
)

type DirectoryPublisher struct {
	UI             *clui.UI
	Log            *zerolog.Logger
	importerClient dsi.ImporterClient
}

func NewDirectoryPublisher(commonCtx *cc.CommonCtx, importerClient dsi.ImporterClient) *DirectoryPublisher {
	return &DirectoryPublisher{
		UI:             commonCtx.UI,
		Log:            commonCtx.Log,
		importerClient: importerClient,
	}
}

func (p *DirectoryPublisher) Publish(ctx context.Context, reader io.Reader) error {
	jsonReader, err := js.NewJSONArrayReader(reader)
	if err != nil {
		return err
	}

	for {
		var message msg.Transform
		err := jsonReader.ReadProtoMessage(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = p.publishMessages(ctx, &message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DirectoryPublisher) publishMessages(ctx context.Context, message *msg.Transform) error {
	var sErr error
	errGroup, iCtx := errgroup.WithContext(ctx)
	stream, err := p.importerClient.Import(iCtx)
	if err != nil {
		return err
	}
	errGroup.Go(receiver(stream))

	// import objects
	for _, object := range message.Objects {
		p.UI.Note().Msgf("object: [%s] type [%s]", object.Key, object.Type)
		sErr = stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Object{
				Object: object,
			},
		})
	}

	// import relations
	for _, relation := range message.Relations {
		p.UI.Note().Msgf("relation: [%s] obj: [%s] subj [%s]", relation.Relation, *relation.Object.Key, *relation.Subject.Key)
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
		p.Log.Err(sErr)
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
