/*
 * Copyright (c) 2026 Gábor Nyíri as known as Halacs
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"fmt"
	"go-qt/files"
	"go-qt/layouts"
	"os"

	qt "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
)

func main() {
	qt.NewQApplication(os.Args)

	// Create the UI defined in layouts/first.ui
	ui := layouts.NewFrameUi()

	fmt.Println(ui.Frame.ObjectName())
	ui.Frame.SetWindowTitle("QT version: " + qt.QLibraryInfo_Version().ToString())

	m := NewModel(func(progress int) {
		mainthread.Start(func() {
			ui.ProgressBar.SetValue(progress)
			if progress >= 100 {
				ui.Label.SetText("Finished")
			}
		})
	})

	// Hook up signals to the widgets from first.ui
	ui.PushButton_start.OnClicked(func() {
		ui.Label.SetText("Started")
		m.Start()
	})

	ui.PushButton_stop.OnClicked(func() {
		ui.Label.SetText("Stopped")
		m.Stop()
	})

	// Load image from embedded files.Folder into GraphicsView
	imgData, err := files.Folder.ReadFile("Qt-logo-neon_900px.png")
	if err != nil {
		panic(err)
	}

	pixmap := qt.NewQPixmap()
	pixmap.LoadFromDataWithData(imgData)

	scene := qt.NewQGraphicsScene4(ui.GraphicsView.QObject)
	item := scene.AddPixmap(pixmap)
	ui.GraphicsView.SetScene(scene)

	ui.GraphicsView.OnResizeEvent(func(super func(event *qt.QResizeEvent), event *qt.QResizeEvent) {
		super(event)
		ui.GraphicsView.FitInView5(item.QGraphicsItem, qt.KeepAspectRatio)
	})

	// Show the window and start the Qt event loop
	ui.Frame.Show()

	// Scale the image to fit inside the GraphicsView initially
	ui.GraphicsView.FitInView5(item.QGraphicsItem, qt.KeepAspectRatio)

	qt.QApplication_Exec()
}
