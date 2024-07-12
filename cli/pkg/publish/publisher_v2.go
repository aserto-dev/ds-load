package publish

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/aserto-dev/go-directory/pkg/convert"
	"github.com/bufbuild/protovalidate-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
)

type DirectoryV2Publisher struct {
	Log            *zerolog.Logger
	importerClient dsi.ImporterClient
	validator      *protovalidate.Validator
	objErr         int
	relErr         int
}

func NewDirectoryV2Publisher(commonCtx *cc.CommonCtx, importerClient dsi.ImporterClient) *DirectoryV2Publisher {
	v, _ := protovalidate.New()

	return &DirectoryV2Publisher{
		Log:            commonCtx.Log,
		importerClient: importerClient,
		validator:      v,
	}
}

func (p *DirectoryV2Publisher) Publish(ctx context.Context, reader io.Reader) error {
	jsonReader, err := js.NewJSONArrayReader(reader)
	if err != nil {
		return err
	}

	for {
		var message msg.Transform
		v2msg := msg.TransformV2{
			Objects:   []*v2.Object{},
			Relations: []*v2.Relation{},
		}
		err := jsonReader.ReadProtoMessage(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for _, object := range message.Objects {
			if err := p.validator.Validate(object); err != nil {
				fmt.Fprintf(os.Stderr, "validation failed, object: [%s] type [%s]\n", object.Id, object.Type)
				continue
			}
			v2msg.Objects = append(v2msg.Objects, convert.ObjectToV2(object))
		}

		for _, relation := range message.Relations {
			if relation.SubjectRelation != "" {
				fmt.Fprintf(os.Stderr, "detected subject relation %s in v2 mode\n", relation.SubjectRelation)
				continue
			}

			if err := p.validator.Validate(relation); err != nil {
				fmt.Fprintf(os.Stderr, "validation failed, relation: [%s] obj: [%s] subj [%s]\n", relation.Relation, relation.ObjectId, relation.SubjectId)
				continue
			}

			v2msg.Relations = append(v2msg.Relations, convert.RelationToV2(relation))
		}

		err = p.publishMessages(ctx, &v2msg)
		if err != nil {
			return err
		}
	}

	if p.objErr != 0 {
		fmt.Fprintf(os.Stderr, "failed to import %d objects\n", p.objErr)
		common.SetExitCode(1)
	}
	if p.relErr != 0 {
		fmt.Fprintf(os.Stderr, "failed to import %d relations\n", p.relErr)
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
	errGroup.Go(p.doneHandler(stream.Context()))

	// import objects
	for _, object := range message.Objects {
		fmt.Fprintf(os.Stdout, "object: [%s] type [%s]\n", object.Key, object.Type)
		sErr := stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Object{
				Object: object,
			},
		})
		p.handleStreamError(sErr)
	}

	// import relations
	for _, relation := range message.Relations {
		fmt.Fprintf(os.Stdout, "relation: [%s] obj: [%s] subj [%s]\n", relation.Relation, *relation.Object.Key, *relation.Subject.Key)
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

			if err != nil {
				return err
			}

			if result != nil {
				if result.Object != nil && result.Object.Error != 0 {
					p.objErr += int(result.Object.Error)
				}
				if result.Relation != nil && result.Relation.Error != 0 {
					p.relErr += int(result.Relation.Error)
				}
			}
		}
	}
}

func (p *DirectoryV2Publisher) doneHandler(ctx context.Context) func() error {
	return func() error {
		<-ctx.Done()
		err := ctx.Err()
		if err != nil && !errors.Is(err, context.Canceled) {
			p.Log.Trace().Err(err).Msg("subscriber-doneHandler")
			return err
		}
		return nil
	}
}
