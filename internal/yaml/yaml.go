package yaml

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

func Pretty(r io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	fmt.Printf("%s", string(out))
	return nil
}
