package main

import (
	"context"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/yuanyu90221/go-twirp-rpc-example/rpc/notes"
)

type notesService struct {
	Notes     []notes.Note
	CurrentId int32
}

func main() {
	notesServer := notes.NewNotesServiceServer(&notesService{})

	mux := http.NewServeMux()
	mux.Handle(notesServer.PathPrefix(), notesServer)

	err := http.ListenAndServe(":8080", notesServer)
	if err != nil {
		panic(err)
	}
}

func (s *notesService) CreateNote(ctx context.Context,
	params *notes.CreateNoteParams) (*notes.Note, error) {
	if len(params.Text) < 4 {
		return nil, twirp.InvalidArgument.Error("Text should be min 4 characters")
	}
	note := notes.Note{
		Id:        s.CurrentId,
		Text:      params.Text,
		CreatedAt: time.Now().UnixMilli(),
	}
	s.Notes = append(s.Notes, note)

	s.CurrentId++
	return &note, nil
}

func (s *notesService) GetAllNotes(
	ctx context.Context,
	params *notes.GetAllNotesParams,
) (*notes.GetAllNotesResult, error) {
	allNotes := make([]*notes.Note, 0)
	for _, note := range s.Notes {
		n := note
		allNotes = append(allNotes, &n)
	}
	return &notes.GetAllNotesResult{
		Notes: allNotes,
	}, nil
}
