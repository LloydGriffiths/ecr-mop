package mop

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/stretchr/testify/assert"
)

var before = func() *time.Time {
	before := time.Now().Add(time.Hour * 24 * -10)
	return &before
}()

var now = func() *time.Time {
	now := time.Now()
	return &now
}()

var igoreTagsFixture = map[string][]ecr.ImageDetail{
	"input": []ecr.ImageDetail{
		{
			ImagePushedAt: now,
			ImageTags:     []string{"foo"},
		},
		{
			ImagePushedAt: now,
			ImageTags:     []string{"latest"},
		},
	},
	"expected": []ecr.ImageDetail{
		{
			ImagePushedAt: now,
			ImageTags:     []string{"latest"},
		},
	},
}

var untaggedFixture = map[string][]ecr.ImageDetail{
	"input": []ecr.ImageDetail{
		{
			ImagePushedAt: now,
			ImageTags:     []string{},
		},
		{
			ImagePushedAt: now,
			ImageTags:     []string{"latest"},
		},
	},
	"expected": []ecr.ImageDetail{
		{
			ImagePushedAt: now,
			ImageTags:     []string{},
		},
	},
}

var staleAfterFixture = map[string][]ecr.ImageDetail{
	"input": []ecr.ImageDetail{
		{
			ImagePushedAt: before,
			ImageTags:     []string{"foo"},
		},
		{
			ImagePushedAt: before,
			ImageTags:     []string{"bar"},
		},
		{
			ImagePushedAt: now,
			ImageTags:     []string{"baz"},
		},
	},
	"expected": []ecr.ImageDetail{
		{
			ImagePushedAt: before,
			ImageTags:     []string{"foo"},
		},
		{
			ImagePushedAt: before,
			ImageTags:     []string{"bar"},
		},
	},
}

func TestStale(t *testing.T) {
	assert := assert.New(t)

	t.Run("ignoreTags", func(t *testing.T) {
		m, err := New("test", 0, false, []string{"foo"})
		assert.Equal(nil, err)
		assert.Equal(igoreTagsFixture["expected"], m.stale(igoreTagsFixture["input"]))
	})

	t.Run("untagged", func(t *testing.T) {
		m, err := New("test", 5, true, []string{})
		assert.Equal(nil, err)
		assert.Equal(untaggedFixture["expected"], m.stale(untaggedFixture["input"]))
	})

	t.Run("staleAfter", func(t *testing.T) {
		m, err := New("test", 5, false, []string{})
		assert.Equal(nil, err)
		assert.Equal(staleAfterFixture["expected"], m.stale(staleAfterFixture["input"]))
	})
}

func TestIdentifiersFrom(t *testing.T) {
	input := []ecr.ImageDetail{
		{
			ImageDigest:   aws.String("foo123"),
			ImagePushedAt: before,
			ImageTags:     []string{"foo"},
		},
		{
			ImageDigest:   aws.String("bar123"),
			ImagePushedAt: before,
			ImageTags:     []string{"bar"},
		},
	}
	expected := []ecr.ImageIdentifier{
		{
			ImageDigest: aws.String("foo123"),
		},
		{
			ImageDigest: aws.String("bar123"),
		},
	}

	assert := assert.New(t)
	m, err := New("test", 5, true, []string{})
	assert.Equal(nil, err)
	assert.Equal(expected, m.identifiersFrom(input))
}
