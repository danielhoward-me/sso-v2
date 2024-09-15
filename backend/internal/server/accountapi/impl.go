package accountapi

import (
	"context"
)

type Api struct{}

var _ StrictServerInterface = (*Api)(nil)

func New() *Api {
	return &Api{}
}

func (Api) GetTest(ctx context.Context, request GetTestRequestObject) (resp GetTestResponseObject, err error) {
	return GetTest200JSONResponse("Hello World"), nil
}
