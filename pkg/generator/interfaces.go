package generator

type EncodeFnRenderer interface {
	RenderEncodeFn(funcPrefix, receivedType string, f Field) string
}
