package aggregator

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication/mock"
	"github.com/stretchr/testify/require"
)

func TestGraphqlResponseGetter_Query(t *testing.T) {

	expectedErr := errors.New("expected error")

	t.Run("native auth client errors should error", func(t *testing.T) {
		t.Parallel()

		responseGetter := &graphqlResponseGetter{
			authClient: &mock.NativeStub{GetAccessTokenCalled: func() (string, error) {
				return "", expectedErr
			}},
		}

		query, err := responseGetter.Query(context.Background(), "", "", "")
		require.Nil(t, query)
		require.Equal(t, expectedErr, err)
	})
	t.Run("invalid URL should error", func(t *testing.T) {
		t.Parallel()

		responseGetter := &graphqlResponseGetter{
			authClient: &mock.NativeStub{GetAccessTokenCalled: func() (string, error) {
				return "accessToken", nil
			}},
		}

		query, err := responseGetter.Query(context.Background(), "invalid URL", "", "")
		require.Nil(t, query)
		require.NotNil(t, err)
		require.IsType(t, err, &url.Error{})
	})
}
