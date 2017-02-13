package unitycompat

import (
	"math"

	"github.com/metaleap/go-util-misc"
	"github.com/metaleap/go-util-num"
)

func Mathf_Abs(value int) int {
	switch {
	case value < 0:
		return -value
	case value == 0:
		return 0
	}
	return value
}

func Mathf_CeilToInt(f float64) int {
	return int(math.Ceil(f))
}

func Mathf_Clamp(val, min, max int) int {
	switch {
	case val < min:
		return min
	case val > max:
		return max
	}
	return val
}

//	Calculates the shortest difference between two given angles.
func Mathf_DeltaAngle(cur, target float64) float64 {
	n := Mathf_Repeat(target-cur, 360)
	if n > 180 {
		n -= 360
	}
	return n
}

//	Returns the largest integer smaller-than or equal-to `f`.
func Mathf_FloorToInt(f float64) int {
	return int(math.Floor(f))
}

func Mathf_Gamma(val, absMax, gamma float64) (r float64) {
	if r = math.Abs(val); r <= absMax {
		r = absMax * math.Pow(r/absMax, gamma)
	}
	if val < 0 {
		r = -r
	}
	return
}

func Mathf_InverseLerp(from, to, val float64) float64 {
	if from < to {
		switch {
		case val < from:
			return 0
		case val > to:
			return 1
		}
		return -from / (to - from)
	}
	switch {
	case from <= to:
		return 0
	case val < to:
		return 1
	case val > from:
		return 0
	}
	return 1 - (val-to)/(from-to)
}

func Mathf_Lerp(from, to, t float64) float64 {
	return unum.Clamp01(t)*(to-from) + from
}

func Mathf_LerpAngle(a, b, t float64) (r float64) {
	if r = Mathf_Repeat(b-a, 360); r > 180 {
		r -= 360
	}
	r = unum.Clamp01(t)*r + a
	return
}

func Mathf_Maxd(values ...float64) (r float64) {
	r = values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > r {
			r = values[i]
		}
	}
	return
}

func Mathf_Maxi(a, b int) int {
	return ugo.Ifi(b > a, b, a)
}

func Mathf_Maxn(values ...int) (r int) {
	r = values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > r {
			r = values[i]
		}
	}
	return
}

func Mathf_Mind(values ...float64) (r float64) {
	r = values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < r {
			r = values[i]
		}
	}
	return
}

func Mathf_Mini(a, b int) int {
	return ugo.Ifi(b < a, b, a)
}

func Mathf_Minn(values ...int) (r int) {
	r = values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < r {
			r = values[i]
		}
	}
	return
}

func Mathf_MoveTowards(cur, target, maxDelta float64) float64 {
	tc := target - cur
	if math.Abs(tc) <= maxDelta {
		return target
	}
	return maxDelta*Mathf_Sign(tc) + cur
}

func Mathf_MoveTowardsAngle(cur, target, maxDelta float64) float64 {
	return Mathf_MoveTowards(cur, cur+Mathf_DeltaAngle(cur, target), maxDelta)
}

func Mathf_PingPong(t, l float64) float64 {
	t = Mathf_Repeat(t, l*2)
	return l - math.Abs(t-l)
}

//	Like `math.Mod(t, l)` but for float64; 'l' must not be 0.
func Mathf_Repeat(t, l float64) float64 {
	return t - math.Floor(t/l)*l
}

//	Returns the next-higher integer if fraction>0.5; if fraction<0.5 returns the next-lower integer; if fraction==0.5, returns the next even integer.
func Mathf_RoundToInt(v float64) int {
	return int(unum.Round(v))
}

func Mathf_Sign(f float64) float64 {
	if f < 0 {
		return -1
	}
	return 1
}

//	Gradually changes a value towards a desired goal over time.
func Mathf_SmoothDamp(current, target float64, currentVelocity *float64, smoothTime, maxSpeed, deltaTime float64) float64 {
	smoothTime = math.Max(0.0001, smoothTime)
	num := 2 / smoothTime
	num2 := num * deltaTime
	num3 := 1 / (1 + num2 + 0.48*num2*num2 + 0.235*num2*num2*num2)
	num4 := current - target
	num5 := target
	num6 := maxSpeed * smoothTime
	num4 = unum.Clamp(num4, -num6, num6)
	target = current - num4
	num7 := (*currentVelocity + num*num4) * deltaTime
	*currentVelocity = (*currentVelocity - num*num7) * num3
	num8 := target + (num4+num7)*num3
	if (num5-current > 0) == (num8 > num5) {
		num8 = num5
		*currentVelocity = (num8 - num5) / deltaTime
	}
	return num8
}

func Mathf_SmoothDampAngle(current, target float64, currentVelocity *float64, smoothTime, maxSpeed, deltaTime float64) float64 {
	return Mathf_SmoothDamp(current, current+Mathf_DeltaAngle(current, target), currentVelocity, smoothTime, maxSpeed, deltaTime)
}

func Mathf_SmoothStep(from, to, t float64) float64 {
	t = unum.Clamp01(t)
	t = -2*t*t*t + 3*t*t
	return to*t + from*(1-t)
}
