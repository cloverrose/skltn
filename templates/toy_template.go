package templates

func init() {
	templateMap["toy"] = toyTemplate
}

const toyTemplate = `Hello {{.ReceiverTyp}}`
