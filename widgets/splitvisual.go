package widgets

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"image"
)

type SplitVisual struct {
}

func (s SplitVisual) Layout(gtx layout.Context, maxHeight unit.Dp, left, right layout.Widget) layout.Dimensions {
	bar := gtx.Dp(unit.Dp(10))
	proportion := float32(0.5)
	leftSize := int(proportion*float32(gtx.Constraints.Max.X) - float32(bar))
	rightOffset := leftSize + bar
	rightSize := gtx.Constraints.Max.X - rightOffset
	barRect := image.Rect(leftSize, 0, rightOffset, gtx.Constraints.Max.X)
	area := clip.Rect(barRect).Push(gtx.Ops)
	area.Pop()

	pt := gtx.Constraints.Max
	if maxHeight != unit.Dp(0) {
		pt.Y = gtx.Dp(maxHeight)
	}
	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftSize, pt.Y))
		left(gtx)
	}

	{
		off := op.Offset(image.Pt(rightOffset, 0)).Push(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightSize, pt.Y))
		right(gtx)
		off.Pop()
	}
	return layout.Dimensions{Size: pt}
}
