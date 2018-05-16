package mop

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/stretchr/testify/assert"
)

func TestStale(t *testing.T) {
	now := aws.Time(time.Now())
	before := aws.Time(time.Now().Add(time.Hour * 24 * -10))
	assert := assert.New(t)

	t.Run("ignoreTags", func(t *testing.T) {
		input := []ecr.ImageDetail{
			{
				ImagePushedAt: now,
				ImageTags:     []string{"foo"},
			},
			{
				ImagePushedAt: now,
				ImageTags:     []string{"latest"},
			},
		}
		expected := []ecr.ImageDetail{
			{
				ImagePushedAt: now,
				ImageTags:     []string{"latest"},
			},
		}

		m, err := New("test", 0, false, []string{"foo"})
		assert.Equal(nil, err)
		assert.Equal(expected, m.stale(input))
	})

	t.Run("untagged", func(t *testing.T) {
		input := []ecr.ImageDetail{
			{
				ImagePushedAt: now,
				ImageTags:     []string{},
			},
			{
				ImagePushedAt: now,
				ImageTags:     []string{"latest"},
			},
		}
		expected := []ecr.ImageDetail{
			{
				ImagePushedAt: now,
				ImageTags:     []string{},
			},
		}

		m, err := New("test", 5, true, []string{})
		assert.Equal(nil, err)
		assert.Equal(expected, m.stale(input))
	})

	t.Run("staleAfter", func(t *testing.T) {
		input := []ecr.ImageDetail{
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
		}
		expected := []ecr.ImageDetail{
			{
				ImagePushedAt: before,
				ImageTags:     []string{"foo"},
			},
			{
				ImagePushedAt: before,
				ImageTags:     []string{"bar"},
			},
		}

		m, err := New("test", 5, false, []string{})
		assert.Equal(nil, err)
		assert.Equal(expected, m.stale(input))
	})
}

func TestIdentifiersFrom(t *testing.T) {
	input := []ecr.ImageDetail{
		{
			ImageDigest: aws.String("foo123"),
			ImageTags:   []string{"foo"},
		},
		{
			ImageDigest: aws.String("bar123"),
			ImageTags:   []string{"bar"},
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
