package fdumper

import (
	"context"
	"errors"
	"fmt"
)

type (
	FnMarshal   func(instance interface{}) ([]byte, error)
	FnUnmarshal func(data []byte, out interface{}) error
	FnFetch     func(ctx context.Context, instance interface{}) error

	Dumper struct {
		FnMarshal
		FnUnmarshal
		FnFetch

		Endpoint

		// Expire int64 `json:"expire"`
	}

	Endpoint interface {
		Dump(ctx context.Context, filepath string, data []byte) error
		Load(ctx context.Context, filepath string) ([]byte, error)
	}
)

var (
	ErrDumpFailed = errors.New("dump failed")
	ErrLoadFailed = errors.New("load failed")
)

func New(marshal FnMarshal, unmarshal FnUnmarshal, fetch FnFetch, ep Endpoint) *Dumper {
	return &Dumper{
		FnMarshal:   marshal,
		FnUnmarshal: unmarshal,
		FnFetch:     fetch,
		Endpoint:    ep,
	}
}

func (dumper *Dumper) Dump(ctx context.Context, filepath string, instance interface{}) error {
	if dumper.Endpoint == nil {
		return fmt.Errorf("there are no endpoint, %w", ErrDumpFailed)
	}
	if dumper.FnMarshal == nil {
		return fmt.Errorf("there are no unmarshaller, %w", ErrDumpFailed)
	}
	data, err := dumper.FnMarshal(instance)
	if err != nil {
		return fmt.Errorf("marshal failed, %w", ErrDumpFailed)
	}
	return dumper.Endpoint.Dump(ctx, filepath, data)
}

func (dumper *Dumper) Load(ctx context.Context, filepath string, out interface{}) error {
	if dumper.Endpoint == nil {
		return fmt.Errorf("there are no endpoint, %w", ErrLoadFailed)
	}
	if dumper.FnUnmarshal == nil {
		return fmt.Errorf("there are no unmarshaller, %w", ErrLoadFailed)
	}
	data, err := dumper.Endpoint.Load(ctx, filepath)
	if err != nil {
		return err
	}
	return dumper.FnUnmarshal(data, out)
}

func (dumper *Dumper) Fetch(ctx context.Context, out interface{}, maxRetry int) error {
	if dumper.FnFetch == nil {
		return fmt.Errorf("got empty fetch method")
	}
	err := dumper.FnFetch(ctx, out)
	for retry := 0; err != nil && retry < maxRetry; retry++ {
		fmt.Println("retry ", retry)
		err = dumper.FnFetch(ctx, out)
	}
	if err != nil {
		return fmt.Errorf("fetching failed, %w", err)
	}
	return err
}

func (dumper *Dumper) FetchAndDump(ctx context.Context, filepath string, out interface{}, maxRetry int) error {
	if err := dumper.Fetch(ctx, out, maxRetry); err != nil {
		return err
	}
	return dumper.Dump(ctx, filepath, out)
}

func (dumper *Dumper) LoadOrFetch(ctx context.Context, fileName string, out interface{}, maxRetry int) error {
	if err := dumper.Load(ctx, fileName, out); err == nil {
		return nil
	}
	return dumper.Fetch(ctx, out, maxRetry)
}

func (dumper *Dumper) LoadOrFetchAndDump(ctx context.Context, fileName string, out interface{}, maxRetry int) error {
	if err := dumper.Load(ctx, fileName, out); err == nil {
		return nil
	}
	return dumper.FetchAndDump(ctx, fileName, out, maxRetry)
}
