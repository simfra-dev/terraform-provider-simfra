package awsclient

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type Client struct {
	Route53       *route53.Client
	Organizations *organizations.Client
}

type Config struct {
	Endpoint     string
	AccessKey    string
	SecretKey    string
	Region       string
	SkipTLSVerify bool
}

type overridesKey struct{}

// WithOverrides attaches X-Simfra-* header overrides to a context for a single API call.
func WithOverrides(ctx context.Context, overrides map[string]string) context.Context {
	return context.WithValue(ctx, overridesKey{}, overrides)
}

func injectOverrides(ctx context.Context, in smithymiddleware.BuildInput, next smithymiddleware.BuildHandler) (smithymiddleware.BuildOutput, smithymiddleware.Metadata, error) {
	if overrides, ok := ctx.Value(overridesKey{}).(map[string]string); ok {
		if req, ok := in.Request.(*smithyhttp.Request); ok {
			for k, v := range overrides {
				req.Header.Set("X-Simfra-"+k, v)
			}
		}
	}
	return next.HandleBuild(ctx, in)
}

// HeaderMiddleware registers the X-Simfra-* header injection middleware on the stack.
func HeaderMiddleware(stack *smithymiddleware.Stack) error {
	return stack.Build.Add(smithymiddleware.BuildMiddlewareFunc("SimfraOverrides", injectOverrides), smithymiddleware.Before)
}

func New(cfg Config) *Client {
	httpClient := &http.Client{}
	if cfg.SkipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}
	}

	awsCfg := aws.Config{
		Region:      cfg.Region,
		Credentials: credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		HTTPClient:  httpClient,
	}

	return &Client{
		Route53: route53.NewFromConfig(awsCfg, func(o *route53.Options) {
			o.BaseEndpoint = &cfg.Endpoint
			o.APIOptions = append(o.APIOptions, HeaderMiddleware)
		}),
		Organizations: organizations.NewFromConfig(awsCfg, func(o *organizations.Options) {
			o.BaseEndpoint = &cfg.Endpoint
			o.APIOptions = append(o.APIOptions, HeaderMiddleware)
		}),
	}
}
