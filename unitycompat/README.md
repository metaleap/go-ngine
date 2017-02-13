# unitycompat
--
    import "github.com/metaleap/go-ngine/unitycompat"


## Usage

#### func  Mathf_Abs

```go
func Mathf_Abs(value int) int
```

#### func  Mathf_CeilToInt

```go
func Mathf_CeilToInt(f float64) int
```

#### func  Mathf_Clamp

```go
func Mathf_Clamp(val, min, max int) int
```

#### func  Mathf_DeltaAngle

```go
func Mathf_DeltaAngle(cur, target float64) float64
```
Calculates the shortest difference between two given angles.

#### func  Mathf_FloorToInt

```go
func Mathf_FloorToInt(f float64) int
```
Returns the largest integer smaller-than or equal-to `f`.

#### func  Mathf_Gamma

```go
func Mathf_Gamma(val, absMax, gamma float64) (r float64)
```

#### func  Mathf_InverseLerp

```go
func Mathf_InverseLerp(from, to, val float64) float64
```

#### func  Mathf_Lerp

```go
func Mathf_Lerp(from, to, t float64) float64
```

#### func  Mathf_LerpAngle

```go
func Mathf_LerpAngle(a, b, t float64) (r float64)
```

#### func  Mathf_Maxd

```go
func Mathf_Maxd(values ...float64) (r float64)
```

#### func  Mathf_Maxi

```go
func Mathf_Maxi(a, b int) int
```

#### func  Mathf_Maxn

```go
func Mathf_Maxn(values ...int) (r int)
```

#### func  Mathf_Mind

```go
func Mathf_Mind(values ...float64) (r float64)
```

#### func  Mathf_Mini

```go
func Mathf_Mini(a, b int) int
```

#### func  Mathf_Minn

```go
func Mathf_Minn(values ...int) (r int)
```

#### func  Mathf_MoveTowards

```go
func Mathf_MoveTowards(cur, target, maxDelta float64) float64
```

#### func  Mathf_MoveTowardsAngle

```go
func Mathf_MoveTowardsAngle(cur, target, maxDelta float64) float64
```

#### func  Mathf_PingPong

```go
func Mathf_PingPong(t, l float64) float64
```

#### func  Mathf_Repeat

```go
func Mathf_Repeat(t, l float64) float64
```
Like `math.Mod(t, l)` but for float64; 'l' must not be 0.

#### func  Mathf_RoundToInt

```go
func Mathf_RoundToInt(v float64) int
```
Returns the next-higher integer if fraction>0.5; if fraction<0.5 returns the
next-lower integer; if fraction==0.5, returns the next even integer.

#### func  Mathf_Sign

```go
func Mathf_Sign(f float64) float64
```

#### func  Mathf_SmoothDamp

```go
func Mathf_SmoothDamp(current, target float64, currentVelocity *float64, smoothTime, maxSpeed, deltaTime float64) float64
```
Gradually changes a value towards a desired goal over time.

#### func  Mathf_SmoothDampAngle

```go
func Mathf_SmoothDampAngle(current, target float64, currentVelocity *float64, smoothTime, maxSpeed, deltaTime float64) float64
```

#### func  Mathf_SmoothStep

```go
func Mathf_SmoothStep(from, to, t float64) float64
```

--
**godocdown** http://github.com/robertkrimen/godocdown
