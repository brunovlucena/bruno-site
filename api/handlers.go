package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	// 🔒 Security package
	"bruno-api/security"
)

// =============================================================================
// 🎯 PROJECT HANDLERS
// =============================================================================

func getProjects(c *gin.Context) {
	query := `
		SELECT id, title, description, description as short_description, type, github_url, live_url, technologies, active
		FROM projects 
		WHERE active = true
		ORDER BY "order" ASC, id ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var githubURL, liveURL sql.NullString
		var technologies pq.StringArray

		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ShortDescription, &p.Type, &githubURL, &liveURL, &technologies, &p.Active); err != nil {
			continue
		}

		if liveURL.Valid {
			p.LiveURL = liveURL.String
		}
		if githubURL.Valid {
			p.GithubURL = githubURL.String
		}
		p.Technologies = []string(technologies)
		projects = append(projects, p)
	}

	c.JSON(http.StatusOK, projects)
}

func getProject(c *gin.Context) {
	id := c.Param("id")

	// Enhanced input validation using security package
	projectID, validationErr := security.ValidateInteger(id, "id", 1, 999999)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
		return
	}

	var p Project
	var githubURL, liveURL sql.NullString
	var technologies pq.StringArray

	query := `
		SELECT id, title, description, description as short_description, type, github_url, live_url, technologies, active
		FROM projects 
		WHERE id = $1 AND active = true
	`

	err := db.QueryRow(query, projectID).Scan(&p.ID, &p.Title, &p.Description, &p.ShortDescription, &p.Type, &githubURL, &liveURL, &technologies, &p.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	if liveURL.Valid {
		p.LiveURL = liveURL.String
	}
	if githubURL.Valid {
		p.GithubURL = githubURL.String
	}
	p.Technologies = []string(technologies)

	c.JSON(http.StatusOK, p)
}

func createProject(c *gin.Context) {
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Enhanced input validation and sanitization using security package
	title, titleErr := security.ValidateAndSanitizeTitle(project.Title)
	if titleErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": titleErr.Message})
		return
	}
	project.Title = title

	description, descErr := security.ValidateAndSanitizeDescription(project.Description)
	if descErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": descErr.Message})
		return
	}
	project.Description = description

	// Validate type
	if project.Type == "" || len(project.Type) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type length"})
		return
	}
	project.Type = security.SanitizeString(project.Type)

	// Validate and sanitize URLs if provided
	if project.GithubURL != "" {
		githubURL, urlErr := security.ValidateAndSanitizeURL(project.GithubURL, "github_url")
		if urlErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Message})
			return
		}
		project.GithubURL = githubURL
	}

	if project.LiveURL != "" {
		liveURL, urlErr := security.ValidateAndSanitizeURL(project.LiveURL, "live_url")
		if urlErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Message})
			return
		}
		project.LiveURL = liveURL
	}

	query := `
		INSERT INTO projects (title, description, type, github_url, live_url, technologies, active, "order")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id int
	err := db.QueryRow(query, project.Title, project.Description, project.Type, project.GithubURL, project.LiveURL, pq.Array(project.Technologies), project.Active, 0).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	project.ID = id
	c.JSON(http.StatusCreated, project)
}

func updateProject(c *gin.Context) {
	id := c.Param("id")

	// Enhanced input validation using security package
	projectID, validationErr := security.ValidateInteger(id, "id", 1, 999999)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
		return
	}

	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Enhanced input validation and sanitization using security package
	title, titleErr := security.ValidateAndSanitizeTitle(project.Title)
	if titleErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": titleErr.Message})
		return
	}
	project.Title = title

	description, descErr := security.ValidateAndSanitizeDescription(project.Description)
	if descErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": descErr.Message})
		return
	}
	project.Description = description

	// Validate type
	if project.Type == "" || len(project.Type) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type length"})
		return
	}
	project.Type = security.SanitizeString(project.Type)

	// Validate and sanitize URLs if provided
	if project.GithubURL != "" {
		githubURL, urlErr := security.ValidateAndSanitizeURL(project.GithubURL, "github_url")
		if urlErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Message})
			return
		}
		project.GithubURL = githubURL
	}

	if project.LiveURL != "" {
		liveURL, urlErr := security.ValidateAndSanitizeURL(project.LiveURL, "live_url")
		if urlErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Message})
			return
		}
		project.LiveURL = liveURL
	}

	query := `
		UPDATE projects 
		SET title = $1, description = $2, type = $3, github_url = $4, live_url = $5, technologies = $6, active = $7
		WHERE id = $8
	`

	result, err := db.Exec(query, project.Title, project.Description, project.Type, project.GithubURL, project.LiveURL, pq.Array(project.Technologies), project.Active, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

func deleteProject(c *gin.Context) {
	id := c.Param("id")

	// Input validation - ensure id is a valid integer
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	projectID, err := strconv.Atoi(id)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	query := `DELETE FROM projects WHERE id = $1`
	result, err := db.Exec(query, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// =============================================================================
// 🛠️ SKILL HANDLERS
// =============================================================================

func getSkills(c *gin.Context) {
	query := `
		SELECT id, name, category, proficiency, icon, "order" 
		FROM skills 
		ORDER BY "order", name
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var skill Skill
		err := rows.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Proficiency, &skill.Icon, &skill.Order)
		if err != nil {
			continue
		}
		skills = append(skills, skill)
	}

	c.JSON(http.StatusOK, skills)
}

func getSkill(c *gin.Context) {
	id := c.Param("id")

	// Input validation - ensure id is a valid integer
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	skillID, err := strconv.Atoi(id)
	if err != nil || skillID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID format"})
		return
	}

	var skill Skill
	query := `SELECT id, name, category, proficiency, icon, "order" FROM skills WHERE id = $1`

	err = db.QueryRow(query, skillID).Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Proficiency, &skill.Icon, &skill.Order)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skill"})
		return
	}

	c.JSON(http.StatusOK, skill)
}

func createSkill(c *gin.Context) {
	var skill Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Input validation and sanitization
	if skill.Name == "" || len(skill.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill name length"})
		return
	}

	if skill.Category == "" || len(skill.Category) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category length"})
		return
	}

	if skill.Proficiency < 1 || skill.Proficiency > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proficiency level (1-5)"})
		return
	}

	if skill.Icon != "" && len(skill.Icon) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid icon length"})
		return
	}

	query := `
		INSERT INTO skills (name, category, proficiency, icon, "order")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := db.QueryRow(query, skill.Name, skill.Category, skill.Proficiency, skill.Icon, skill.Order).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create skill"})
		return
	}

	skill.ID = id
	c.JSON(http.StatusCreated, skill)
}

func updateSkill(c *gin.Context) {
	id := c.Param("id")

	// Input validation - ensure id is a valid integer
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	skillID, err := strconv.Atoi(id)
	if err != nil || skillID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID format"})
		return
	}

	var skill Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Input validation and sanitization
	if skill.Name == "" || len(skill.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill name length"})
		return
	}

	if skill.Category == "" || len(skill.Category) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category length"})
		return
	}

	if skill.Proficiency < 1 || skill.Proficiency > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proficiency level (1-5)"})
		return
	}

	if skill.Icon != "" && len(skill.Icon) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid icon length"})
		return
	}

	query := `
		UPDATE skills 
		SET name = $1, category = $2, proficiency = $3, icon = $4, "order" = $5
		WHERE id = $6
	`

	result, err := db.Exec(query, skill.Name, skill.Category, skill.Proficiency, skill.Icon, skill.Order, skillID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skill"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill updated successfully"})
}

func deleteSkill(c *gin.Context) {
	id := c.Param("id")

	// Input validation - ensure id is a valid integer
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	skillID, err := strconv.Atoi(id)
	if err != nil || skillID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID format"})
		return
	}

	query := `DELETE FROM skills WHERE id = $1`
	result, err := db.Exec(query, skillID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}

// =============================================================================
// 💼 EXPERIENCE HANDLERS
// =============================================================================

func getExperiences(c *gin.Context) {
	query := `
		SELECT id, title, company, start_date, end_date, current, description, technologies, "order", active
		FROM experience 
		WHERE active = true
		ORDER BY "order" DESC, start_date DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch experiences"})
		return
	}
	defer rows.Close()

	var experiences []Experience
	for rows.Next() {
		var exp Experience
		var technologies pq.StringArray
		err := rows.Scan(&exp.ID, &exp.Title, &exp.Company, &exp.StartDate, &exp.EndDate, &exp.Current, &exp.Description, &technologies, &exp.Order, &exp.Active)
		if err != nil {
			continue
		}

		exp.Technologies = []string(technologies)

		experiences = append(experiences, exp)
	}

	c.JSON(http.StatusOK, experiences)
}

func getExperience(c *gin.Context) {
	id := c.Param("id")

	var exp Experience
	var technologies pq.StringArray
	query := `SELECT id, title, company, start_date, end_date, current, description, technologies, "order", active FROM experience WHERE id = $1 AND active = true`

	err := db.QueryRow(query, id).Scan(&exp.ID, &exp.Title, &exp.Company, &exp.StartDate, &exp.EndDate, &exp.Current, &exp.Description, &technologies, &exp.Order, &exp.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch experience"})
		return
	}

	exp.Technologies = []string(technologies)

	c.JSON(http.StatusOK, exp)
}

func createExperience(c *gin.Context) {
	var exp Experience
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	technologiesJSON, _ := json.Marshal(exp.Technologies)
	query := `
		INSERT INTO experience (title, company, start_date, end_date, current, description, technologies, "order", active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var id int
	err := db.QueryRow(query, exp.Title, exp.Company, exp.StartDate, exp.EndDate, exp.Current, exp.Description, technologiesJSON, exp.Order, true).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create experience"})
		return
	}

	exp.ID = id
	c.JSON(http.StatusCreated, exp)
}

func updateExperience(c *gin.Context) {
	id := c.Param("id")
	var exp Experience
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	technologiesJSON, _ := json.Marshal(exp.Technologies)
	query := `
		UPDATE experience 
		SET title = $1, company = $2, start_date = $3, end_date = $4, current = $5, description = $6, technologies = $7, "order" = $8, active = $9
		WHERE id = $10
	`

	result, err := db.Exec(query, exp.Title, exp.Company, exp.StartDate, exp.EndDate, exp.Current, exp.Description, technologiesJSON, exp.Order, exp.Active, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update experience"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Experience updated successfully"})
}

func deleteExperience(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM experience WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete experience"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Experience deleted successfully"})
}

// =============================================================================
// 📄 CONTENT HANDLERS
// =============================================================================

func getContent(c *gin.Context) {
	query := `SELECT id, type, value FROM content ORDER BY type, id`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch content"})
		return
	}
	defer rows.Close()

	var contents []Content
	for rows.Next() {
		var content Content
		err := rows.Scan(&content.ID, &content.Type, &content.Value)
		if err != nil {
			continue
		}
		contents = append(contents, content)
	}

	c.JSON(http.StatusOK, contents)
}

func getContentByType(c *gin.Context) {
	contentType := c.Param("type")

	query := `SELECT id, type, value FROM content WHERE type = $1`

	rows, err := db.Query(query, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch content"})
		return
	}
	defer rows.Close()

	var contents []Content
	for rows.Next() {
		var content Content
		err := rows.Scan(&content.ID, &content.Type, &content.Value)
		if err != nil {
			continue
		}
		contents = append(contents, content)
	}

	c.JSON(http.StatusOK, contents)
}

func createContent(c *gin.Context) {
	var content Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	query := `
		INSERT INTO content (type, value)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := db.QueryRow(query, content.Type, content.Value).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create content"})
		return
	}

	content.ID = id
	c.JSON(http.StatusCreated, content)
}

func updateContent(c *gin.Context) {
	id := c.Param("id")
	var content Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	query := `UPDATE content SET type = $1, value = $2 WHERE id = $3`

	result, err := db.Exec(query, content.Type, content.Value, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update content"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content updated successfully"})
}

func deleteContent(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM content WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete content"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content deleted successfully"})
}

// =============================================================================
// 👤 ABOUT HANDLERS
// =============================================================================

func getAbout(c *gin.Context) {
	var description string
	err := db.QueryRow("SELECT value->>'description' FROM content WHERE key = 'about'").Scan(&description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch about data"})
		return
	}

	aboutData := AboutData{
		Description: description,
		Highlights: []struct {
			Icon string `json:"icon"`
			Text string `json:"text"`
		}{},
	}

	c.JSON(http.StatusOK, aboutData)
}

func updateAbout(c *gin.Context) {
	var aboutData AboutData
	if err := c.ShouldBindJSON(&aboutData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	query := `UPDATE content SET value = jsonb_set(value, '{description}', $1) WHERE key = 'about'`
	result, err := db.Exec(query, aboutData.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update about data"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "About content not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "About data updated successfully"})
}

// =============================================================================
// 📞 CONTACT HANDLERS
// =============================================================================

func getContact(c *gin.Context) {
	contactData := ContactData{
		Email:        getContentValue("contact", "email", "bruno@lucena.cloud"),
		Location:     getContentValue("contact", "location", "Brazil"),
		LinkedIn:     getContentValue("contact", "linkedin", "https://www.linkedin.com/in/bvlucena"),
		GitHub:       getContentValue("contact", "github", "https://github.com/brunovlucena"),
		Availability: getContentValue("contact", "availability", "Open to new opportunities"),
	}

	c.JSON(http.StatusOK, contactData)
}

func updateContact(c *gin.Context) {
	var contactData ContactData
	if err := c.ShouldBindJSON(&contactData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update each field individually
	fields := map[string]string{
		"email":        contactData.Email,
		"location":     contactData.Location,
		"linkedin":     contactData.LinkedIn,
		"github":       contactData.GitHub,
		"availability": contactData.Availability,
	}

	for field, value := range fields {
		query := `UPDATE content SET value = jsonb_set(value, $1, $2) WHERE key = 'contact'`
		_, err := db.Exec(query, "{"+field+"}", value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact data"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact data updated successfully"})
}

// =============================================================================
// 🛠️ UTILITY FUNCTIONS
// =============================================================================

func getContentValue(key, field, defaultValue string) string {
	var value string
	query := "SELECT value->>$1 FROM content WHERE key = $2"
	err := db.QueryRow(query, field, key).Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}

// =============================================================================
// 🔒 SECURITY HELPER FUNCTIONS
// =============================================================================

func isValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if host is not empty
	if parsedURL.Host == "" {
		return false
	}

	return true
}
