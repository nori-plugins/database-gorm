package main

import (
	"context"
	"errors"
	"log"
	"strconv"

	gormLogger "gorm.io/gorm/logger"

	p "github.com/nori-io/common/v5/pkg/domain/plugin"

	"github.com/nori-io/common/v5/pkg/domain/registry"

	"github.com/nori-plugins/database-gorm/internal/hook"

	"github.com/nori-io/common/v5/pkg/domain/config"
	em "github.com/nori-io/common/v5/pkg/domain/enum/meta"
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	m "github.com/nori-io/common/v5/pkg/meta"
	i "github.com/nori-io/interfaces/database/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dialects = [3]string{"mysql", "postgres", "sqlite"}
	logModes = [4]string{"silent", "error", "warn", "info"}
)

func New() p.Plugin {
	log.Println("NEW!!!")
	return &plugin{}
}

type plugin struct {
	db     *gorm.DB
	config conf
	logger logger.FieldLogger
}

type conf struct {
	dsn     string
	dialect string
	logMode string
}

func (p *plugin) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	p.logger = log

	var isValidDialect, isValidLogMode bool

	p.config.logMode = config.String("logMode", "log mode: silent, error, warn, info")()
	p.config.dsn = config.String("dsn", "database connection string")()
	p.config.dialect = config.String("dialect", "sql dialect: mysql, postgres, sqlite")()
	for _, v := range dialects {
		if v == p.config.dialect {
			isValidDialect = true
		}
	}

	for i, v := range logModes {
		if v == p.config.dialect {
			isValidLogMode = true
			p.config.logMode = string(i)
		}
	}

	if !isValidDialect {
		return errors.New("Dialect is wrong. You should use on of sql dialects: mysql, postgres, sqlite")
	}

	if !isValidLogMode {
		p.config.logMode = string(0)
	}

	return nil
}

func (p *plugin) Instance() interface{} {
	return p.db
}

func (p *plugin) Meta() meta.Meta {
	return m.Meta{
		ID: m.ID{
			ID:      "nori/sql/gorm",
			Version: "1.21.11",
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

func (p *plugin) Start(ctx context.Context, registry registry.Registry) error {
	var err error

	switch p.config.dialect {
	case "mysql":
		p.db, err = gorm.Open(mysql.Open(p.config.dsn), &gorm.Config{Logger: &hook.Logger{Origin: p.logger}})
	case "postgres":
		p.db, err = gorm.Open(postgres.Open(p.config.dsn), &gorm.Config{})
		p.logger.Error(err.Error())

	case "sqllite":
		p.db, err = gorm.Open(sqlite.Open(p.config.dsn), &gorm.Config{})
	}
	if err != nil {
		return err
	}

	logLevel, err := strconv.Atoi(p.config.logMode)
	if err != nil {
		return err
	}
	p.db.Logger.LogMode(gormLogger.LogLevel(logLevel))

	return err
}

func (p *plugin) Stop(ctx context.Context, registry registry.Registry) error {
	p.logger.Info("%s", "STOP BEGINNING")
	db, err := p.db.DB()
	if err != nil {
		p.logger.Error(err.Error())
		return err
	}
	err = db.Close()
	p.logger.Info("%s", "STOP ENDING")
	return err
}
