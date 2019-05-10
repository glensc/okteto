package model

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type serialMount struct {
	SendOnly bool   `json:"sendonly,omitempty" yaml:"sendonly,omitempty"`
	Source   string `json:"source,omitempty" yaml:"source,omitempty"`
	Path     string `json:"path" yaml:"path,omitempty"`
	Target   string `json:"target,omitempty" yaml:"target,omitempty"` //TODO: decrecated
	Size     string `json:"size,omitempty" yaml:"size,omitempty"`
}

// UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (e *EnvVar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	parts := strings.SplitN(raw, "=", 2)
	e.Name = parts[0]
	if len(parts) == 2 {
		if strings.HasPrefix(parts[1], "$") {
			e.Value = os.ExpandEnv(parts[1])
			return nil
		}

		e.Value = parts[1]
		return nil
	}

	val := os.ExpandEnv(parts[0])
	if val != parts[0] {
		e.Value = val
	}

	return nil
}

// MarshalYAML Implements the marshaler interface of the yaml pkg.
func (e *EnvVar) MarshalYAML() (interface{}, error) {
	return e.Name + "=" + e.Value, nil
}

// UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (f *Forward) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	parts := strings.SplitN(raw, ":", 2)

	localPort, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("Cannot convert remote port '%s' in port-forward '%s'", parts[0], raw)
	}

	if len(parts) != 2 {
		parts = append(parts, parts[0])
	}

	remotePort, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("Cannot convert remote port '%s' in port-forward '%s'", parts[1], raw)
	}
	f.Local = localPort
	f.Remote = remotePort
	return nil
}

// MarshalYAML Implements the marshaler interface of the yaml pkg.
func (f Forward) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%d:%d", f.Local, f.Remote), nil
}
