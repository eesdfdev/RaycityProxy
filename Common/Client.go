package Common

type IClient interface {
	RecvCallback(payload []byte)
	OnExit()
}
