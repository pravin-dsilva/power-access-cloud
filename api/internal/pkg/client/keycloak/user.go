package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

// GetUser returns a user with the given userid
func (g *GoCloak) GetUser(ctx context.Context, token, realm string, userID string) (*gocloak.User, error) {
	const errMessage = "could not get user"

	var result *gocloak.User

	resp, err := g.client.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

func (g *GoCloak) getAdminRealmURL(realm string, path ...string) string {
	path = append([]string{g.baseURL, "admin/realms", realm}, path...)
	return makeURL(path...)
}
