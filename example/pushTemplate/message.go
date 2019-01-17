package pushTemplate

type Message interface {
	OpenId() string
	Mobile() string
}
