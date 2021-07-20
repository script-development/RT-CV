package models

type Period struct {
	Start   string `json:"start"` // iso 8601 time
	End     string `json:"end"`   // iso 8601 time
	Present bool   `json:"present"`
}

type Cv struct {
	Title                string           `json:"title"`
	ReferenceNumber      string           `json:"referenceNumber"`
	LastChanged          string           `json:"lastChanged"` // iso 8601 time
	Educations           []Education      `json:"educations"`
	Courses              []Course         `json:"courses"`
	WorkExperiences      []WorkExperience `json:"workExperiences"`
	PreferredJobs        []string         `json:"preferredJobs"`
	Languages            []Language       `json:"languages"`
	Competences          []Competence     `json:"competences"`
	Interests            []Interest       `json:"interests"`
	PersonalDetails      PersonalDetails  `json:"personalDetails"`
	PersonalPresentation string           `json:"personalPresentation"`
	DriversLicenses      []string         `json:"driversLicenses"`
	CreatedAt            *string          `json:"createdAt"` // iso 8601 time
}

type Education struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// TODO find difference between iscompleted and hasdiploma
	IsCompleted bool   `json:"isCompleted"`
	HasDiploma  bool   `json:"hasDiploma"`
	Period      Period `json:"period"`
	StartDate   string `json:"startDate"` // iso 8601 time
	EndDate     string `json:"endDate"`   // iso 8601 time
	Institute   string `json:"institute"`
	SubsectorID int    `json:"subsectorID"`
}

type Course struct {
	Name           string `json:"name"`
	NormalizedName string `json:"normalizedName"`
	StartDate      string `json:"startDate"` // iso 8601 time
	EndDate        string `json:"endDate"`   // iso 8601 time
	IsCompleted    bool   `json:"isCompleted"`
	Institute      string `json:"institute"`
	Description    string `json:"description"`
}

type WorkExperience struct {
	Description       string `json:"description"`
	Profession        string `json:"profession"`
	Period            Period `json:"period"`
	StartDate         string `json:"startDate"` // iso 8601 time
	EndDate           string `json:"endDate"`   // iso 8601 time
	StillEmployed     bool   `json:"stillEmployed"`
	Employer          string `json:"employer"`
	WeeklyHoursWorked int    `json:"weeklyHoursWorked"`
}

type LanguageProficiency int

type Language struct {
	Name         string              `json:"name"`
	LevelSpoken  LanguageProficiency `json:"levelSpoken"`
	LevelWritten LanguageProficiency `json:"levelWritten"`
}

type Competence struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Interest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PersonalDetails struct {
	Initials          string `json:"initials"`
	FirstName         string `json:"firstName"`
	SurNamePrefix     string `json:"surNamePrefix"`
	SurName           string `json:"surName"`
	Dob               string `json:"dob"` // iso 8601 time
	Gender            string `json:"gender"`
	StreetName        string `json:"streetName"`
	HouseNumber       string `json:"houseNumber"`
	HouseNumberSuffix string `json:"houseNumberSuffix"`
	Zip               string `json:"zip"`
	City              string `json:"city"`
	Country           string `json:"country"`
	PhoneNumber       string `json:"phoneNumber"`
	Email             string `json:"email"`
}
