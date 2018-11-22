package consumer

import (
	"time"

	ui "github.com/gizak/termui"
)

func Run() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// build
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, list("Top elements", aggrClicks.TopList)),
			ui.NewCol(4, 0, list("Top prod name ($)", aggrBuyAmountName.TopList)),
			ui.NewCol(4, 0, list("Top prod color ($)", aggrBuyColor.TopList)),
		),
		ui.NewRow(
			ui.NewCol(4, 0, list("Modules #errors", aggrErrorModules.TopList)),
			ui.NewCol(4, 0, list("Top Cart name (view)", aggrAddToCartName.TopList)),
			ui.NewCol(4, 0, list("Top Cart color (view)", aggrAddToCartColor.TopList)),
		),
		ui.NewRow(
			ui.NewCol(8, 0, list("Last errs", latestErrors.Get)),
			ui.NewCol(4, 0, list("Top pages", aggrPageViews.TopList)),
		),
	)

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	//q exits
	exit := func(ui.Event) {
		ui.StopLoop()
	}
	//ui.Handle("q", exit)
	ui.Handle("/sys/kbd", exit)
	//Ctrl+C exits
	//ui.Handle("/sys/kbd/C-c", exit)

	ui.Loop()
}

func list(name string, content func() []string) *ui.List {
	ls := ui.NewList()
	ls.Items = content()
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = name
	ls.Height = 7
	ls.Width = 25
	ls.Y = 0

	go func() {
		update := time.NewTicker(time.Millisecond * 600)
		for {
			<-update.C
			ls.Items = content()
			ui.Render(ls)
		}
	}()
	//ui.Handle("/timer/1s", func(e ui.Event) {
	//	ls.Items = strs()
	//})

	return ls
}
