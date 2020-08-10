package main

import (
	"context"
	"errors"
	"github.com/nori-io/sql-gorm/internal/hook"

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

type pluginConfig struct {
	dsn     string
	dialect string
	logMode bool
}

var (
	Plugin   plugin.Plugin = &service{}
	dialects               = [4]string{"mssql", "mysql", "postgres", "sqlite"}
)

func (p *service) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {

	var isValidDialect bool

	p.logger = log
	p.config.logMode = config.Bool("sql.gorm.logMode", "log mode: true or false")()
	p.config.dsn = config.String("sql.gorm.dsn", "database connection string")()
	p.config.dialect = config.String("sql.gorm.dialect", "sql dialect: mssql, mysql, postgres, sqlite")()

	for _, v := range dialects {
		if v == p.config.dialect {
			isValidDialect = true
		}
	}

	if !isValidDialect {
		return errors.New("Dialect is wrong. You should use on of sql dialects: mssql, mysql, postgres, sqlite")
	}
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
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name:        "Nori: ORM GORM",
			Description: "This plugin implements instance of ORM GORM",
		},
		Core: meta.Core{
			VersionConstraint: "^0.2.0",
		},
		Interface: i.GormInterface,
		License: []meta.License{
			{
				Title: "GPLv3",
				Type:  "GPLv3",
				URI:   "https://www.gnu.org/licenses/"},
		},
		Links:      nil,
		Repository: meta.Repository{
			Type: "git",
			URI:  "https://github.com/nori-io/sql-gorm",
		},
		Tags:       []string{"orm", "gorm", "sql", "database", "db"},
	}
}

func (p *service) Start(ctx context.Context, registry plugin.Registry) error {
	var err error
	p.db, err = gorm.Open(p.config.dialect, p.config.dsn)
	if err != nil {
		p.logger.Error(err.Error())
	} else {
		p.db.LogMode(p.config.logMode)
		if p.config.logMode == true {
			p.db.SetLogger(&hook.Logger{Origin: p.logger})
		}
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
