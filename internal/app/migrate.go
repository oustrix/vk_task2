package app

func (a *App) migrate() error {
	return a.db.AutoMigrate()
}
