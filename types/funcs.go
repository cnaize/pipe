package types

func MergeSendOut(dst, src *SendOut) *SendOut {
	if src.DirMakeAll != nil {
		dst.DirMakeAll = src.DirMakeAll
	}
	if src.FileOpen != nil {
		dst.FileOpen = src.FileOpen
	}
	if src.FileCreate != nil {
		dst.FileCreate = src.FileCreate
	}
	if src.Sha256 != nil {
		dst.Sha256 = src.Sha256
	}
	return dst
}
