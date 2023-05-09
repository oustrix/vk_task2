package app

import "vk_task2/internal/repository/postgres/service"

func (a *App) migrate() error {
	return a.db.AutoMigrate(
		&service.Service{})
}
