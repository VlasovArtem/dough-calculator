//go:generate mockgen -source=config_manager.go -destination=mocks/config_manager.go -package=mocks

package domain

import "dough-calculator/internal/config"

type ConfigManager interface {
	ParseConfig() error
	GetConfig() config.Config
}
