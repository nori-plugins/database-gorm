package main

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nori-io/common/v3/config"
	"github.com/nori-io/common/v3/logger"
	"github.com/nori-io/common/v3/meta"
	"github.com/nori-io/common/v3/plugin"
	i "github.com/nori-io/interfaces/public/sql/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type service struct {
	db     *gorm.DB
	config *pluginConfig
	logger logger.FieldLogger
}

//one parameter with format  host:port
type pluginConfig struct {
	connectionString string
	dialect          string
}

var (
	Plugin plugin.Plugin = &service{}
)

func (p *service) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	p.config.connectionString = config.String("connectionString", "Connection string consists needed data for connect to database")()
	p.config.dialect = config.String("dialect", "dialect")()
	return nil
}

func (p *service) Instance() interface{} {
	return p.db
}

func (p *service) Meta() meta.Meta {
	return &meta.Data{
		ID: meta.ID{
			ID:      "sql/gorm",
			Version: "1.9.15",
		},
		Author: meta.Author{
			Name: "Nori.io",
			URI:  "https://nori.io/",
		},
		Core: meta.Core{
			VersionConstraint: "=0.2.0",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name:        "Nori: ORM GORM",
			Description: "This plugin implements instance of ORM GORM",
		},
		Interface: i.GormInterface,
		License: []meta.License{
			{
				Title: "GPLv3",
				Type:  "GPLv3",
				URI:   "https://www.gnu.org/licenses/"},
		},
		Tags: []string{"orm", "gorm", "sql", "database", "db"},
	}
}

func (p *service) Start(ctx context.Context, registry plugin.Registry) error {

	_, err := gorm.Open(p.config.dialect, p.config.connectionString)
	if err != nil {
		p.logger.Error(err.Error())
	}

	return err
}

func (p *service) Stop(ctx context.Context, registry plugin.Registry) error {
	err := p.db.Close()
	if err != nil {
		p.logger.Error(err.Error())
	}

	return err
}
