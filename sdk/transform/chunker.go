package transform

import (
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
)

type Chunker interface {
	WriteChunks(writer *js.JSONArrayWriter, directoryObject *msg.Transform) error
}

type chunker struct {
	maxChunkSize int
}

func NewChunker(maxChunkSize int) Chunker {
	return &chunker{
		maxChunkSize: maxChunkSize,
	}
}

func (c *chunker) WriteChunks(writer *js.JSONArrayWriter, directoryObject *msg.Transform) error {
	chunks := c.prepareChunks(directoryObject)

	for _, chunk := range chunks {
		err := writer.WriteProtoMessage(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *chunker) prepareChunks(directoryObject *msg.Transform) []*msg.Transform {
	var chunks []*msg.Transform
	var freeChunk *msg.Transform

	for _, obj := range directoryObject.Objects {
		freeChunk, chunks = c.nextFreeChunk(chunks)
		freeChunk.Objects = append(freeChunk.Objects, obj)
	}

	for _, rel := range directoryObject.Relations {
		freeChunk, chunks = c.nextFreeChunk(chunks)
		freeChunk.Relations = append(freeChunk.Relations, rel)
	}

	return chunks
}

func (c *chunker) nextFreeChunk(chunks []*msg.Transform) (*msg.Transform, []*msg.Transform) {
	if len(chunks) == 0 {
		chunks = c.addEmptyChunk(chunks)
	}

	lastChunk := chunks[len(chunks)-1]
	if c.isRoomInChunk(lastChunk) {
		return lastChunk, chunks
	}

	chunks = c.addEmptyChunk(chunks)

	return chunks[len(chunks)-1], chunks
}

func (c *chunker) isRoomInChunk(chunk *msg.Transform) bool {
	return len(chunk.Objects)+len(chunk.Relations) < c.maxChunkSize
}

func (c *chunker) addEmptyChunk(chunks []*msg.Transform) []*msg.Transform {
	return append(chunks, &msg.Transform{Objects: []*v2.Object{}, Relations: []*v2.Relation{}})
}
