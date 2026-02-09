{{define "markdown"}}
# {{.Headline}}

{{if .OpenToWork}}**Status:** Aberto a propostas 🟢{{end}}

{{if .Bio}}
{{.Bio}}
{{end}}

---

**Info:** {{seniorityLabel (printf "%s" .Seniority)}} • {{.YearsOfExp}} anos de experiência • {{locationLabel (printf "%s" .Location)}} {{if .ContractType}}• {{.ContractType}}{{end}}

{{if .SocialLinks.LinkedIn}}- **LinkedIn:** [{{.SocialLinks.LinkedIn}}]({{.SocialLinks.LinkedIn}}){{end}}
{{if .SocialLinks.GitHub}}- **GitHub:** [{{.SocialLinks.GitHub}}]({{.SocialLinks.GitHub}}){{end}}
{{if .SocialLinks.Website}}- **Website:** [{{.SocialLinks.Website}}]({{.SocialLinks.Website}}){{end}}

{{if .Skills}}
## Habilidades
{{join ", " .Skills}}
{{end}}

{{if .Experiences}}
## Experiência Profissional
{{range .Experiences}}
### **{{.Role}}** | {{.Company}}
*{{formatMonthYear .StartDate}} - {{formatMonthYearPtr .EndDate}}*

{{.Description}}

{{if .TechStack}}**Tech Stack:** {{join ", " .TechStack}}{{end}}

---
{{end}}
{{end}}

{{if .Projects}}
## Projetos
{{range .Projects}}
{{if .ShowOnResume}}
### {{.Name}}
{{.Description}}

{{if .Tags}}**Tags:** {{join ", " .Tags}}{{end}}

{{if .RepoURL}}- [Código Fonte]({{.RepoURL}}){{end}}
{{if .LiveURL}}- [Demo Online]({{.LiveURL}}){{end}}

{{end}}
{{end}}
{{end}}

{{if .Educations}}
## Formação Acadêmica
{{range .Educations}}
- **{{.Degree}} em {{.Field}}** - {{.Institution}}
  *{{formatMonthYear .StartDate}} - {{formatMonthYearPtr .EndDate}}*
{{end}}
{{end}}

---
*Gerado em {{currentDate}} via DevPortfolio*
{{end}}