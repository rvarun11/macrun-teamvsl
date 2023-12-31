package amqp

const (
	// exchangeKind       = "direct"
	// exchangeDurable    = true
	// exchangeAutoDelete = false
	// exchangeInternal   = false
	// exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	// publishMandatory = false
	// publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = true
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)
