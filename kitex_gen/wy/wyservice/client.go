// Code generated by Kitex v0.13.0. DO NOT EDIT.

package wyservice

import (
	"context"
	wy "github.com/ChangKaiChen/lx-music-api-server/kitex_gen/wy"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	WyMusicUrl(ctx context.Context, req *wy.WyRequest, callOptions ...callopt.Option) (r *wy.WyResponse, err error)
}

// NewClient creates a start for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kWyServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a start for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kWyServiceClient struct {
	*kClient
}

func (p *kWyServiceClient) WyMusicUrl(ctx context.Context, req *wy.WyRequest, callOptions ...callopt.Option) (r *wy.WyResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.WyMusicUrl(ctx, req)
}
