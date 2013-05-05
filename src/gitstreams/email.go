package main

type Email struct {
	from      string
	to        []string
	subject   string
	text      string
	variables map[string]string
}

func (e *Email) From() string                 { return e.from }
func (e *Email) To() []string                 { return e.to }
func (e *Email) Subject() string              { return e.subject }
func (e *Email) Text() string                 { return e.text }
func (e *Email) Variables() map[string]string { return e.variables }
func (e *Email) Cc() []string                 { return nil }
func (e *Email) Bcc() []string                { return nil }
func (e *Email) Html() string                 { return "" }
func (e *Email) Headers() map[string]string   { return nil }
func (e *Email) Options() map[string]string   { return nil }
