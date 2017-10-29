package config

type (
	Specification interface {
		Username() string
		Password() string
		Workers() int
	}

	specs struct {
		User       string `envconfig:"username" required:"true"`
		Pass       string `envconfig:"password" required:"true"`
		WorkersNum int    `envconfig:"workers" default:"4"`
	}
)

func NewSpecification() Specification {
	return &specs{}
}

func (s *specs) Username() string {
	return s.User
}

func (s *specs) Password() string {
	return s.Pass
}

func (s *specs) Workers() int {
	return s.WorkersNum
}
