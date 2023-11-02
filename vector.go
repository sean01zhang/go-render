package main

type vec2 [2]float32
type vec3 [3]float32
type vec4 [4]float32

func (a *vec3) Plus(b vec3) vec3 {
	return vec3{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func (a *vec3) Neg() vec3 {
	return vec3{-a[0], -a[1], -a[2]}
}

func (a *vec3) Minus(b vec3) vec3 {
	return a.Plus(b.Neg())
}

func (a *vec3) Dot(b vec3) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func (a *vec3) Cross(b vec3) vec3 {
	return vec3{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}
