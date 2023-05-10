// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: dfsSrv.proto

package pb

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for DfsSrv service

func NewDfsSrvEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for DfsSrv service

type DfsSrvService interface {
	Upload(ctx context.Context, in *Args, opts ...client.CallOption) (*Result, error)
}

type dfsSrvService struct {
	c    client.Client
	name string
}

func NewDfsSrvService(name string, c client.Client) DfsSrvService {
	return &dfsSrvService{
		c:    c,
		name: name,
	}
}

func (c *dfsSrvService) Upload(ctx context.Context, in *Args, opts ...client.CallOption) (*Result, error) {
	req := c.c.NewRequest(c.name, "DfsSrv.Upload", in)
	out := new(Result)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DfsSrv service

type DfsSrvHandler interface {
	Upload(context.Context, *Args, *Result) error
}

func RegisterDfsSrvHandler(s server.Server, hdlr DfsSrvHandler, opts ...server.HandlerOption) error {
	type dfsSrv interface {
		Upload(ctx context.Context, in *Args, out *Result) error
	}
	type DfsSrv struct {
		dfsSrv
	}
	h := &dfsSrvHandler{hdlr}
	return s.Handle(s.NewHandler(&DfsSrv{h}, opts...))
}

type dfsSrvHandler struct {
	DfsSrvHandler
}

func (h *dfsSrvHandler) Upload(ctx context.Context, in *Args, out *Result) error {
	return h.DfsSrvHandler.Upload(ctx, in, out)
}
