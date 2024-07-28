package server

type Server interface {
	RegisterV1Routes()
	Start()
}
