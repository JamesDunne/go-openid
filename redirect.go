package openid

import (
	"net/url"
	"strings"
)

func RedirectUrl(id, callbackUrl, realm string, opArgs map[string]string) (string, error) {
	return redirectUrl(id, callbackUrl, realm, opArgs, urlGetter)
}

func redirectUrl(id, callbackUrl, realm string, opArgs map[string]string, getter httpGetter) (string, error) {
	opEndpoint, opLocalId, claimedId, err := discover(id, getter)
	if err != nil {
		return "", err
	}
	return buildRedirectUrl(opEndpoint, opLocalId, claimedId, callbackUrl, realm, opArgs)
}

func buildRedirectUrl(opEndpoint, opLocalId, claimedId, returnTo, realm string, opArgs map[string]string) (string, error) {
	values := make(url.Values)
	values.Add("openid.ns", "http://specs.openid.net/auth/2.0")
	values.Add("openid.mode", "checkid_setup")
	values.Add("openid.return_to", returnTo)

	if len(claimedId) > 0 {
		values.Add("openid.claimed_id", claimedId)
		if len(opLocalId) > 0 {
			values.Add("openid.identity", opLocalId)
		} else {
			values.Add("openid.identity",
				"http://specs.openid.net/auth/2.0/identifier_select")
		}
	} else {
		values.Add("openid.identity",
			"http://specs.openid.net/auth/2.0/identifier_select")
	}

	if len(realm) > 0 {
		values.Add("openid.realm", realm)
	}

	// ADDED(jsd): Pass extra openid args
	if opArgs != nil {
		for key, val := range opArgs {
			values.Add(key, val)
		}
	}

	if strings.Contains(opEndpoint, "?") {
		return opEndpoint + "&" + values.Encode(), nil
	}
	return opEndpoint + "?" + values.Encode(), nil
}
