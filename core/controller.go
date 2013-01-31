package core

import (
	"math"

	unum "github.com/metaleap/go-util/num"
)

//	Encapsulates a position-and-direction and provides methods
//	manipulating these with respect to each other (e.g. "move forward"
//	some entity that is rotated facing some arbitrary direction).
type Controller struct {
	//	The position being manipulated by this Controller.
	//	When manipulating this manually (outside the TurnXyz() / MoveXyz() methods),
	//	do so in between calling the BeginUpdate() and EndUpdate() methods.
	Pos unum.Vec3

	//	Indicates which axis is consider "upward". This is typically
	//	the Y-axis, denoted by the default value (0, 1, 0).
	//	When manipulating this manually (outside the TurnXyz() / MoveXyz() methods),
	//	do so in between calling the BeginUpdate() and EndUpdate() methods.
	UpAxis unum.Vec3

	Params *ControllerParams

	thrApp struct {
		mat unum.Mat4
	}
	thrPrep struct {
		mat unum.Mat4
	}

	dir            unum.Vec3
	autoUpdate     bool
	hAngle, vAngle float64
}

func (me *Controller) applyTranslation() {
	if me.autoUpdate {
		thrApp.ctlTmps.posNeg.SetFromNeg(&me.Pos)
		thrApp.ctlTmps.matLook.LookAt(&thrApp.numBag, &me.dir, &me.UpAxis)
		thrApp.ctlTmps.matTrans.Translation(&thrApp.ctlTmps.posNeg)
		me.thrApp.mat.SetFromMult4(&thrApp.ctlTmps.matLook, &thrApp.ctlTmps.matTrans)
	}
}

func (me *Controller) applyRotation() {
	if me.autoUpdate {
		thrApp.ctlTmps.axis.Set(0, 1, 0)
		me.dir.Set(1, 0, 0)
		me.dir.RotateDeg(&thrApp.numBag, me.hAngle, &thrApp.ctlTmps.axis)
		me.dir.Normalize()

		thrApp.ctlTmps.axis.SetFromCross(&me.dir)
		thrApp.ctlTmps.axis.Normalize()
		me.dir.RotateDeg(&thrApp.numBag, me.vAngle, &thrApp.ctlTmps.axis)
		me.dir.Normalize()

		me.UpAxis = *me.dir.Cross(&thrApp.ctlTmps.axis)
		me.UpAxis.Normalize()
	}
}

//	Temporarily suspends all matrix re-calculations typically occuring inside
//	the MoveXyz() / TurnXyz() methods. Call this prior to multiple subsequent
//	calls to any combination of those methods, and/or prior to manually modifying
//	the Pos, Dir or UpAxis fields of me. Immediately afterwards, be sure to call
//	EndUpdate() to apply all changes in a final matrix re-calculation.
func (me *Controller) BeginUpdate() {
	me.autoUpdate = false
}

func (me *Controller) CopyFrom(ctl *Controller) {
	thrApp.ctlTmps.tmpCopy = *ctl
	thrApp.ctlTmps.tmpCopy.thrPrep = me.thrPrep
	*me = thrApp.ctlTmps.tmpCopy
}

//	The direction being manipulated by this Controller.
//	CAUTION: this returns a pointer to the direction vector to avoid a copy, but it's
//	NOT meant to be modified, as the vector is re-computed by the TurnFoo() methods.
func (me *Controller) Dir() *unum.Vec3 {
	return &me.dir
}

//	Applies all changes made to Pos, Dir or UpAxis since BeginUpdate() was last
//	called, and recalculates this Controller's final 4x4 transformation matrix.
//	Also resumes all matrix re-calculations typically occuring inside the
//	MoveXyz() / TurnXyz() methods that were suspended since BeginUpdate().
func (me *Controller) EndUpdate() {
	me.autoUpdate = true
	me.applyRotation()
	me.applyTranslation()
}

func (me *Controller) init() {
	me.Params = Core.Options.Misc.DefaultControllerParams
	me.autoUpdate, me.dir.Z, me.UpAxis.Y = true, 1, 1
	unum.Mat4Identities(&me.thrPrep.mat, &me.thrApp.mat)
	htarget := &unum.Vec3{me.dir.X, 0, me.dir.Z}
	htarget.Normalize()
	if htarget.Z >= 0 {
		if htarget.X >= 0 {
			me.hAngle = 360 - unum.RadToDeg(math.Asin(htarget.Z))
		} else {
			me.hAngle = 180 + unum.RadToDeg(math.Asin(htarget.Z))
		}
	} else {
		if htarget.X >= 0 {
			me.hAngle = unum.RadToDeg(math.Asin(-htarget.Z))
		} else {
			me.hAngle = 90 + unum.RadToDeg(math.Asin(-htarget.Z))
		}
	}
	me.vAngle = -unum.RadToDeg(math.Asin(me.dir.Y))
	me.applyRotation()
	me.applyTranslation()
}

//	Recomputes Pos with regards to UpAxis and Dir to effect a "move backward".
func (me *Controller) MoveBackward() {
	me.Pos.SetFromAddMult1(&me.Pos, &me.dir, me.StepSizeMove())
	me.applyTranslation()
}

//	Recomputes Pos with regards to UpAxis to effect a "move downward".
func (me *Controller) MoveDown() {
	me.Pos.SetFromSubMult1(&me.Pos, &me.UpAxis, me.StepSizeMove())
	me.applyTranslation()
}

//	Recomputes Pos with regards to UpAxis and Dir to effect a "move forward".
func (me *Controller) MoveForward() {
	me.Pos.SetFromSubMult1(&me.Pos, &me.dir, me.StepSizeMove())
	me.applyTranslation()
}

//	Recomputes Pos with regards to UpAxis and Dir to effect a "move left-ward".
func (me *Controller) MoveLeft() {
	me.Pos.SetFromAddMult1(&me.Pos, me.dir.CrossNormalized(&me.UpAxis), me.StepSizeMove())
	me.applyTranslation()
}

//	Recomputes Pos with regards to UpAxis and Dir to effect a "move right-ward".
func (me *Controller) MoveRight() {
	me.Pos.SetFromAddMult1(&me.Pos, me.UpAxis.CrossNormalized(&me.dir), me.StepSizeMove())
	me.applyTranslation()
}

//	Recomputes Pos with regards to UpAxis to effect a "move upward".
func (me *Controller) MoveUp() {
	me.Pos.SetFromAddMult1(&me.Pos, &me.UpAxis, me.StepSizeMove())
	me.applyTranslation()
}

func (me *Controller) rotH(deg float64) {
	me.hAngle += deg
	me.applyRotation()
	me.applyTranslation()
}

func (me *Controller) rotV(deg float64) {
	me.vAngle += deg
	me.applyRotation()
	me.applyTranslation()
}

//	Returns the current distance that a single MoveXyz() call (per loop iteration) would move.
//	(Loop.TickDelta * me.Params.MoveSpeed * me.Params.MoveSpeedupFactor)
func (me *Controller) StepSizeMove() float64 {
	return Loop.TickDelta * me.Params.MoveSpeed * me.Params.MoveSpeedupFactor
}

//	Returns the current degrees that a single TurnXyz() call (per loop iteration) would turn.
//	(Loop.TickDelta * me.Params.TurnSpeed * me.Params.TurnSpeedupFactor)
func (me *Controller) StepSizeTurn() float64 {
	return Loop.TickDelta * me.Params.TurnSpeed * me.Params.TurnSpeedupFactor
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn downward" by me.StepSizeTurn() degrees.
func (me *Controller) TurnDown() {
	me.TurnDownBy(me.StepSizeTurn())
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn downward" by the specified degrees.
func (me *Controller) TurnDownBy(deg float64) {
	if me.vAngle > me.Params.MinTurnDown {
		me.rotV(-deg)
	}
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn left-ward" by me.StepSizeTurn() degrees.
func (me *Controller) TurnLeft() {
	me.TurnLeftBy(me.StepSizeTurn())
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn left-ward" by the specified degrees.
func (me *Controller) TurnLeftBy(deg float64) {
	me.rotH(deg)
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn right-ward" by me.StepSizeTurn() degrees.
func (me *Controller) TurnRight() {
	me.TurnRightBy(me.StepSizeTurn())
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn right-ward" by the specified degrees.
func (me *Controller) TurnRightBy(deg float64) {
	me.rotH(-deg)
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn upward" by me.StepSizeTurn() degrees.
func (me *Controller) TurnUp() {
	me.TurnUpBy(me.StepSizeTurn())
}

//	Recomputes Dir with regards to UpAxis and Pos to effect a "turn upward" by the specified degress.
func (me *Controller) TurnUpBy(deg float64) {
	if me.vAngle < me.Params.MaxTurnUp {
		me.rotV(deg)
	}
}

type ControllerParams struct {
	//	Speed of "moving" in the MoveXyz() methods, in units per second.
	//	Defaults to 2.
	MoveSpeed float64

	//	A factor multiplied with MoveSpeed in the MoveXyz() methods. Defaults to 1.
	MoveSpeedupFactor float64

	//	Speed of "turning" in the TurnXyz() methods, in degrees per second.
	//	Defaults to 90.
	TurnSpeed float64

	//	A factor multiplied with TurnSpeed in the TurnXyz() methods. Defaults to 1.
	TurnSpeedupFactor float64

	//	The maximum degree that TurnUp() allows. Defaults to 90.
	MaxTurnUp float64

	//	The minimum degree that TurnDown() allows. Defaults to -90.
	MinTurnDown float64
}

func NewControllerParams() (me *ControllerParams) {
	me = &ControllerParams{}
	me.MoveSpeed, me.MoveSpeedupFactor, me.TurnSpeed, me.TurnSpeedupFactor = 2, 1, 90, 1
	me.MaxTurnUp, me.MinTurnDown = 90, -90
	return
}
