package publish

import (
	"context"
	"io"

	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	dsiv3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
)

type DirectoryPublisher struct {
	UI             *clui.UI
	Log            *zerolog.Logger
	importerClient dsiv3.ImporterClient
	objErr         int
	relErr         int
}

func NewDirectoryPublisher(commonCtx *cc.CommonCtx, importerClient dsiv3.ImporterClient) *DirectoryPublisher {
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

	if p.objErr != 0 {
		p.UI.Problem().Msgf("failed to import %d objects", p.objErr)
		common.SetExitCode(1)
	}
	if p.relErr != 0 {
		p.UI.Problem().Msgf("failed to import %d relations", p.relErr)
		common.SetExitCode(1)
	}

	return nil
}

func (p *DirectoryPublisher) publishMessages(ctx context.Context, message *msg.Transform) error {
	errGroup, iCtx := errgroup.WithContext(ctx)
	stream, err := p.importerClient.Import(iCtx)
	if err != nil {
		return err
	}
	errGroup.Go(p.receiver(stream))

	// import objects
	for _, object := range message.Objects {
		p.UI.Note().Msgf("object: [%s] type [%s]", object.Id, object.Type)
		sErr := stream.Send(&dsiv3.ImportRequest{
			Msg: &dsiv3.ImportRequest_Object{
				Object: object,
			},
		})
		p.handleStreamError(sErr)
	}

	// import relations
	for _, relation := range message.Relations {
		p.UI.Note().Msgf("relation: [%s] obj: [%s] subj [%s]", relation.Relation, relation.ObjectId, relation.SubjectId)
		sErr := stream.Send(&dsiv3.ImportRequest{
			Msg: &dsiv3.ImportRequest_Relation{
				Relation: relation,
			},
		})
		p.handleStreamError(sErr)
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	err = errGroup.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (p *DirectoryPublisher) handleStreamError(err error) {
	if err == nil {
		return
	}

	p.Log.Err(err)
	common.SetExitCode(1)
}

func (p *DirectoryPublisher) receiver(stream dsiv3.Importer_ImportClient) func() error {
	return func() error {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if result.Object.Error != 0 {
				p.objErr += int(result.Object.Error)
			}
			if result.Relation.Error != 0 {
				p.relErr += int(result.Relation.Error)
			}

			if err != nil {
				return err
			}
		}
	}
}
