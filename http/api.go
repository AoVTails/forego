package http

import (
	"encoding/json"
	"net/http"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

type Doable interface {
	Do(ctx.C) error
}

func RegisterAPI[T Doable](c ctx.C, s *Server, obj T) error {
	handler, err := api.NewServer(c, obj)
	if err != nil {
		return err
	}
	f := func(c ctx.C, in []byte, r *http.Request) ([]byte, error) {
		req := &api.JSON{}
		err := json.Unmarshal(in, &req.Data)
		// TODO auth

		obj, err := handler.Recv(c, req)
		if err != nil {
			return nil, err
		}
		err = obj.Do(c)
		if err != nil {
			return nil, err
		}

		res := &api.JSON{}
		err = handler.Send(c, obj, res)
		if err != nil {
			return nil, err
		}
		return json.Marshal(res.Data)
	}

	urls := handler.URLs()
	if len(urls) == 0 {
		return ctx.NewErrorf(c, "no URL to register for %T", obj)
	}

	for _, u := range handler.URLs() {
		log.Debugf(c, "registering to %q", u.Path)
		s.HandleRequest(u.Path, f)
	}
	return nil
}
