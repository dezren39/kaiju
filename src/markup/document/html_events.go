/*****************************************************************************/
/* html_events.go                                                            */
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

package document

import (
	"kaiju/systems/events"
	"kaiju/ui"
)

func tryMap(attr string, elm *DocElement, evt *events.Event, funcMap map[string]func(*DocElement)) {
	if funcName := elm.HTML.Attribute(attr); len(funcName) > 0 {
		if f, ok := funcMap[funcName]; ok {
			evt.Add(func() { f(elm) })
		}
	}
}

func setupEvents(elm *DocElement, funcMap map[string]func(*DocElement)) {
	tryMap("onclick", elm, elm.UI.Event(ui.EventTypeClick), funcMap)
	tryMap("onmouseover", elm, elm.UI.Event(ui.EventTypeEnter), funcMap)
	tryMap("onmouseenter", elm, elm.UI.Event(ui.EventTypeEnter), funcMap)
	tryMap("onmouseleave", elm, elm.UI.Event(ui.EventTypeExit), funcMap)
	tryMap("onmouseexit", elm, elm.UI.Event(ui.EventTypeExit), funcMap)
	tryMap("onmousedown", elm, elm.UI.Event(ui.EventTypeDown), funcMap)
	tryMap("onmouseup", elm, elm.UI.Event(ui.EventTypeUp), funcMap)
	tryMap("onmousewheel", elm, elm.UI.Event(ui.EventTypeScroll), funcMap)
	tryMap("onchange", elm, elm.UI.Event(ui.EventTypeChange), funcMap)
}
