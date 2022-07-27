package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/carrelage/services"
	"github.com/skatekrak/utils/helpers"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/valyala/fasthttp"
)

func addDefaultRoleToUser(userID string) {
	resp, err := userroles.AddRoleToUser(userID, "user", nil)

	if err != nil {
		log.Printf("Couldn't assign default role to userID(%s): %s", userID, err.Error())
		return
	}

	if resp.UnknownRoleError != nil {
		log.Println("Default role doesn't exists")
		return
	}

	if resp.OK.DidUserAlreadyHaveRole {
		log.Printf("userID(%s) already had the default role", userID)
	}
}

func signInPost(authService *services.AuthService, originalImplementation epmodels.APIInterface) {
	originalSignInPOST := *originalImplementation.SignInPOST

	(*originalImplementation.SignInPOST) = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
		response, err := originalSignInPOST(formFields, options, userContext)
		if err != nil {
			return epmodels.SignInPOSTResponse{}, err
		}

		if response.OK != nil {
			user := response.OK.User

			authService.CreateUserAndProfileIfNotExists(user.ID)
		}

		return response, err
	}
}

func signUpPost(authService *services.AuthService, originalImplementation epmodels.APIInterface) {
	originalSignUpPOST := *originalImplementation.SignUpPOST

	(*originalImplementation.SignUpPOST) = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
		response, err := originalSignUpPOST(formFields, options, userContext)
		if err != nil {
			return epmodels.SignUpPOSTResponse{}, err
		}

		if response.OK != nil {
			user := response.OK.User

			authService.CreateUserAndProfileIfNotExists(user.ID)
			addDefaultRoleToUser(user.ID)
		}

		return response, err
	}
}

// Wrapper to override SuperTokens APIs
func overrideAPIs(authService *services.AuthService) func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
	return func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
		signUpPost(authService, originalImplementation)

		return originalImplementation
	}
}

func InitSuperTokens(authService *services.AuthService) {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: os.Getenv("SUPERTOKENS_CONNECTION_URI"),
			APIKey:        os.Getenv("SUPERTOKENS_API_KEY"),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "carrelage",
			APIDomain:       os.Getenv("API_DOMAIN"),
			WebsiteDomain:   os.Getenv("WEBSITE_DOMAIN"),
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: overrideAPIs(authService),
				},
			}),
			userroles.Init(nil),
			session.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}
}

// Middleware that check if user is logged in
func Logged(options *sessmodels.VerifySessionOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return adaptor.HTTPHandlerFunc(http.HandlerFunc(session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.SetUserContext(r.Context())
			if err := c.Next(); err != nil {
				if err = supertokens.ErrorHandler(err, r, rw); err != nil {
					rw.WriteHeader(500)
					_, _ = rw.Write([]byte(err.Error()))
				}
				return
			}

			c.Response().Header.VisitAll(func(key, value []byte) {
				if string(key) == fasthttp.HeaderContentType {
					rw.Header().Set(string(key), string(value))
				}
			})
			rw.WriteHeader(c.Response().StatusCode())
			_, _ = rw.Write(c.Response().Body())
		})))(c)
	}
}

/**
Throw an Unauthorized error if the user requesting the resource
has not the given role
**/
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionContainer := session.GetSessionFromRequestContext(c.UserContext())
		userId := sessionContainer.GetUserID()

		resp, err := userroles.GetRolesForUser(userId, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Couldn't get your roles")
		}

		if helpers.Has(resp.OK.Roles, role) {
			return c.Next()
		} else {
			return fiber.NewError(fiber.StatusUnauthorized)
		}
	}
}
