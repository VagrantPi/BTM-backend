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

type SumsubDataResApplicantIdDoc struct {
	IdDocType    string `json:"idDocType,omitempty"`
	Country      string `json:"country,omitempty"`
	FirstName    string `json:"firstName,omitempty"`
	FirstNameEn  string `json:"firstNameEn,omitempty"`
	MiddleName   string `json:"middleName,omitempty"`
	MiddleNameEn string `json:"middleNameEn,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	LastNameEn   string `json:"lastNameEn,omitempty"`
	DateOfBirth  string `json:"dob,omitempty"` // yyyy-mm-dd format
	Number       string `json:"number,omitempty"`
	IssuedDate   string `json:"issuedDate,omitempty"` // yyyy-mm-dd format
}

type SumsubDataResApplicantInfo struct {
	FirstName    string                        `json:"firstName,omitempty"`
	FirstNameEn  string                        `json:"firstNameEn,omitempty"`
	MiddleName   string                        `json:"middleName,omitempty"`
	MiddleNameEn string                        `json:"middleNameEn,omitempty"`
	LastName     string                        `json:"lastName,omitempty"`
	LastNameEn   string                        `json:"lastNameEn,omitempty"`
	Dob          string                        `json:"dob,omitempty"` // yyyy-mm-dd format
	Gender       string                        `json:"gender,omitempty"`
	Country      string                        `json:"country,omitempty"`
	IdDocs       []SumsubDataResApplicantIdDoc `json:"idDocs,omitempty"`
}

type SumsubDataResApplicant struct {
	ID             string                     `json:"id,omitempty"`
	CreatedAt      string                     `json:"createdAt,omitempty"`
	Key            string                     `json:"key,omitempty"`
	ClientID       string                     `json:"clientId,omitempty"`
	InspectionID   string                     `json:"inspectionId,omitempty"`
	ExternalUserID string                     `json:"externalUserId,omitempty"`
	Info           SumsubDataResApplicantInfo `json:"info,omitempty"`
	FixedInfo      SumsubDataResApplicantInfo `json:"fixedInfo,omitempty"`
	Phone          string                     `json:"phone,omitempty"`
	Review         struct {
		ElapsedSincePendingMs int    `json:"elapsedSincePendingMs,omitempty"`
		ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs,omitempty"`
		Reprocessing          bool   `json:"reprocessing,omitempty"`
		CreateDate            string `json:"createDate,omitempty"`
		ReviewDate            string `json:"reviewDate,omitempty"`
		StartDate             string `json:"startDate,omitempty"`
		ReviewResult          struct {
			ReviewAnswer      string   `json:"reviewAnswer,omitempty"`
			RejectLabels      []string `json:"rejectLabels,omitempty"`
			ModerationComment string   `json:"moderationComment,omitempty"`
			ClientComment     string   `json:"clientComment,omitempty"`
		} `json:"reviewResult,omitempty"`
		ReviewStatus           string `json:"reviewStatus,omitempty"`
		NotificationFailureCnt int    `json:"notificationFailureCnt,omitempty"`
		Priority               int    `json:"priority,omitempty"`
		LevelName              string `json:"levelName,omitempty"`
	} `json:"review,omitempty"`
	Lang       string   `json:"lang,omitempty"`
	Type       string   `json:"type,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	RiskLabels struct {
		Email      []string `json:"email,omitempty"`
		Phone      []string `json:"phone,omitempty"`
		Device     []string `json:"device,omitempty"`
		CrossCheck []string `json:"crossCheck,omitempty"`
		Selfie     []string `json:"selfie,omitempty"`
		Aml        []string `json:"aml,omitempty"`
		Person     []string `json:"person,omitempty"`
	} `json:"riskLabels,omitempty"`
	Questionnaires []struct {
		Id       string  `json:"id,omitempty"`
		Score    float64 `json:"score,omitempty"`
		Sections struct {
			BasicInfo struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					LastName struct {
						Value string `json:"value,omitempty"`
					} `json:"lastName,omitempty"`
					FirstName struct {
						Value string `json:"value,omitempty"`
					} `json:"firstName,omitempty"`
					IDCategory struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"ID_category,omitempty"`
					IDLocation struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"ID_location,omitempty"`
					Adress struct {
						Value string `json:"value,omitempty"`
					} `json:"adress,omitempty"`
				}
			} `json:"Basic_info,omitempty"`
			JiBenZiLiao struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					KaiHuMuDe struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"kaiHuMuDe,omitempty"`
					NianShouRuJiJu struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"nianShouRuJiJu,omitempty"`
					ZhiCheng struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"zhiCheng,omitempty"`
					ZhiYe struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"zhiYe,omitempty"`
					XianRenZhiGongSiHuoJ struct {
						Value string `json:"value,omitempty"`
					} `json:"xianRenZhiGongSiHuoJ,omitempty"`
					RenZhiNianShu struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"renZhiNianShu,omitempty"`
					ZhangHaoYongTu struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"zhangHaoYongTu,omitempty"`
					RenZhiGongSiLeiXing struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"renZhiGongSiLeiXing,omitempty"`
				}
			} `json:"jiBenZiLiao,omitempty"`
			GaoZhiShiXiang struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					ShuoMingNeiRong struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"shuoMingNeiRong,omitempty"`
				}
			} `json:"gaoZhiShiXiang,omitempty"`
			ZiChanPingGu struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					GeRenZongZiChan struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"geRenZongZiChan,omitempty"`
					ZiJinLaiYuan struct {
						Values []string `json:"values,omitempty"`
						Score  float64  `json:"score,omitempty"`
					} `json:"ziJinLaiYuan,omitempty"`
					ZiChanZhengMingWenJi struct {
						Values []string `json:"values,omitempty"`
					} `json:"ziChanZhengMingWenJi,omitempty"`
				}
			} `json:"ziChanPingGu,omitempty"`
		}
	} `json:"questionnaires,omitempty"`
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

type SumsubTag string

func (e SumsubTag) String() string { return string(e) }

const (
	SumsubTagCib SumsubTag = "告誡名單"
)
