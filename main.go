package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
)

var (
	clientName = "1password-aws-sso"
	clientType = "public"
	grantType  = "urn:ietf:params:oauth:grant-type:device_code"
)

func launchBrowser(url string) {
	var cmd string
	var args []string

	switch o := runtime.GOOS; o {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		log.Fatalf("Unsupported OS: %s", o)
	}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var profileName string
	var accessToken string

	ctx := context.Background()

	flag.StringVar(&profileName, "profile", "", "AWS profile to use")
	flag.Parse()

	if profileName == "" {
		log.Fatal("no profile specified")
	}

	profile, err := config.LoadSharedConfigProfile(ctx, profileName)
	if err != nil {
		log.Fatalf("failed to load profile %s: %s", profileName, err)
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profileName))
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	oidcClient := ssooidc.NewFromConfig(cfg)
	ssoClient := sso.NewFromConfig(cfg)

	registerClientInput := ssooidc.RegisterClientInput{ClientName: &clientName, ClientType: &clientType}
	registerClientOutput, err := oidcClient.RegisterClient(ctx, &registerClientInput)
	if err != nil {
		log.Fatalf("failed to register client: %s", err)
	}

	startDeveiceAuthorizationOutput, err := oidcClient.StartDeviceAuthorization(
		ctx,
		&ssooidc.StartDeviceAuthorizationInput{ClientId: registerClientOutput.ClientId, ClientSecret: registerClientOutput.ClientSecret, StartUrl: &profile.SSOStartURL},
	)
	if err != nil {
		log.Fatalf("failed to start device authorization: %s", err)
	}

	createTokenInput := ssooidc.CreateTokenInput{ClientId: registerClientOutput.ClientId, ClientSecret: registerClientOutput.ClientSecret, DeviceCode: startDeveiceAuthorizationOutput.DeviceCode, GrantType: &grantType}

	launchBrowser(*startDeveiceAuthorizationOutput.VerificationUriComplete)

	for {
		createTokenOutput, err := oidcClient.CreateToken(ctx, &createTokenInput)
		if err != nil {
			var authorizationPendingException *types.AuthorizationPendingException
			if errors.As(err, &authorizationPendingException) {
				time.Sleep(1 * time.Second)
				continue
			} else {
				log.Fatalf("failed to create token: %s", err)
			}
		}

		accessToken = *createTokenOutput.AccessToken
		break
	}

	rci := &sso.GetRoleCredentialsInput{AccountId: &profile.SSOAccountID, RoleName: &profile.SSORoleName, AccessToken: &accessToken}
	creds, err := ssoClient.GetRoleCredentials(ctx, rci)
	if err != nil {
		log.Fatalf("failed to get role credentials: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(creds.RoleCredentials); err != nil {
		log.Fatalf("failed to encode credentials: %s", err)
	}
}
