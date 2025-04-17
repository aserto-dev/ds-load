package publish

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	dsi3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
	"github.com/aserto-dev/go-directory/pkg/validator"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
)

const (
	objectsCounter   string = "object"
	relationsCounter string = "relation"
)

type DirectoryPublisher struct {
	Log            *zerolog.Logger
	importerClient dsi3.ImporterClient
	errs           bool
	objCounter     *dsi3.ImportCounter
	relCounter     *dsi3.ImportCounter
}

func NewDirectoryPublisher(commonCtx *cc.CommonCtx, importerClient dsi3.ImporterClient) *DirectoryPublisher {
	return &DirectoryPublisher{
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
		if errors.Is(err, io.EOF) {
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

	if p.objCounter != nil {
		printCounter(os.Stdout, p.objCounter)
	}

	if p.relCounter != nil {
		printCounter(os.Stdout, p.relCounter)
	}

	if p.errs {
		fmt.Fprintf(os.Stderr, "import failure\n")
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

	opCode := message.GetOpCode()
	if opCode == dsi3.Opcode_OPCODE_UNKNOWN {
		opCode = dsi3.Opcode_OPCODE_SET
	}

	// import objects
	for _, object := range message.GetObjects() {
		if err := validator.Object(object); err != nil {
			fmt.Fprintf(os.Stderr, "validation failed, object: [%s] type [%s]\n", object.GetId(), object.GetType())
			continue
		}

		if (opCode == dsi3.Opcode_OPCODE_DELETE || opCode == dsi3.Opcode_OPCODE_DELETE_WITH_RELATIONS) && object.GetType() == "group" {
			continue
		}

		fmt.Fprintf(os.Stdout, "object: [%s] type [%s]\n", object.GetId(), object.GetType())
		sErr := stream.Send(&dsi3.ImportRequest{
			Msg: &dsi3.ImportRequest_Object{
				Object: object,
			},
			OpCode: opCode,
		})

		p.handleStreamError(sErr)
	}

	// import relations
	for _, relation := range message.GetRelations() {
		if err := validator.Relation(relation); err != nil {
			fmt.Fprintf(os.Stderr, "validation failed, relation: [%s] obj: [%s] subj [%s]\n",
				relation.GetRelation(), relation.GetObjectId(), relation.GetSubjectId())
			continue
		}

		fmt.Fprintf(os.Stdout, "relation: [%s] obj: [%s] subj [%s]\n", relation.GetRelation(), relation.GetObjectId(), relation.GetSubjectId())
		sErr := stream.Send(&dsi3.ImportRequest{
			Msg: &dsi3.ImportRequest_Relation{
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

	p.Log.Err(err).Msg("stream error")
	common.SetExitCode(1)
}

func (p *DirectoryPublisher) receiver(stream dsi3.Importer_ImportClient) func() error {
	return func() error {
		for {
			result, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}

			if err != nil {
				return err
			}

			if result != nil {
				switch m := result.GetMsg().(type) {
				case *dsi3.ImportResponse_Status:
					p.errs = true

					printStatus(os.Stderr, m.Status)
				case *dsi3.ImportResponse_Counter:
					switch m.Counter.GetType() {
					case objectsCounter:
						p.objCounter = m.Counter
					case relationsCounter:
						p.relCounter = m.Counter
					}
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

func printStatus(w io.Writer, status *dsi3.ImportStatus) {
	fmt.Fprintf(w, "%-9s : %s - %s (%d)\n",
		"error",
		status.GetMsg(),
		codes.Code(status.GetCode()).String(),
		status.GetCode())
}

func printCounter(w io.Writer, ctr *dsi3.ImportCounter) {
	fmt.Fprintf(w, "%-9s : %d (set:%d delete:%d error:%d)\n",
		ctr.GetType(),
		ctr.GetRecv(),
		ctr.GetSet(),
		ctr.GetDelete(),
		ctr.GetError(),
	)
}
