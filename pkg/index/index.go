package index

import (
	"github.com/elolpuer/Blog/pkg/models"
)

func Page() *models.IndexResp{
	var Resp = new(models.IndexResp)
	Resp.Index = "Index page"
	Resp.Title = "Index Page"
	return Resp
}
