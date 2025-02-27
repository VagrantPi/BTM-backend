package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SumsubDataInfoAddress struct {
	Street           string `json:"street"`
	StreetEn         string `json:"streetEn"`
	Town             string `json:"town"`
	TownEn           string `json:"townEn"`
	Country          string `json:"country"`
	FormattedAddress string `json:"formattedAddress"`
	// LocationPosition interface{} `json:"locationPosition"`
}

type SumsubDataInfoIdDocsAddress struct {
	Street           string `json:"street"`
	StreetEn         string `json:"streetEn"`
	Town             string `json:"town"`
	TownEn           string `json:"townEn"`
	Country          string `json:"country"`
	FormattedAddress string `json:"formattedAddress"`
	// LocationPosition interface{} `json:"locationPosition"`
}

type SumsubDataInfoIdDocs struct {
	IdDocType    string                      `json:"idDocType"`
	Country      string                      `json:"country"`
	FirstName    string                      `json:"firstName"`
	FirstNameEn  string                      `json:"firstNameEn"`
	LastName     string                      `json:"lastName"`
	LastNameEn   string                      `json:"lastNameEn"`
	IssuedDate   string                      `json:"issuedDate"`
	Number       string                      `json:"number"`
	Dob          string                      `json:"dob"`
	Gender       string                      `json:"gender"`
	PlaceOfBirth string                      `json:"placeOfBirth"`
	Address      SumsubDataInfoIdDocsAddress `json:"address"`
}

type SumsubDataInfo struct {
	FirstName      string                  `json:"firstName"`
	FirstNameEn    string                  `json:"firstNameEn"`
	LastName       string                  `json:"lastName"`
	LastNameEn     string                  `json:"lastNameEn"`
	Dob            string                  `json:"dob"`
	Gender         string                  `json:"gender"`
	PlaceOfBirth   string                  `json:"placeOfBirth"`
	PlaceOfBirthEn string                  `json:"placeOfBirthEn"`
	Country        string                  `json:"country"`
	Addresses      []SumsubDataInfoAddress `json:"addresses"`
	IdDocs         []SumsubDataInfoIdDocs  `json:"idDocs"`
}

type SumsubDataAgreement struct {
	CreatedAt  string   `json:"createdAt"`
	AcceptedAt string   `json:"acceptedAt"`
	Source     string   `json:"source"`
	RecordIds  []string `json:"recordIds"`
}

type SumsubDataRequiredIdDocsDocSets struct {
	IdDocSetType  string   `json:"idDocSetType"`
	Types         []string `json:"types"`
	VideoRequired string   `json:"videoRequired"`
	CaptureMode   string   `json:"captureMode"`
	UploaderMode  string   `json:"uploaderMode"`
}

type SumsubDataRequiredIdDocs struct {
	IncludedCountries []string                          `json:"includedCountries"`
	DocSets           []SumsubDataRequiredIdDocsDocSets `json:"docSets"`
}

type SumsubDataReviewReviewResult struct {
	ReviewAnswer string `json:"reviewAnswer"`
}

type SumsubDataReview struct {
	ReviewId              string `json:"reviewId"`
	AttemptId             string `json:"attemptId"`
	AttemptCnt            int    `json:"attemptCnt"`
	ElapsedSincePendingMs int    `json:"elapsedSincePendingMs"`
	ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs"`
	Reprocessing          bool   `json:"reprocessing"`
	LevelName             string `json:"levelName"`
	// LevelAutoCheckMode    interface{}                  `json:"levelAutoCheckMode"`
	CreateDate   string                       `json:"createDate"`
	ReviewDate   string                       `json:"reviewDate"`
	ReviewResult SumsubDataReviewReviewResult `json:"reviewResult"`
	ReviewStatus string                       `json:"reviewStatus"`
	Priority     int                          `json:"priority"`
}

type SumsubDataRiskLabels struct {
	AttemptId string   `json:"attemptId"`
	CreatedAt string   `json:"createdAt"`
	Device    []string `json:"device"`
	Selfie    []string `json:"selfie"`
}

type SumsubData struct {
	Id                string                   `json:"id"`
	CreatedAt         string                   `json:"createdAt"`
	CreatedBy         string                   `json:"createdBy"`
	Key               string                   `json:"key"`
	ClientId          string                   `json:"clientId"`
	InspectionId      string                   `json:"inspectionId"`
	ExternalUserId    string                   `json:"externalUserId"`
	Info              SumsubDataInfo           `json:"info"`
	ApplicantPlatform string                   `json:"applicantPlatform"`
	IpCountry         string                   `json:"ipCountry"`
	AuthCode          string                   `json:"authCode"`
	Agreement         SumsubDataAgreement      `json:"agreement"`
	RequiredIdDocs    SumsubDataRequiredIdDocs `json:"requiredIdDocs"`
	Review            SumsubDataReview         `json:"review"`
	Lang              string                   `json:"lang"`
	Type              string                   `json:"type"`
	RiskLabels        SumsubDataRiskLabels     `json:"riskLabels"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (s *SumsubData) Scan(value interface{}) error {
	if value == nil {
		*s = SumsubData{}
		return nil
	}

	bytesValue, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid scan SumsubData")
	}

	data := SumsubData{}
	if err := json.Unmarshal(bytesValue, &data); err != nil {
		return errors.New("Invalid scan SumsubData unmarshal")
	}

	*s = data
	return nil
}

// Value return json value, implement driver.Valuer interface
func (s SumsubData) Value() (driver.Value, error) {
	return json.Marshal(s)
}
