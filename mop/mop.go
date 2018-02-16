package mop

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// Mop represents an ECR mop.
type Mop struct {
	ecr          *ecr.ECR
	ignoreTags   []string
	repository   string
	staleAfter   time.Duration
	wipeUntagged bool
}

// Result represents the result of a mop wipe.
type Result struct {
	Removed []ecr.ImageIdentifier
	Failed  []ecr.ImageFailure
}

// New creates an ECR mop.
func New(repository string, staleAfter int, wipeUntagged bool, ignoreTags []string) (*Mop, error) {
	c, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	return &Mop{
		ecr:          ecr.New(c),
		ignoreTags:   ignoreTags,
		repository:   repository,
		staleAfter:   time.Hour * 24 * time.Duration(staleAfter),
		wipeUntagged: wipeUntagged,
	}, nil
}

// Wipe removes stale images from the provided ECR repositories.
func (m *Mop) Wipe() (*Result, error) {
	wres := &Result{}
	reqd := m.ecr.DescribeImagesRequest(&ecr.DescribeImagesInput{RepositoryName: &m.repository})
	page := reqd.Paginate()

	for page.Next() {
		if page.Err() != nil {
			return wres, page.Err()
		}

		r, err := m.delete(m.identifiersFrom(m.stale(page.CurrentPage().ImageDetails)))
		if err != nil {
			return wres, err
		}

		wres.Removed = append(wres.Removed, r.ImageIds...)
		wres.Failed = append(wres.Failed, r.Failures...)
	}

	return wres, nil
}

func (m *Mop) delete(images []ecr.ImageIdentifier) (*ecr.BatchDeleteImageOutput, error) {
	if len(images) == 0 {
		return &ecr.BatchDeleteImageOutput{}, nil
	}

	r, err := m.ecr.BatchDeleteImageRequest(&ecr.BatchDeleteImageInput{RepositoryName: &m.repository, ImageIds: images}).Send()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *Mop) stale(images []ecr.ImageDetail) (stale []ecr.ImageDetail) {
	for _, image := range images {
		if m.shouldIgnore(image) {
			continue
		}
		if m.wipeUntagged && len(image.ImageTags) == 0 {
			stale = append(stale, image)
		}
		if image.ImagePushedAt.Before(time.Now().Add(-m.staleAfter)) {
			stale = append(stale, image)
		}
	}
	return
}

func (m *Mop) identifiersFrom(images []ecr.ImageDetail) (identifiers []ecr.ImageIdentifier) {
	for _, image := range images {
		identifiers = append(identifiers, ecr.ImageIdentifier{ImageDigest: image.ImageDigest})
	}
	return
}

func (m *Mop) shouldIgnore(image ecr.ImageDetail) bool {
	for _, ignoreTag := range m.ignoreTags {
		for _, imageTag := range image.ImageTags {
			if ignoreTag == imageTag {
				return true
			}
		}
	}
	return false
}
