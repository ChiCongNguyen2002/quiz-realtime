package database

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type DBGroup struct {
	Master   *gorm.DB
	Replicas []*gorm.DB
}

func NewPostgres(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := configurePool(sqlDB, cfg); err != nil {
		return nil, err
	}

	return db, nil
}

func NewDBGroup(masterCfg Config, replicaCfgs []Config) (*DBGroup, error) {
	master, err := NewPostgres(masterCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master: %w", err)
	}

	replicas := make([]*gorm.DB, 0, len(replicaCfgs))
	for i, cfg := range replicaCfgs {
		replica, err := NewPostgres(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to replica %d: %w", i, err)
		}
		replicas = append(replicas, replica)
	}

	return &DBGroup{
		Master:   master,
		Replicas: replicas,
	}, nil
}

func configurePool(db *sql.DB, cfg Config) error {
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	return nil
}

func (g *DBGroup) MasterDB() *gorm.DB {
	return g.Master
}

func (g *DBGroup) ReplicaDB() *gorm.DB {
	if len(g.Replicas) == 0 {
		return g.Master
	}
	return g.Replicas[0]
}

func (g *DBGroup) Close() error {
	sqlDB, err := g.Master.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}

	for _, replica := range g.Replicas {
		sqlDB, err := replica.DB()
		if err != nil {
			continue
		}
		sqlDB.Close()
	}

	return nil
}
