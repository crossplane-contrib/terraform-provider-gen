package generator

type EncodeFnGenerator interface {
	GenerateEncodeFn(funcPrefix, receivedType string, f Field) string
}

type DecodeFnGenerator interface {
	GenerateDecodeFn(funcPrefix, receivedType string, f Field) string
}

type MergeFnGenerator interface {
	GenerateMergeFn(funcPrefix, receivedType string, f Field, isSpec bool) string
}
