package matrix

import (
	"math"
)

type VectorComponent = int
type QuaternionComponent = int

const (
	Vx VectorComponent = iota
	Vy
	Vz
	Vw
)

const (
	Qw QuaternionComponent = iota
	Qx
	Qy
	Qz
)

type tFloatingPoint interface {
	~float32 | ~float64
}

type tSigned interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type tUnsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type tInteger interface {
	tSigned | tUnsigned
}

type tNumber interface {
	tInteger | tFloatingPoint
}

type tVector interface {
	Vec2 | Vec3 | Vec4 | Quaternion
}

type tMatrix interface {
	Mat3 | Mat4
}

func rad2Deg[T tFloatingPoint](radian T) T {
	return radian * (180.0 / math.Pi)
}

func deg2Rad[T tFloatingPoint](degree T) T {
	return degree * (math.Pi / 180.0)
}

func clamp[T tFloatingPoint](current, minimum, maximum T) T {
	return T(max(minimum, min(maximum, current)))
}