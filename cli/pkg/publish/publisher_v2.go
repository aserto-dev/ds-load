package publish

import (
	"context"
	"io"

	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/aserto-dev/go-directory/pkg/convert"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
)

type DirectoryV2Publisher struct {
	UI             *clui.UI
	Log            *zerolog.Logger
	importerClient dsi.ImporterClient
	objErr         int
	relErr         int
}

func NewDirectoryV2Publisher(commonCtx *cc.CommonCtx, importerClient dsi.ImporterClient) *DirectoryV2Publisher {
	return &DirectoryV2Publisher{
		UI:             commonCtx.UI,
		Log:            commonCtx.Log,
		importerClient: importerClient,
	}
}

func (p *DirectoryV2Publisher) Publish(ctx context.Context, reader io.Reader) error {
	jsonReader, err := js.NewJSONArrayReader(reader)
	if err != nil {
		return err
	}

	for {
		var message msg.Transform
		var v2msg msg.TransformV2
		err := jsonReader.ReadProtoMessage(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for _, r := range message.Relations {
			if r.SubjectRelation != "" {
				p.UI.Problem().Msgf("detected subject relation %s in v2 mode", r.SubjectRelation)
			}
		}
		v2msg.Objects = convert.ObjectArrayToV2(message.Objects)
		v2msg.Relations = convert.RelationArrayToV2(message.Relations)

		err = p.publishMessages(ctx, &v2msg)
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

func (p *DirectoryV2Publisher) publishMessages(ctx context.Context, message *msg.TransformV2) error {
	errGroup, iCtx := errgroup.WithContext(ctx)
	stream, err := p.importerClient.Import(iCtx)
	if err != nil {
		return err
	}
	errGroup.Go(p.receiver(stream))

	// import objects
	for _, object := range message.Objects {
		p.UI.Note().Msgf("object: [%s] type [%s]", object.Key, object.Type)
		sErr := stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Object{
				Object: object,
			},
		})
		p.handleStreamError(sErr)
	}

	// import relations
	for _, relation := range message.Relations {
		p.UI.Note().Msgf("relation: [%s] obj: [%s] subj [%s]", relation.Relation, *relation.Object.Key, *relation.Subject.Key)
		sErr := stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Relation{
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

func (p *DirectoryV2Publisher) handleStreamError(err error) {
	if err == nil {
		return
	}
	p.Log.Err(err)
	common.SetExitCode(1)
}

func (p *DirectoryV2Publisher) receiver(stream dsi.Importer_ImportClient) func() error {
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
