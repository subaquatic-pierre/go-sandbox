package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/widgets"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Exactly 1 Alert argument is required")
		return
	}

	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)

	alertTxt := os.Args[1]

	// create a window
	// with a minimum size of 250*200
	// and sets the title to "Hello Widgets Example"
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(250, 200)
	window.SetWindowTitle("ALERT")

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	// create a line edit
	// with a custom placeholder text
	// and add it to the central widgets layout
	label := widgets.NewQLabel(nil, 0)
	label.SetText(alertTxt)
	widget.Layout().AddWidget(label)

	// create a button
	// connect the clicked signal
	// and add it to the central widgets layout
	// button := widgets.NewQPushButton2("and click me!", nil)
	// button.ConnectClicked(func(bool) {
	// 	widgets.QMessageBox_Information(nil, "OK", input.Text(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	// })
	// widget.Layout().AddWidget(button)

	// make the window visible
	window.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}
