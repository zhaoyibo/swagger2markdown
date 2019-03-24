package tool

import "bytes"

type StringBuilder struct {
	buf bytes.Buffer
}

func (sb *StringBuilder) Append(str string) *StringBuilder {
	sb.buf.WriteString(str)
	return sb
}

func (sb *StringBuilder) Br2() *StringBuilder {
	sb.buf.WriteString("\n")
	sb.buf.WriteString("\n")
	return sb
}

func (sb *StringBuilder) Br() *StringBuilder {
	sb.buf.WriteString("\n")
	return sb
}

func (sb *StringBuilder) String() string {
	return sb.buf.String()
}
