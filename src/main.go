/*****************************************************************************/
/* main.go                                                                   */
/*****************************************************************************/
/*                           This file is part of:                           */
/*                                KAIJU ENGINE                               */
/*                          https://kaijuengine.org                          */
/*****************************************************************************/
/* MIT License                                                               */
/*                                                                           */
/* Copyright (c) 2023-present Kaiju Engine contributors (CONTRIBUTORS.md).   */
/* Copyright (c) 2015-2023 Brent Farris.                                     */
/*                                                                           */
/* May all those that this source may reach be blessed by the LORD and find  */
/* peace and joy in life.                                                    */
/* Everyone who drinks of this water will be thirsty again; but whoever      */
/* drinks of the water that I will give him shall never thirst; John 4:13-14 */
/*                                                                           */
/* Permission is hereby granted, free of charge, to any person obtaining a   */
/* copy of this software and associated documentation files (the "Software"),*/
/* to deal in the Software without restriction, including without limitation */
/* the rights to use, copy, modify, merge, publish, distribute, sublicense,  */
/* and/or sell copies of the Software, and to permit persons to whom the     */
/* Software is furnished to do so, subject to the following conditions:      */
/*                                                                           */
/* The above copyright, blessing, biblical verse, notice and                 */
/* this permission notice shall be included in all copies or                 */
/* substantial portions of the Software.                                     */
/*                                                                           */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS   */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.    */
/* IN NO EVENT SHALL THE /* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY   */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE     */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                             */
/*****************************************************************************/

package main

import (
	"fmt"
	"kaiju/bootstrap"
	"kaiju/editor/ui/hierarchy"
	"kaiju/engine"
	"kaiju/host_container"
	"kaiju/matrix"
	"kaiju/profiler"
	"kaiju/systems/console"
	tests "kaiju/tests/rendering_tests"
	"kaiju/tools/html_preview"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func addConsole(host *engine.Host) {
	console.For(host).AddCommand("EntityCount", func(*engine.Host, string) string {
		return fmt.Sprintf("Entity count: %d", len(host.Entities()))
	})
	html_preview.SetupConsole(host)
	hierarchy.SetupConsole(host)
	profiler.SetupConsole(host)
	tests.SetupConsole(host)
}

func main() {
	container := host_container.New("Kaiju")
	go container.Run(engine.DefaultWindowWidth, engine.DefaultWindowHeight)
	<-container.PrepLock
	container.RunFunction(func() {
		container.Host.Camera.SetPosition(matrix.Vec3{0.0, 0.0, 2.0})
		addConsole(container.Host)
	})
	bootstrap.Main(container)
}
