package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"io"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type (
	IMiniApp interface {
		ListAppsCategories() ([]*entity.MiniAppCategory, error)

		ListCollectedApps(userId int64) ([]*entity.MiniAppListDto, error)

		IsAppsCollected(appIds []string, userId int64) (map[string]bool, error)

		DoCollectApp(appId string, userId int64) error

		DoUnCollectApp(appId string, userId int64) error

		// ListApps category=0 则不筛选
		ListApps(category int64) ([]*entity.MiniAppListDto, error)

		ListAppsByUuids(uuids []string) ([]*entity.MiniAppListDto, error)

		ListAppsByUser(userId, category int64) ([]*entity.MiniAppListDto, error)

		GetAppById(id int64) (*entity.MiniApp, bool, error)

		GetAppOutputById(outputId string) (*entity.MiniAppOutput, bool, error)

		CheckOutputExists(outputId string) (bool, error)

		CheckAppRanByUser(appId string, userId int64) (bool, error)

		CheckAppExists(uuid string) (bool, error)

		GetAppDetailByUuid(uuid string) (*entity.MiniAppDetailDto, bool, error)

		GetAppByUuid(uuid string) (*entity.MiniApp, bool, error)

		SaveApp(app *entity.MiniApp) (*entity.MiniApp, error)

		CreateOutput(output *entity.MiniAppOutput) error

		RunApp(userId int64, uuid string, param entity.MiniAppRunParam) (*ChannelStreamTrunkReader[*entity.MiniAppOutputStreamChunk], error)

		PageOpenedAppOutputsByAppId(query page.StreamQuery, uuid string) (*page.StreamResult[*entity.MiniAppOutputDto], int64, error)

		PageUsersApps(query page.StreamQuery, userId int64) (*page.StreamResult[*entity.MiniAppListDto], error)

		ListOpenedAppOutputsByAppId(uuid string) ([]*entity.MiniAppOutputDto, error)

		ListAIModels() ([]*entity.MiniAppAiModel, error)

		ListAIModelCategories() ([]*entity.MiniAppAiModelCategory, error)

		DeleteApp(userId int64, appId string) error

		CountUserCreatedApps(userId int64) (int64, error)

		CountUserRanApps(userId int64) (int64, error)

		CountUsersAppsRuntimes(userId int64) (int64, error)
	}
)

type ChannelStreamTrunkReader[T comparable] struct {
	ch  chan T
	err chan error
}

func NewChannelStreamTrunkReader[T comparable](ch chan T) *ChannelStreamTrunkReader[T] {
	return &ChannelStreamTrunkReader[T]{
		ch:  ch,
		err: make(chan error, 1),
	}
}

func (g ChannelStreamTrunkReader[T]) SetInterrupt(err error) {
	g.err <- err
}

func (g ChannelStreamTrunkReader[T]) Close() error {
	return nil
}

func (g ChannelStreamTrunkReader[T]) Read() (chunk T, err error) {
	select {
	case err = <-g.err:
		return

	case c := <-g.ch:
		if c == chunk {
			err = io.EOF
			return
		}

		chunk = c
		return
	}
}
