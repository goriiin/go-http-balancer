package global

func (a *App) Run() {

	a.router.HandleFunc("/data", a.dataDelivery.Update)

}
