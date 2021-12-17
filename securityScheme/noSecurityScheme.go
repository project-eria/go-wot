package securityScheme

type NoSecurity struct {
	Security
}

func NewNoSecurity() NoSecurity {
	security := Security{
		Scheme: "nosec",
	}
	return NoSecurity{
		Security: security,
	}
}
