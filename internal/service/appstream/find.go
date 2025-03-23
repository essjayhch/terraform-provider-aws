// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appstream

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindImageBuilderByName(ctx context.Context, conn *appstream.Client, name string) (*awstypes.ImageBuilder, error) {
	input := &appstream.DescribeImageBuildersInput{
		Names: []string{name},
	}

	return findImageBuilder(ctx, conn, input)
}

func findImageBuilders(ctx context.Context, conn *appstream.Client, input *appstream.DescribeImageBuildersInput) ([]awstypes.ImageBuilder, error) {
	var output []awstypes.ImageBuilder

	err := describeImageBuildersPages(ctx, conn, input, func(page *appstream.DescribeImageBuildersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		output = append(output, page.ImageBuilders...)

		return !lastPage
	})

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func findImageBuilder(ctx context.Context, conn *appstream.Client, input *appstream.DescribeImageBuildersInput) (*awstypes.ImageBuilder, error) {
	output, err := findImageBuilders(ctx, conn, input)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSingleValueResult(output)
}

func FindUserByTwoPartKey(ctx context.Context, conn *appstream.Client, username, authType string) (*awstypes.User, error) {
	input := &appstream.DescribeUsersInput{
		AuthenticationType: awstypes.AuthenticationType(authType),
	}

	return findUser(ctx, conn, input, func(v *awstypes.User) bool {
		return aws.ToString(v.UserName) == username
	})
}

func findUsers(ctx context.Context, conn *appstream.Client, input *appstream.DescribeUsersInput, filter tfslices.Predicate[*awstypes.User]) ([]awstypes.User, error) {
	var output []awstypes.User

	err := describeUsersPages(ctx, conn, input, func(page *appstream.DescribeUsersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Users {
			if filter(&v) {
				output = append(output, v)
			}
		}

		return !lastPage
	})

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func findUser(ctx context.Context, conn *appstream.Client, input *appstream.DescribeUsersInput, filter tfslices.Predicate[*awstypes.User]) (*awstypes.User, error) {
	output, err := findUsers(ctx, conn, input, filter)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSingleValueResult(output)
}

// findImages finds all images from a describe images input
func findImages(ctx context.Context, conn *appstream.Client, input *appstream.DescribeImagesInput) ([]awstypes.Image, error) {
	var output []awstypes.Image

	pages := appstream.NewDescribeImagesPaginator(conn, input)

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: input,
			}
		}

		if err != nil {
			return nil, err
		}

		output = append(output, page.Images...)
	}

	return output, nil
}
