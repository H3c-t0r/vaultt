package identity

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
)

func TestPopulate_Basic(t *testing.T) {
	var tests = []struct {
		mode              int
		name              string
		input             string
		output            string
		err               error
		entityName        string
		metadata          map[string]string
		aliasAccessor     string
		aliasID           string
		aliasName         string
		nilEntity         bool
		validityCheckOnly bool
		aliasMetadata     map[string]string
		groupName         string
		groupMetadata     map[string]string
		groupMemberships  []string
	}{
		// time.* tests. Keep tests with time.Now() at the front to avoid false
		// positives due to the second changing during the test
		{
			name:   "time now",
			input:  "{{time.now}}",
			output: strconv.Itoa(int(time.Now().Unix())),
		},
		{
			name:   "time plus",
			input:  "{{time.now.plus.1h}}",
			output: strconv.Itoa(int(time.Now().Unix() + (60 * 60))),
		},
		{
			name:   "time plus",
			input:  "{{time.now.minus.5m}}",
			output: strconv.Itoa(int(time.Now().Unix() - (5 * 60))),
		},
		{
			name:  "invalid operator",
			input: "{{time.now.divide.5m}}",
			err:   errors.New("invalid time operator \"divide\""),
		},
		{
			name:  "time missing operand",
			input: "{{time.now.plus}}",
			err:   errors.New("missing time operand"),
		},

		{
			name:   "no_templating",
			input:  "path foobar {",
			output: "path foobar {",
		},
		{
			name:  "only_closing",
			input: "path foobar}} {",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:  "closing_in_front",
			input: "path }} {{foobar}} {",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:  "closing_in_back",
			input: "path {{foobar}} }}",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:   "basic",
			input:  "path /{{identity.entity.id}}/ {",
			output: "path /entityID/ {",
		},
		{
			name:       "multiple",
			input:      "path {{identity.entity.name}} {\n\tval = {{identity.entity.metadata.foo}}\n}",
			entityName: "entityName",
			metadata:   map[string]string{"foo": "bar"},
			output:     "path entityName {\n\tval = bar\n}",
		},
		{
			name:     "multiple_bad_name",
			input:    "path {{identity.entity.name}} {\n\tval = {{identity.entity.metadata.foo}}\n}",
			metadata: map[string]string{"foo": "bar"},
			err:      ErrTemplateValueNotFound,
		},
		{
			name:  "unbalanced_close",
			input: "path {{identity.entity.id}} {\n\tval = {{ent}}ity.metadata.foo}}\n}",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:  "unbalanced_open",
			input: "path {{identity.entity.id}} {\n\tval = {{ent{{ity.metadata.foo}}\n}",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:      "no_entity_no_directives",
			input:     "path {{identity.entity.id}} {\n\tval = {{ent{{ity.metadata.foo}}\n}",
			err:       ErrNoEntityAttachedToToken,
			nilEntity: true,
		},
		{
			name:      "no_entity_no_diretives",
			input:     "path name {\n\tval = foo\n}",
			output:    "path name {\n\tval = foo\n}",
			nilEntity: true,
		},
		{
			name:          "alias_id_name",
			input:         "path {{ identity.entity.name}} {\n\tval = {{identity.entity.aliases.foomount.id}}\n}",
			entityName:    "entityName",
			aliasAccessor: "foomount",
			aliasID:       "aliasID",
			metadata:      map[string]string{"foo": "bar"},
			output:        "path entityName {\n\tval = aliasID\n}",
		},
		{
			name:          "alias_id_name_bad_selector",
			input:         "path foobar {\n\tval = {{identity.entity.aliases.foomount}}\n}",
			aliasAccessor: "foomount",
			err:           errors.New("invalid alias selector"),
		},
		{
			name:          "alias_id_name_bad_accessor",
			input:         "path \"foobar\" {\n\tval = {{identity.entity.aliases.barmount.id}}\n}",
			aliasAccessor: "foomount",
			err:           errors.New("alias not found"),
		},
		{
			name:          "alias_id_name",
			input:         "path \"{{identity.entity.name}}\" {\n\tval = {{identity.entity.aliases.foomount.metadata.zip}}\n}",
			entityName:    "entityName",
			aliasAccessor: "foomount",
			aliasID:       "aliasID",
			metadata:      map[string]string{"foo": "bar"},
			aliasMetadata: map[string]string{"zip": "zap"},
			output:        "path \"entityName\" {\n\tval = zap\n}",
		},
		{
			name:       "group_name",
			input:      "path \"{{identity.groups.ids.groupID.name}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			output:     "path \"groupName\" {\n\tval = entityName\n}",
		},
		{
			name:       "group_bad_id",
			input:      "path \"{{identity.groups.ids.hroupID.name}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			err:        errors.New("entity is not a member of group \"hroupID\""),
		},
		{
			name:       "group_id",
			input:      "path \"{{identity.groups.names.groupName.id}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			output:     "path \"groupID\" {\n\tval = entityName\n}",
		},
		{
			name:       "group_bad_name",
			input:      "path \"{{identity.groups.names.hroupName.id}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			err:        errors.New("entity is not a member of group \"hroupName\""),
		},
		{
			name:     "metadata_object_disallowed",
			input:    "{{identity.entity.metadata}}",
			metadata: map[string]string{"foo": "bar"},
			err:      ErrTemplateValueNotFound,
		},
		{
			name:          "alias_metadata_object_disallowed",
			input:         "{{identity.entity.aliases.foomount.metadata}}",
			aliasAccessor: "foomount",
			aliasMetadata: map[string]string{"foo": "bar"},
			err:           ErrTemplateValueNotFound,
		},
		{
			name:             "group_names_disallowed",
			input:            "{{identity.entity.group_names}}",
			groupMemberships: []string{"foo", "bar"},
			err:              ErrTemplateValueNotFound,
		},
		{
			name:             "group_ids_disallowed",
			input:            "{{identity.entity.group_ids}}",
			groupMemberships: []string{"foo", "bar"},
			err:              ErrTemplateValueNotFound,
		},

		// missing selector cases
		{
			mode:   JSONTemplating,
			name:   "entity id",
			input:  "{{identity.entity.id}}",
			output: `"entityID"`,
		},
		{
			mode:       JSONTemplating,
			name:       "entity name",
			input:      "{{identity.entity.name}}",
			entityName: "entityName",
			output:     `"entityName"`,
		},
		{
			mode:   JSONTemplating,
			name:   "entity name missing",
			input:  "{{identity.entity.name}}",
			output: `""`,
		},
		{
			mode:     JSONTemplating,
			name:     "one metadata key",
			input:    "{{identity.entity.metadata.color}}",
			metadata: map[string]string{"foo": "bar", "color": "green"},
			output:   `"green"`,
		},
		{
			mode:     JSONTemplating,
			name:     "one metadata key not found",
			input:    "{{identity.entity.metadata.size}}",
			metadata: map[string]string{"foo": "bar", "color": "green"},
			output:   `""`,
		},
		{
			mode:     JSONTemplating,
			name:     "all entity metadata",
			input:    "{{identity.entity.metadata}}",
			metadata: map[string]string{"foo": "bar", "color": "green"},
			output:   `{"color":"green","foo":"bar"}`,
		},
		{
			mode:             JSONTemplating,
			name:             "group_names",
			input:            "{{identity.entity.group_names}}",
			groupMemberships: []string{"foo", "bar"},
			output:           `["foo","bar"]`,
		},
		{
			mode:             JSONTemplating,
			name:             "group_ids",
			input:            "{{identity.entity.group_ids}}",
			groupMemberships: []string{"foo", "bar"},
			output:           `["foo_0","bar_1"]`,
		},
		{
			mode:          JSONTemplating,
			name:          "one alias metadata key",
			input:         "{{identity.entity.aliases.aws_123.metadata.color}}",
			aliasAccessor: "aws_123",
			aliasMetadata: map[string]string{"foo": "bar", "color": "green"},
			output:        `"green"`,
		},
		{
			mode:          JSONTemplating,
			name:          "one alias metadata key not found",
			input:         "{{identity.entity.aliases.aws_123.metadata.size}}",
			aliasAccessor: "aws_123",
			aliasMetadata: map[string]string{"foo": "bar", "color": "green"},
			output:        `""`,
		},
		{
			mode:          JSONTemplating,
			name:          "one alias metadata, accessor not found",
			input:         "{{identity.entity.aliases.aws_123.metadata.size}}",
			aliasAccessor: "not_gonna_match",
			aliasMetadata: map[string]string{"foo": "bar", "color": "green"},
			output:        `""`,
		},
		{
			mode:          JSONTemplating,
			name:          "all alias metadata",
			input:         "{{identity.entity.aliases.aws_123.metadata}}",
			aliasAccessor: "aws_123",
			aliasMetadata: map[string]string{"foo": "bar", "color": "green"},
			output:        `{"color":"green","foo":"bar"}`,
		},
		{
			mode:          JSONTemplating,
			name:          "all alias metadata, accessor not found",
			input:         "{{identity.entity.aliases.aws_123.metadata}}",
			aliasAccessor: "not_gonna_match",
			aliasMetadata: map[string]string{"foo": "bar", "color": "green"},
			output:        `{}`,
		},
	}

	for _, test := range tests {
		var entity *Entity
		if !test.nilEntity {
			entity = &Entity{
				ID:       "entityID",
				Name:     test.entityName,
				Metadata: test.metadata,
			}
		}
		if test.aliasAccessor != "" {
			entity.Aliases = []*Alias{
				&Alias{
					MountAccessor: test.aliasAccessor,
					ID:            test.aliasID,
					Name:          test.aliasName,
					Metadata:      test.aliasMetadata,
				},
			}
		}
		var groups []*Group
		if test.groupName != "" {
			groups = append(groups, &Group{
				ID:          "groupID",
				Name:        test.groupName,
				Metadata:    test.groupMetadata,
				NamespaceID: namespace.RootNamespace.ID,
			})
		}

		if test.groupMemberships != nil {
			for i, groupName := range test.groupMemberships {
				groups = append(groups, &Group{
					ID:   fmt.Sprintf("%s_%d", groupName, i),
					Name: groupName,
				})
			}
		}

		subst, out, err := PopulateString(&PopulateStringInput{
			Mode:              test.mode,
			ValidityCheckOnly: test.validityCheckOnly,
			String:            test.input,
			Entity:            entity,
			Groups:            groups,
			Namespace:         namespace.RootNamespace,
		})
		if err != nil {
			if test.err == nil {
				t.Fatalf("%s: expected success, got error: %v", test.name, err)
			}
			if err.Error() != test.err.Error() {
				t.Fatalf("%s: got error: %v", test.name, err)
			}
		}
		if out != test.output {
			t.Fatalf("%s: bad output: %s, expected: %s", test.name, out, test.output)
		}
		if err == nil && !subst && out != test.input {
			t.Fatalf("%s: bad subst flag", test.name)
		}
	}
}

func TestPopulate_FullObject(t *testing.T) {
	var testEntity = &Entity{
		ID:   "abc-123",
		Name: "Entity Name",
		Metadata: map[string]string{
			"color":         "green",
			"size":          "small",
			"non-printable": "\"\n\t",
		},
		Aliases: []*Alias{
			{
				MountAccessor: "aws_123",
				Metadata: map[string]string{
					"service": "ec2",
					"region":  "west",
				},
			},
		},
	}

	var testGroups = []*Group{
		{ID: "a08b0c02", Name: "g1"},
		{ID: "239bef91", Name: "g2"},
	}

	template := `
			{
			    "id": {{identity.entity.id}},
			    "name": {{identity.entity.name}},
			    "all metadata": {{identity.entity.metadata}},
			    "one metadata key": {{identity.entity.metadata.color}},
			    "one metadata key not found": {{identity.entity.metadata.asldfk}},
			    "alias metadata": {{identity.entity.aliases.aws_123.metadata}},
			    "alias not found metadata": {{identity.entity.aliases.blahblah.metadata}},
			    "one alias metadata key": {{identity.entity.aliases.aws_123.metadata.service}},
			    "one not found alias metadata key": {{identity.entity.aliases.blahblah.metadata.service}},
			    "group names": {{identity.entity.group_names}},
			    "group ids": {{identity.entity.group_ids}},
			    "time now": {{time.now}},
			    "time plus": {{time.now.plus.10s}},
			    "time minus": {{time.now.minus.1h}},
			    "repeated and": {"nested element": {{identity.entity.name}}}
			}`

	expected := `
			{
			    "id": "abc-123",
			    "name": "Entity Name",
			    "all metadata": {"color":"green","non-printable":"\"\n\t","size":"small"},
			    "one metadata key": "green",
			    "one metadata key not found": "",
			    "alias metadata": {"region":"west","service":"ec2"},
			    "alias not found metadata": {},
			    "one alias metadata key": "ec2",
			    "one not found alias metadata key": "",
			    "group names": ["g1","g2"],
			    "group ids": ["a08b0c02","239bef91"],
			    "time now": <now>,
			    "time plus": <plus>,
			    "time minus": <minus>,
			    "repeated and": {"nested element": "Entity Name"}
			}`

	now := time.Now().Unix()
	minus := now - (60 * 60)
	plus := now + 10

	expected = strings.ReplaceAll(expected, "<now>", strconv.FormatInt(now, 10))
	expected = strings.ReplaceAll(expected, "<plus>", strconv.FormatInt(plus, 10))
	expected = strings.ReplaceAll(expected, "<minus>", strconv.FormatInt(minus, 10))

	input := &PopulateStringInput{
		Mode:   JSONTemplating,
		String: template,
		Entity: testEntity,
		Groups: testGroups,
	}
	_, out, err := PopulateString(input)
	if err != nil {
		t.Fatal(err)
	}

	if out != expected {
		t.Fatalf("expected:\n%s\n\ngot:\n%s", expected, out)
	}
}
