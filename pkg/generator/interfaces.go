package generator

type EncodeFnGenerator interface {
	GenerateEncodeFn(funcPrefix, receivedType string, f Field) string
}
