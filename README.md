# go rpc example with twirp package[^1]

Twirp is an RPC framework from Twitch which just like gRPC uses Protobufs and is much easier to use.

## Create Project

```shell=
mkdir $your_project_name
cd $your_project_name
go mod init github.com/$your_github_name/$your_project_name
```
project layout

```shell=
go-twirp-rpc-example/
|- client/
|- server/
|- rpc/
   |- notes/
```

* client: this folder place rpc client connect code
* server: this folder place rpc server implementation code
* rpc/notes: this folder place proto relative code

## cli tool for compile proto

```shell=
go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
## install protobuf and twirp module 
```shell=
go get github.com/twitchtv/twirp
go get google.golang.org/protobuf
```

## Step 1: Create Proto for communicate

```proto=
syntax = "proto3";

package yuanyu;
option go_package = "rpc/notes";

message Note {
  int32 id = 1;
  string text = 2;
  int64 created_at = 3;
}

message CreateNoteParams {
  string text = 1;
}

message GetAllNotesParams {
}

message GetAllNotesResult {
  repeated Note notes = 1;
}

service NotesService {
  rpc CreateNote(CreateNoteParams) returns (Note);
  rpc GetAllNotes(GetAllNotesParams) returns (GetAllNotesResult);
}
```

## Step2 Generating Go code for proto

```shell=
protoc --twirp_out=. --go_out=. rpc/notes/service.proto
```

***Notice*** the path for proto files must be correct

## Step3 implement the service code

1. Define a struct noteService to hold a list of notes

```go
import (
	"github.com/yuanyu90221/go-twirp-rpc-example/rpc/notes"
)

type notesService struct {
	Notes     []notes.Note
	CurrentId int32
}
```
2. Implement the CreateNote function on notesService 

```go
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
```

3. Implement the GetAllNotes function on notesService 

```go
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
```

4. Implement HTTP Server for serve over HTTP

```go
func main() {
	notesServer := notes.NewNotesServiceServer(&notesService{})

	mux := http.NewServeMux()
	mux.Handle(notesServer.PathPrefix(), notesServer)

	err := http.ListenAndServe(":8080", notesServer)
	if err != nil {
		panic(err)
	}
}
```

## Write Client Code to connect use grpc

1. setup connect client

```go
client := notes.NewNotesServiceProtobufClient("http://localhost:8080",
		&http.Client{})
```

2. call createNote

```go
ctx := context.Background()

_, err := client.CreateNote(ctx, &notes.CreateNoteParams{Text: "Hello World"})
if err != nil {
  log.Fatal(err)
}
```
3. call GetAllNotes

```go
allNotes, err := client.GetAllNotes(ctx, &notes.GetAllNotesParams{})
if err != nil {
  log.Fatal(err)
}
for _, note := range allNotes.Notes {
  log.Println(note)
}
```
## call rpc over HTTP POST
```shell
curl --request POST   --url http://localhost:8080/twirp/yuanyu.NotesService/GetAllNotes   --header 'Content-Type: application/json'   --data '{}'
```
## Reference

[^1]: https://thedevelopercafe.com/articles/rpc-in-go-using-twitchs-twirp-3dcb78ece775