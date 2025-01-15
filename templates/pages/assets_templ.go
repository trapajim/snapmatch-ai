// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/trapajim/snapmatch-ai/templates/models"
	"github.com/trapajim/snapmatch-ai/templates/partials"
)

func Assets(assets models.Assets) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- drawer component --> <div id=\"drawer-example\" class=\"fixed top-0 left-0 z-40 h-screen p-6 overflow-y-auto transition-transform -translate-x-full bg-white w-80 shadow-lg rounded-r-lg dark:bg-gray-800\" tabindex=\"-1\" aria-labelledby=\"drawer-label\"><!-- Drawer Header --><div class=\"flex items-center justify-between mb-6\"><h5 id=\"drawer-label\" class=\"inline-flex items-center text-lg font-semibold text-gray-700 dark:text-gray-300\"><svg width=\"20px\" height=\"20px\" class=\"fill-gray-200 mr-2\" viewBox=\"0 0 48 48\" xmlns=\"http://www.w3.org/2000/svg\"><title>ai</title><g id=\"Layer_2\" data-name=\"Layer 2\"><g id=\"invisible_box\" data-name=\"invisible box\"><rect width=\"48\" height=\"48\" fill=\"none\"></rect></g> <g id=\"Q3_icons\" data-name=\"Q3 icons\"><g><path d=\"M45.6,18.7,41,14.9V7.5a1,1,0,0,0-.6-.9L30.5,2.1h-.4l-.6.2L24,5.9,18.5,2.2,17.9,2h-.4L7.6,6.6a1,1,0,0,0-.6.9v7.4L2.4,18.7a.8.8,0,0,0-.4.8v9H2a.8.8,0,0,0,.4.8L7,33.1v7.4a1,1,0,0,0,.6.9l9.9,4.5h.4l.6-.2L24,42.1l5.5,3.7.6.2h.4l9.9-4.5a1,1,0,0,0,.6-.9V33.1l4.6-3.8a.8.8,0,0,0,.4-.7V19.4h0A.8.8,0,0,0,45.6,18.7Zm-5.1,6.8H42v1.6l-3.5,2.8-.4.3-.4-.2a1.4,1.4,0,0,0-2,.7,1.5,1.5,0,0,0,.6,2l.7.3h0v5.4l-6.6,3.1-4.2-2.8-.7-.5V25.5H27a1.5,1.5,0,0,0,0-3H25.5V9.7l.7-.5,4.2-2.8L37,9.5v5.4h0l-.7.3a1.5,1.5,0,0,0-.6,2,1.4,1.4,0,0,0,1.3.9l.7-.2.4-.2.4.3L42,20.9v1.6H40.5a1.5,1.5,0,0,0,0,3ZM21,25.5h1.5V38.3l-.7.5-4.2,2.8L11,38.5V33.1h0l.7-.3a1.5,1.5,0,0,0,.6-2,1.4,1.4,0,0,0-2-.7l-.4.2-.4-.3L6,27.1V25.5H7.5a1.5,1.5,0,0,0,0-3H6V20.9l3.5-2.8.4-.3.4.2.7.2a1.4,1.4,0,0,0,1.3-.9,1.5,1.5,0,0,0-.6-2L11,15h0V9.5l6.6-3.1,4.2,2.8.7.5V22.5H21a1.5,1.5,0,0,0,0,3Z\"></path> <path d=\"M13.9,9.9a1.8,1.8,0,0,0,0,2.2l2.6,2.5v2.8l-4,4v5.2l4,4v2.8l-2.6,2.5a1.8,1.8,0,0,0,0,2.2,1.5,1.5,0,0,0,1.1.4,1.5,1.5,0,0,0,1.1-.4l3.4-3.5V29.4l-4-4V22.6l4-4V13.4L16.1,9.9A1.8,1.8,0,0,0,13.9,9.9Z\"></path> <path d=\"M31.5,14.6l2.6-2.5a1.8,1.8,0,0,0,0-2.2,1.8,1.8,0,0,0-2.2,0l-3.4,3.5v5.2l4,4v2.8l-4,4v5.2l3.4,3.5a1.7,1.7,0,0,0,2.2,0,1.8,1.8,0,0,0,0-2.2l-2.6-2.5V30.6l4-4V21.4l-4-4Z\"></path></g></g></g></svg> <span>Jobs</span></h5><button type=\"button\" data-drawer-hide=\"drawer-example\" aria-controls=\"drawer-example\" class=\"text-gray-500 hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 flex items-center justify-center dark:hover:bg-gray-600 dark:hover:text-white\"><svg class=\"w-4 h-4\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 14 14\" aria-hidden=\"true\"><path stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6\"></path></svg> <span class=\"sr-only\">Close menu</span></button></div><!-- Drawer Content --><div class=\"space-y-4\"><div><h6 class=\"text-gray-900 text-lg font-medium dark:text-gray-200\">Grouping</h6><p class=\"text-sm text-gray-600 dark:text-gray-400\">This job will group images into categories based on their content.</p></div><form hx-post=\"/assets/predict\" hx-swap=\"none\" hx-indicator=\"next\"><div class=\"mb-4\"><label for=\"categories\" class=\"block text-sm font-medium text-gray-900 dark:text-gray-200\">Categories</label><p class=\"text-xs text-gray-600 dark:text-gray-400 mb-2\">(Optional) Specify categories separated by commas, e.g., \"Electronics, Food, Shoes.\"</p><textarea name=\"categories\" id=\"categories\" class=\"w-full px-3 py-2 border rounded-lg text-gray-900 bg-gray-50 border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-800 dark:text-gray-100 dark:border-gray-600 dark:focus:ring-blue-500 dark:focus:border-blue-500\" placeholder=\"Electronics, Food, Shoes\" rows=\"3\"></textarea></div><button type=\"submit\" class=\"w-full text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 transition dark:bg-blue-500 dark:hover:bg-blue-600 dark:focus:ring-blue-800\">Start Categorization Job</button></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = partials.Indicator().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</div></div><div class=\"container mx-auto py-6\"><div class=\"flex items-center justify-between\"><h1 class=\"text-2xl font-bold text-gray-800 mb-4\">Assets</h1><div class=\"text-center\"><button class=\"text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800\" type=\"button\" data-drawer-target=\"drawer-example\" data-drawer-show=\"drawer-example\" aria-controls=\"drawer-example\">Jobs</button></div></div><div class=\"mb-4\"><form hx-get=\"/assets\" hx-disabled-elt=\"#find\" hx-indicator=\"next\" hx-target=\"#assets-wrapper\" class=\"flex flex-col sm:flex-row items-start sm:items-center space-y-4 sm:space-y-0 sm:space-x-4\"><div class=\"flex flex-col flex-grow\"><label for=\"query\" class=\"flex items-center text-sm font-medium text-gray-700\">Search Query</label> <input type=\"text\" name=\"query\" id=\"query\" placeholder=\"Search assets...\" class=\"flex-grow px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring focus:ring-blue-200\"></div><div class=\"flex flex-col\"><label for=\"similarity\" class=\"flex items-center text-sm font-medium text-gray-700\">Similarity <span class=\"ml-1 relative group\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\" class=\"w-4 h-4 text-gray-500 cursor-pointer\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M13 16h-1v-4h-.01M12 8h.01m-9 4a9 9 0 1118 0 9 9 0 01-18 0z\"></path></svg><div class=\"absolute top-full left-1/2 -translate-x-1/2 mt-1 w-48 p-2 bg-gray-800 text-white text-xs rounded-lg shadow-lg opacity-0 max-h-0 overflow-hidden group-hover:opacity-100 group-hover:max-h-40 transition-all duration-300\">Choose the level of similarity for the search results. Higher similarity will return closer matches.</div></span></label> <select id=\"similarity\" name=\"similarity\" class=\"mt-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring focus:ring-blue-200\"><option value=\"high\">High Similarity</option> <option value=\"medium\" selected>Medium Similarity</option> <option value=\"low\">Low Similarity</option></select></div><button type=\"submit\" id=\"find\" class=\"px-4 self-end py-2 bg-blue-500 text-white font-medium rounded-lg hover:bg-blue-600 transition\">Search</button></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = partials.Indicator().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><div id=\"assets-wrapper\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = partials.Assets(assets).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = Page("Assets").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
