package gpms

type Provider interface {
	Send(subject string, opts ...SendOption) error
	withConfig(c *config) Provider
}

type Option interface {
	apply(*config)
}

type SendOption interface {
	apply(*sendConfig)
}

type optionFunc func(*config)
type sendOptionFunc func(*sendConfig)

func (f optionFunc) apply(c *config) {
	f(c)
}

func (f sendOptionFunc) apply(c *sendConfig) {
	f(c)
}

func newSendConfig(opts ...SendOption) *sendConfig {
	cfg := &sendConfig{}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

func WithAPIKey(key string) Option {
	return optionFunc(func(cfg *config) {
		cfg.APIKey = key
	})
}

func WithStaticFrom(from *Email) Option {
	return optionFunc(func(c *config) {
		c.staticFrom = from
	})
}

func WithStaticTo(to []*Email) Option {
	return optionFunc(func(c *config) {
		c.staticTo = to
	})
}

func WithStaticCC(cc []*Email) Option {
	return optionFunc(func(c *config) {
		c.staticCC = cc
	})
}

func WithStaticBCC(bcc []*Email) Option {
	return optionFunc(func(c *config) {
		c.staticBCC = bcc
	})
}

func WithFrom(from *Email) SendOption {
	return sendOptionFunc(func(c *sendConfig) {
		c.from = from
	})
}

func WithVars(vars map[string]interface{}) SendOption {
	return sendOptionFunc(func(c *sendConfig) {
		c.vars = vars
	})
}

func WithTo(to []*Email) SendOption {
	return sendOptionFunc(func(c *sendConfig) {
		c.to = to
	})
}

func WithCC(cc []*Email) SendOption {
	return sendOptionFunc(func(c *sendConfig) {
		c.cc = cc
	})
}

func WithBCC(bcc []*Email) SendOption {
	return sendOptionFunc(func(c *sendConfig) {
		c.bcc = bcc
	})
}

func WithTemplateID(id string) SendOption {
	return sendOptionFunc(func(c *sendConfig) {
		c.templateID = id
	})
}

type config struct {
	APIKey     string
	staticTo   []*Email
	staticFrom *Email
	staticCC   []*Email
	staticBCC  []*Email
}

type sendConfig struct {
	to         []*Email
	from       *Email
	cc         []*Email
	bcc        []*Email
	templateID string
	vars       map[string]interface{}
}

func Load[T Provider](opts ...Option) Provider {
	var p T
	cfg := &config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return p.withConfig(cfg)
}
