package keycloak

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	urlSeparator string = "/"
)

type GoCloak struct {
	client  *gocloak.GoCloak
	baseURL string
}

func NewGoCloakClient(baseURL string) *GoCloak {
	return &GoCloak{
		client:  gocloak.NewClient(baseURL),
		baseURL: baseURL,
	}
}

func (g *GoCloak) GetAccountGroups(ctx context.Context, token, realm string) ([]*gocloak.Group, error) {
	const errMessage = "could not get groups"

	var result []*gocloak.Group

	resp, err := g.client.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getRealmURL(realm, "account", "groups"))
	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

func (g *GoCloak) getRealmURL(realm string, path ...string) string {
	path = append([]string{g.baseURL, "realms", realm}, path...)
	return makeURL(path...)
}

func makeURL(path ...string) string {
	return strings.Join(path, urlSeparator)
}

func checkForError(resp *resty.Response, err error, errMessage string) error {
	if err != nil {
		return &gocloak.APIError{
			Code:    0,
			Message: errors.Wrap(err, errMessage).Error(),
			Type:    gocloak.ParseAPIErrType(err),
		}
	}

	if resp == nil {
		return &gocloak.APIError{
			Message: "empty response",
			Type:    gocloak.ParseAPIErrType(err),
		}
	}

	if resp.IsError() {
		var msg string

		if e, ok := resp.Error().(*gocloak.HTTPErrorResponse); ok && e.NotEmpty() {
			msg = fmt.Sprintf("%s: %s", resp.Status(), e)
		} else {
			msg = resp.Status()
		}

		return &gocloak.APIError{
			Code:    resp.StatusCode(),
			Message: msg,
			Type:    gocloak.ParseAPIErrType(err),
		}
	}

	return nil
}
