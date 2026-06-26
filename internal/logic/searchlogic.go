package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
	"github.com/aiqoder/monitor-lite-api/model"
	"github.com/aiqoder/monitor-lite-api/utils"
)

type SearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// domestic false 表示国外
func (l *SearchLogic) Search(req *types.TvJsonReq, domestic bool) ([]types.Tv, error) {
	var tvs []types.Tv

	if req.Mode == "so" {
		cTvs, err := utils.Search(req.TvName, "name")

		if err != nil {
			return nil, err
		}

		for i := 0; i < len(cTvs); i++ {
			tv := cTvs[i]

			//if domestic && utils.ChinaDisabledStr(tv.Name) {
			//	continue
			//}

			tvs = append(tvs, types.Tv{
				Name: tv.Name,
				Url:  tv.Url,
			})
		}
	} else if req.Mode == "me" {
		var list []model.Tv
		var err error

		if req.TvName == "auto check" {
			list, err = l.svcCtx.TvModel.NextTvOffset()
		} else {
			list, err = l.svcCtx.TvModel.TableListByTvName(req.TvName)
		}

		if err != nil {
			return nil, err
		}

		for i := 0; i < len(list); i++ {
			tv := list[i]

			tvs = append(tvs, types.Tv{
				ID:   tv.ID,
				Name: tv.Name,
				Url:  tv.Url,
			})
		}
	} else if req.Mode == "re" {
		cTvs, err := utils.Search(req.TvName, "regex")

		if err != nil {
			return nil, err
		}

		for i := 0; i < len(cTvs); i++ {
			tv := cTvs[i]
			tvs = append(tvs, types.Tv{
				Name: tv.Name,
				Url:  tv.Url,
			})
		}
	}

	return tvs, nil
}
