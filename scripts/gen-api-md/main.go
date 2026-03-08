// gen-api-md reads internal/docs/swagger.json and writes API.md at the project root.
// Run via: go run ./scripts/gen-api-md
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type swaggerSpec struct {
	Info struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Host        string                     `json:"host"`
	BasePath    string                     `json:"basePath"`
	Paths       map[string]pathItem        `json:"paths"`
	Definitions map[string]definition      `json:"definitions"`
	Tags        []tagDef                   `json:"tags"`
}

type pathItem map[string]operation

type operation struct {
	Summary    string                `json:"summary"`
	Tags       []string              `json:"tags"`
	Parameters []parameter           `json:"parameters"`
	Security   []map[string][]string `json:"security"`
}

type parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Type        string  `json:"type"`
	Schema      *schema `json:"schema"`
}

type schema struct {
	Ref string `json:"$ref"`
}

type definition struct {
	Required   []string             `json:"required"`
	Properties map[string]property  `json:"properties"`
}

type property struct {
	Type      string   `json:"type"`
	Enum      []string `json:"enum"`
	Items     *schema  `json:"items"`
	MinLength *int     `json:"minLength"`
	MaxLength *int     `json:"maxLength"`
	Minimum   *int     `json:"minimum"`
	Maximum   *int     `json:"maximum"`
}

type tagDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type endpoint struct {
	method  string
	path    string
	summary string
	secured bool
	params  []parameter
}

func main() {
	data, err := os.ReadFile("internal/docs/swagger.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading swagger.json: %v\n", err)
		os.Exit(1)
	}

	var spec swaggerSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing swagger.json: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create("API.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating API.md: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	w := func(format string, args ...interface{}) {
		fmt.Fprintf(f, format, args...)
	}

	// ── Header ──────────────────────────────────────────────────────────────
	w("# Budget API\n\n")
	w("> Auto-generated from swagger spec. Run `make api-docs` to regenerate.\n\n")

	w("## Base URL\n\n")
	w("`http://%s%s`\n\n", spec.Host, spec.BasePath)

	w("## Authentication\n\n")
	w("| Scope | Method | Header |\n")
	w("|-------|--------|--------|\n")
	w("| `/api/v1/*` | Session token | `Authorization: Bearer <token>` |\n")
	w("| `/admin/*` | Static API key | `X-API-Key: <your-api-key>` |\n")
	w("| `POST /api/v1/auth/login` | None (public) | — |\n\n")
	w("Obtain a token via `POST /api/v1/auth/login`. Tokens expire after 30 days. Service accounts use a permanent token seeded from `SERVICE_USER_TOKEN` env var.\n\n")

	w("## Domain Conventions\n\n")
	w("- **Amounts**: string of integer pence. `\"1050\"` = £10.50. Negative = expense, positive = income.\n")
	w("- **Dates**: ISO 8601, `YYYY-MM-DD`\n")
	w("- **IDs**: integer\n\n")

	// ── Group endpoints by tag ───────────────────────────────────────────────
	tagOrder := make([]string, 0, len(spec.Tags))
	tagsSeen := map[string]bool{}
	for _, t := range spec.Tags {
		tagOrder = append(tagOrder, t.Name)
		tagsSeen[t.Name] = true
	}

	byTag := map[string][]endpoint{}

	sortedPaths := make([]string, 0, len(spec.Paths))
	for p := range spec.Paths {
		sortedPaths = append(sortedPaths, p)
	}
	sort.Strings(sortedPaths)

	methodOrder := []string{"get", "post", "put", "patch", "delete"}

	for _, path := range sortedPaths {
		item := spec.Paths[path]
		for _, method := range methodOrder {
			op, ok := item[method]
			if !ok {
				continue
			}
			tag := "other"
			if len(op.Tags) > 0 {
				tag = op.Tags[0]
			}
			if !tagsSeen[tag] {
				tagOrder = append(tagOrder, tag)
				tagsSeen[tag] = true
			}
			byTag[tag] = append(byTag[tag], endpoint{
				method:  strings.ToUpper(method),
				path:    path,
				summary: op.Summary,
				secured: len(op.Security) > 0,
				params:  op.Parameters,
			})
		}
	}

	// ── Endpoints ────────────────────────────────────────────────────────────
	w("## Endpoints\n\n")
	for _, tag := range tagOrder {
		eps, ok := byTag[tag]
		if !ok {
			continue
		}
		w("### %s\n\n", capitalize(tag))
		w("| Method | Path | Auth | Description |\n")
		w("|--------|------|------|-------------|\n")
		for _, e := range eps {
			auth := authLabel(e.path, e.secured)
			w("| `%s` | `%s` | %s | %s |\n", e.method, e.path, auth, e.summary)
		}
		w("\n")

		// Query parameters for any endpoint that has them
		for _, e := range eps {
			var queryParams []parameter
			for _, p := range e.params {
				if p.In == "query" {
					queryParams = append(queryParams, p)
				}
			}
			if len(queryParams) == 0 {
				continue
			}
			w("**`%s %s`** query parameters:\n\n", e.method, e.path)
			w("| Parameter | Type | Required | Description |\n")
			w("|-----------|------|----------|-------------|\n")
			for _, p := range queryParams {
				req := "no"
				if p.Required {
					req = "yes"
				}
				w("| `%s` | %s | %s | %s |\n", p.Name, p.Type, req, p.Description)
			}
			w("\n")
		}
	}

	// ── Request Schemas ──────────────────────────────────────────────────────
	w("## Request Schemas\n\n")

	defNames := make([]string, 0, len(spec.Definitions))
	for name := range spec.Definitions {
		defNames = append(defNames, name)
	}
	sort.Strings(defNames)

	for _, name := range defNames {
		def := spec.Definitions[name]
		shortName := strings.TrimPrefix(name, "model.")

		requiredSet := map[string]bool{}
		for _, r := range def.Required {
			requiredSet[r] = true
		}

		w("### %s\n\n", shortName)
		w("| Field | Type | Required | Notes |\n")
		w("|-------|------|----------|-------|\n")

		fieldNames := make([]string, 0, len(def.Properties))
		for f := range def.Properties {
			fieldNames = append(fieldNames, f)
		}
		sort.Strings(fieldNames)

		for _, field := range fieldNames {
			prop := def.Properties[field]
			req := "no"
			if requiredSet[field] {
				req = "yes"
			}
			typ := prop.Type
			if prop.Items != nil {
				typ = "array[integer]"
			}
			var notes []string
			if len(prop.Enum) > 0 {
				notes = append(notes, "one of: "+strings.Join(prop.Enum, ", "))
			}
			if prop.MinLength != nil && prop.MaxLength != nil {
				notes = append(notes, fmt.Sprintf("len %d–%d", *prop.MinLength, *prop.MaxLength))
			} else if prop.MinLength != nil {
				notes = append(notes, fmt.Sprintf("min len %d", *prop.MinLength))
			} else if prop.MaxLength != nil {
				notes = append(notes, fmt.Sprintf("max len %d", *prop.MaxLength))
			}
			if prop.Minimum != nil && prop.Maximum != nil {
				notes = append(notes, fmt.Sprintf("range %d–%d", *prop.Minimum, *prop.Maximum))
			} else if prop.Minimum != nil {
				notes = append(notes, fmt.Sprintf("min %d", *prop.Minimum))
			} else if prop.Maximum != nil {
				notes = append(notes, fmt.Sprintf("max %d", *prop.Maximum))
			}
			w("| `%s` | %s | %s | %s |\n", field, typ, req, strings.Join(notes, "; "))
		}
		w("\n")
	}

	// ── Example ──────────────────────────────────────────────────────────────
	w("## Example\n\n")
	w("### Login and make a request\n\n")
	w("**1. Login**\n\n")
	w("```http\n")
	w("POST /api/v1/auth/login\n")
	w("Content-Type: application/json\n\n")
	w("{\n")
	w("  \"email\": \"user@example.com\",\n")
	w("  \"password\": \"yourpassword\"\n")
	w("}\n")
	w("```\n\n")
	w("**Response** `200 OK`\n\n")
	w("```json\n")
	w("{\n")
	w("  \"data\": {\n")
	w("    \"token\": \"a3f8c2...\",\n")
	w("    \"expires_at\": \"2026-04-06T10:00:00Z\",\n")
	w("    \"user_id\": 1,\n")
	w("    \"email\": \"user@example.com\"\n")
	w("  },\n")
	w("  \"error\": null\n")
	w("}\n")
	w("```\n\n")
	w("**2. Use the token**\n\n")
	w("```http\n")
	w("GET /api/v1/transactions\n")
	w("Authorization: Bearer a3f8c2...\n")
	w("```\n\n")
	w("---\n\n")
	w("### Create a transaction\n\n")
	w("**Request**\n\n")
	w("```http\n")
	w("POST /api/v1/transactions\n")
	w("Authorization: Bearer a3f8c2...\n")
	w("Content-Type: application/json\n\n")
	w("{\n")
	w("  \"amount\": \"-1050\",\n")
	w("  \"t_date\": \"2026-03-07\",\n")
	w("  \"note\": \"Coffee\",\n")
	w("  \"tag_ids\": [1, 3]\n")
	w("}\n")
	w("```\n\n")
	w("> `\"amount\": \"-1050\"` = £10.50 expense. Positive values are income.\n\n")
	w("**Response** `200 OK`\n\n")
	w("```json\n")
	w("{\n")
	w("  \"id\": 42,\n")
	w("  \"amount\": \"-1050\",\n")
	w("  \"t_date\": \"2026-03-07\",\n")
	w("  \"note\": \"Coffee\",\n")
	w("  \"tag_ids\": [1, 3],\n")
	w("  \"created_at\": \"2026-03-07T10:15:30Z\"\n")
	w("}\n")
	w("```\n")

	fmt.Println("API.md generated successfully")
}

// authLabel returns the auth column value for an endpoint based on its path.
func authLabel(path string, secured bool) string {
	if strings.HasPrefix(path, "/admin/") {
		return "X-API-Key"
	}
	if path == "/auth/login" {
		return "None"
	}
	if secured {
		return "Bearer"
	}
	return ""
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
