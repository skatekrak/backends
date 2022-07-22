package auth

import (
	"net/http"
	"os"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/valyala/fasthttp"
)

func InitSuperTokens() {
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
			emailpassword.Init(nil),
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
