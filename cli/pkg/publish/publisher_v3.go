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
	"github.com/bufbuild/protovalidate-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	dsiv3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
)

type DirectoryPublisher struct {
	Log            *zerolog.Logger
	importerClient dsiv3.ImporterClient
	validator      *protovalidate.Validator
	objErr         int
	relErr         int
}

func NewDirectoryPublisher(commonCtx *cc.CommonCtx, importerClient dsiv3.ImporterClient) *DirectoryPublisher {
	v, _ := protovalidate.New()

	return &DirectoryPublisher{
		Log:            commonCtx.Log,
		importerClient: importerClient,
		validator:      v,
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
		fmt.Fprintf(os.Stderr, "failed to import %d objects\n", p.objErr)
		common.SetExitCode(1)
	}
	if p.relErr != 0 {
		fmt.Fprintf(os.Stderr, "failed to import %d relations\n", p.relErr)
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
	errGroup.Go(p.doneHandler(stream.Context()))

	opCode := message.OpCode
	if opCode == dsiv3.Opcode_OPCODE_UNKNOWN {
		opCode = dsiv3.Opcode_OPCODE_SET
	}

	// import objects
	for _, object := range message.Objects {
		if err := p.validator.Validate(object); err != nil {
			fmt.Fprintf(os.Stderr, "validation failed, object: [%s] type [%s]\n", object.Id, object.Type)
			continue
		}
		if (opCode == dsiv3.Opcode_OPCODE_DELETE || opCode == dsiv3.Opcode_OPCODE_DELETE_WITH_RELATIONS) && object.Type == "group" {
			continue
		}
		fmt.Fprintf(os.Stdout, "object: [%s] type [%s]\n", object.Id, object.Type)
		sErr := stream.Send(&dsiv3.ImportRequest{
			Msg: &dsiv3.ImportRequest_Object{
				Object: object,
			},
			OpCode: opCode,
		})
		p.handleStreamError(sErr)
	}

	// import relations
	for _, relation := range message.Relations {
		if err := p.validator.Validate(relation); err != nil {
			fmt.Fprintf(os.Stderr, "validation failed, relation: [%s] obj: [%s] subj [%s]\n", relation.Relation, relation.ObjectId, relation.SubjectId)
			continue
		}
		fmt.Fprintf(os.Stdout, "relation: [%s] obj: [%s] subj [%s]\n", relation.Relation, relation.ObjectId, relation.SubjectId)
		sErr := stream.Send(&dsiv3.ImportRequest{
			Msg: &dsiv3.ImportRequest_Relation{
				Relation: relation,
			},
			OpCode: opCode,
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

func (p *DirectoryPublisher) doneHandler(ctx context.Context) func() error {
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
