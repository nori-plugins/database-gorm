package plugin

import (
	"context"
	"errors"

	"github.com/nori-io/common/v4/pkg/domain/registry"

	"github.com/nori-plugins/database-orm-gorm/internal/hook"

	"github.com/jinzhu/gorm"
	"github.com/nori-io/common/v4/pkg/domain/config"
	em "github.com/nori-io/common/v4/pkg/domain/enum/meta"
	"github.com/nori-io/common/v4/pkg/domain/logger"
	"github.com/nori-io/common/v4/pkg/domain/meta"
	p "github.com/nori-io/common/v4/pkg/domain/plugin"
	m "github.com/nori-io/common/v4/pkg/meta"
	i "github.com/nori-io/interfaces/database/orm/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	Plugin   p.Plugin = plugin{}
	dialects          = [4]string{"mssql", "mysql", "postgres", "sqlite"}
)

type plugin struct {
	db     *gorm.DB
	config conf
	logger logger.FieldLogger
}

type conf struct {
	dsn     string
	dialect string
	logMode bool
}

func (p plugin) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
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

func (p plugin) Instance() interface{} {
	return p.db
}

func (p plugin) Meta() meta.Meta {
	return m.Meta{
		ID: m.ID{
			ID:      "sql/gorm",
			Version: "1.9.15",
		},
		Author: m.Author{
			Name: "Nori.io",
			URL:  "https://nori.io/",
		},
		Dependencies: []meta.Dependency{},
		Description: m.Description{
			Title:       "",
			Description: "This plugin implements instance of ORM GORM",
		},
		Interface: i.GormInterface,
		License:   []meta.License{},
		Links:     []meta.Link{},
		Repository: m.Repository{
			Type: em.Git,
			URL:  "https://github.com/nori-io/sql-gorm",
		},
		Tags: []string{"orm", "gorm", "sql", "database", "db"},
	}
}

func (p plugin) Start(ctx context.Context, registry registry.Registry) error {
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

func (p plugin) Stop(ctx context.Context, registry registry.Registry) error {
	err := p.db.Close()
	if err != nil {
		p.logger.Error(err.Error())
	}

	return err
}
