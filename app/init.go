package app

import (
	"lelv/app/controllers"
	"lelv/app/models/blog"
	"lelv/app/models/conversation"
	"log"
	"strings"
	"time"

	"github.com/revel/revel"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.TemplateFuncs["htmlcnt"] = func(s string) string {
		r := strings.Replace(s, "&", "&amp;", -1)
		r = strings.Replace(r, " ", "&nbsp;", -1)
		r = strings.Replace(r, "<", "&lt", -1)
		r = strings.Replace(r, ">", "&gt", -1)
		r = strings.Replace(r, "'", "＇", -1)
		r = strings.Replace(r, "\n", "<br>", -1)
		return r
	}

	revel.TemplateFuncs["strcat"] = func(s1, s2 string) string {
		return s1 + s2
	}

	revel.TemplateFuncs["shortTime"] = func(ts string) string {
		return strings.Split(ts, " ")[0]
	}

	revel.TemplateFuncs["getDate"] = func(ts string) string {
		s0 := strings.Split(ts, " ")[0]
		s1 := strings.Split(s0, "-")
		return s1[1] + "|" + s1[2]
	}

	revel.TemplateFuncs["join"] = func(a []string, s string) string {
		return strings.Join(a, s)
	}

	revel.TemplateFuncs["descendingCommentByID"] = func(data []blog.Comment) []blog.Comment {
		for i := 0; i < len(data); i++ {
			for j := 0; j < len(data)-i-1; j++ {
				if data[j].ID < data[j+1].ID {
					data[j], data[j+1] = data[j+1], data[j]
				}
			}
		}
		return data
	}

	revel.TemplateFuncs["descendingByDateTime"] = func(data interface{}) interface{} {
		switch data.(type) {
		case []blog.Comment:
			for i := 0; i < len(data.([]blog.Comment)); i++ {
				for j := 0; j < len(data.([]blog.Comment))-i-1; j++ {
					t1, err := time.Parse("2006-01-02 15:04:05", data.([]blog.Comment)[j].TimeStamp)
					t2, err := time.Parse("2006-01-02 15:04:05", data.([]blog.Comment)[j+1].TimeStamp)

					if err == nil && t1.Before(t2) {
						data.([]blog.Comment)[j], data.([]blog.Comment)[j+1] = data.([]blog.Comment)[j+1], data.([]blog.Comment)[j]
					}
				}
			}
		case []conversation.Conversation:
			for i := 0; i < len(data.([]conversation.Conversation)); i++ {
				for j := 0; j < len(data.([]conversation.Conversation))-i-1; j++ {
					t1, err := time.Parse("2006-01-02 15:04:05", data.([]conversation.Conversation)[j].TimeStamp)
					t2, err := time.Parse("2006-01-02 15:04:05", data.([]conversation.Conversation)[j+1].TimeStamp)

					if err == nil && t1.Before(t2) {
						data.([]conversation.Conversation)[j], data.([]conversation.Conversation)[j+1] = data.([]conversation.Conversation)[j+1], data.([]conversation.Conversation)[j]
					}
				}
			}
		case []conversation.Message:
			for i := 0; i < len(data.([]conversation.Message)); i++ {
				for j := 0; j < len(data.([]conversation.Message))-i-1; j++ {
					t1, err := time.Parse("2006-01-02 15:04:05", data.([]conversation.Message)[j].TimeStamp)
					t2, err := time.Parse("2006-01-02 15:04:05", data.([]conversation.Message)[j+1].TimeStamp)

					if err == nil && t1.Before(t2) {
						data.([]conversation.Message)[j], data.([]conversation.Message)[j+1] = data.([]conversation.Message)[j+1], data.([]conversation.Message)[j]
					}
				}
			}

		default:
			log.Println("不支持的排序类型")
		}

		return data
	}

	revel.TemplateFuncs["getLastObj"] = func(data []conversation.Message) conversation.Message {
		return data[len(data)-1]
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)

	revel.InterceptFunc(check, revel.BEFORE, &controllers.App{})
	revel.InterceptFunc(check, revel.BEFORE, &controllers.User{})
	revel.InterceptFunc(check, revel.BEFORE, &controllers.Blog{})
	revel.InterceptFunc(check, revel.BEFORE, &controllers.Admin{})
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
