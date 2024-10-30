package models

type URL struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

type ProjectItem struct {
	ID          string   `json:"id"`
	Visible     bool     `json:"visible"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Summary     string   `json:"summary"`
	Keywords    []string `json:"keywords"`
	URL         URL      `json:"url"`
}

type Projects struct {
	Name          string        `json:"name"`
	Columns       int           `json:"columns"`
	SeparateLinks bool          `json:"separateLinks"`
	Visible       bool          `json:"visible"`
	ID            string        `json:"id"`
	Items         []ProjectItem `json:"items"`
}

type SkillItem struct {
	ID       string   `json:"id"`
	Visible  bool     `json:"visible"`
	Name     string   `json:"name"`
	Level    int8     `json:"level"`
	Keywords []string `json:"keywords"`
}

type Skill struct {
	Name          string      `json:"name"`
	Columns       int         `json:"columns"`
	SeparateLinks bool        `json:"separateLinks"`
	Visible       bool        `json:"visible"`
	ID            string      `json:"id"`
	Items         []SkillItem `json:"items"`
}

type Experience struct {
	Name          string           `json:"name"`
	Columns       int              `json:"columns"`
	SeparateLinks bool             `json:"separateLinks"`
	Visible       bool             `json:"visible"`
	ID            string           `json:"id"`
	Items         []ExperienceItem `json:"items"`
}

type ExperienceItem struct {
	ID       string `json:"id"`
	Visible  bool   `json:"visible"`
	Company  string `json:"company"`
	Position string `json:"position"`
	Location string `json:"location"`
	Date     string `json:"date"`
	Summary  string `json:"summary"`
	Url      URL    `json:"url"`
}
type EducationItem struct {
	ID          string `json:"id"`
	Visible     bool   `json:"visible"`
	Institution string `json:"institution"`
	StudyType   string `json:"studyType"`
	Area        string `json:"area"`
	Score       string `json:"score"`
	Date        string `json:"date"`
	Summary     string `json:"summary"`
	URL         URL    `json:"url"`
}

type Education struct {
	Name          string          `json:"name"`
	Columns       int             `json:"columns"`
	SeparateLinks bool            `json:"separateLinks"`
	Visible       bool            `json:"visible"`
	ID            string          `json:"id"`
	Items         []EducationItem `json:"items"`
}

type Sections struct {
	Summary struct {
		Name          string `json:"name"`
		Columns       int    `json:"columns"`
		SeparateLinks bool   `json:"separateLinks"`
		Visible       bool   `json:"visible"`
		ID            string `json:"id"`
		Content       string `json:"content"`
	} `json:"summary"`
	Experience Experience `json:"experience"`
	Education  Education  `json:"education"`
	Skill      Skill      `json:"skills"`
	Projects   Projects   `json:"projects"`
}

type JsonData struct {
	Basics struct {
		Name     string `json:"name"`
		Headline string `json:"headline"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Location string `json:"location"`
		Url      URL    `json:"url"`
	} `json:"basics"`
	Sections Sections `json:"sections"`
}
