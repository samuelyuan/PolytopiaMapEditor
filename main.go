package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/samuelyuan/PolytopiaMapEditor/internal/ui"
)

func main() {
	e := ui.NewEditor()

	a := app.NewWithID("io.fyne.polytopiamapeditor")
	w := a.NewWindow("Polytopia Map Editor")
	e.BuildUI(w)
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}
