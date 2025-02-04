// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/trapajim/snapmatch-ai/templates/models"

func Jobs(jobs models.Jobs, expiry string) templ.Component {
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"container mx-auto py-6\"><div class=\"flex items-center justify-between\"><h1 class=\"text-2xl font-bold text-gray-800 mb-4\">Jobs</h1></div><div id=\"job-wrapper\"><table class=\"table-auto w-full bg-white shadow-md rounded-lg\"><thead class=\"bg-gray-100 border-b\"><tr><th class=\"px-4 py-2 text-left text-gray-600 font-medium\">ID</th><th class=\"px-4 py-2 text-left text-gray-600 font-medium\">Name</th><th class=\"px-4 py-2 text-left text-gray-600 font-medium\">Status</th></tr></thead> <tbody>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, job := range jobs {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<tr class=\"@if index % 2 == 0 { bg-gray-50 } else { bg-white } hover:bg-gray-100\"><td class=\"px-4 py-3 text-gray-800 text-sm\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(job.ID)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/jobs.templ`, Line: 23, Col: 60}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</td><td class=\"px-4 py-3 text-gray-800 text-sm\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(job.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/jobs.templ`, Line: 24, Col: 62}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</td><td class=\"px-4 py-3 text-sm\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if job.Status == "Running" {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<span class=\"px-2 py-1 text-xs font-semibold text-blue-700 bg-blue-100 rounded\">Running</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if job.Status == "Pending" {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "<span class=\"px-2 py-1 text-xs font-semibold text-gray-700 bg-blue-100 rounded\">Pending</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if job.Status == "Failed" {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<span class=\"px-2 py-1 text-xs font-semibold text-red-700 bg-red-100 rounded\">Failed</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else if job.Status == "Success" {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "<span class=\"px-2 py-1 text-xs font-semibold text-green-700 bg-green-100 rounded\">Success</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "<span class=\"px-2 py-1 text-xs font-semibold text-gray-700 bg-gray-100 rounded\">Unknown</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</td></tr>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</tbody></table></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = Page("Jobs", expiry).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
