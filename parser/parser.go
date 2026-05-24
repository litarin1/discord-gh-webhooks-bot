package parser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/mitchellh/mapstructure"
)

type categoryConfig struct {
	id     string `mapstructure:"id"`
	format string `mapstructure:"format"`
}

type route struct {
	channel  *string         `mapstructure:"channel"`
	category *categoryConfig `mapstructure:"category"`
	webhook  *string         `mapstructure:"webhook"`
}
type Payload struct {
	route    `mapstructure:",squash"`
	children map[string]route `mapstructure:"children"`
}

var schemaRoutes *jsonschema.Resolved = nil

func prepareSchema() error {
	if schemaRoutes == nil {
		bytes, err := os.ReadFile("schema.json")
		if err != nil {
			return err
		}

		var schema jsonschema.Schema
		if err := json.Unmarshal(bytes, &schema); err != nil {
			return err
		}

		schema.Schema = ""
		resolvedSchema, err := schema.Resolve(nil)
		if err != nil {
			return err
		}
		schemaRoutes = resolvedSchema
	}
	return nil
}
func validate(data map[string]any) error {
	if err := prepareSchema(); err != nil {
		return err
	}
	return schemaRoutes.Validate(data)
}

func Parse(jsonData []byte) (Payload, error) {
	var payload Payload

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &payload,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return Payload{}, err
	}

	var m map[string]any
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return Payload{}, err
	}
	if err := validate(m); err != nil {
		return Payload{}, err
	}
	if err := decoder.Decode(m); err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func (p *Payload) GetRoute(author, repo_name string) route {
	if author == "" {
		slog.Error("GetRoute(): author is empty!")
		return route{}
	}
	if repo_name != "" {
		if repo_name != "" {
			if r, found := p.children[author+"/"+repo_name]; found {
				return r
			}
		}
	}
	if r, found := p.children[author]; found {
		return r
	}
	return p.route
}

func (r route) String() string {
	var sb strings.Builder
	if r.channel != nil {
		fmt.Fprintf(&sb, " Channel %v;", *r.channel)
	}
	if r.category != nil {
		fmt.Fprintf(&sb, " Category %v;", *r.category)
	}
	if r.webhook != nil {
		fmt.Fprintf(&sb, " Webhook %v;", *r.webhook)
	}
	return sb.String()
}

func (p Payload) String() string {
	return fmt.Sprint("Default route:", p.route, "Children:", p.children)
}
