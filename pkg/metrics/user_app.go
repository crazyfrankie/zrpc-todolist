package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	UserRegisterCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "user_register_total",
		Help: "The number of user login",
	})

	UserLoginCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "user_login_total",
		Help: "The number of user login",
	})
)

func RegistryUser() {
	registry.MustRegister(UserRegisterCounter, UserLoginCounter)
}
